package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/influxdata/influxdb1-client"
	client "github.com/influxdata/influxdb1-client/v2"
)

// Struct attributes must be capitalized for GORM to detect them
type animal struct {
	Species  string
	Nickname string
	Zone     int
	Age      int
}

func main() {
	// Connect to the database
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "dinoadmin",
		Password: "dinoadmin",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	db := "dino"

	// Specimens in the reserve
	nicknameTags := []string{"T-Rex", "Raptor", "Velo"}
	speciesTags := []string{"Tyrannousaurus rex", "Velociraptor", "Velociraptor"}

	rand.Seed(time.Now().Unix())

	channel := make(chan *client.Point, 100)

	for i := 0; i < len(nicknameTags); i++ {
		// Generate random data points for each specimen
		go generateHealthMetrics(speciesTags[i], nicknameTags[i], channel)
	}

	// Create batch for data points
	// Precision will restrict the time granularity of data points may seem duplicates
	// This program can create data points very fast so 'ns' is needed
	batch, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db,
		Precision: "ns",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Data points will be generated until a signal is captured
	// TODO Possible race condition when accessing 'wg' because it is shared between main thread and goroutine?
	wg := sync.WaitGroup{}

	for detectSignal := checkStopOSSignals(&wg); !*detectSignal; {
		p := <-channel
		batch.AddPoint(p)

		if len(batch.Points()) >= 5000 || *detectSignal {
			fmt.Printf("Writing %d items to database...\n", len(batch.Points()))
			if err = c.Write(batch); err != nil {
				log.Fatal(err)
			}

			batch, err = client.NewBatchPoints(client.BatchPointsConfig{
				Database:  db,
				Precision: "ns",
			})
		}
	}

	wg.Wait()

	// Now query the data that has been inserted
	parameters := map[string]interface{}{
		"nickname": "Velo",
	}
	result, err := query(c, db, "select * from health where nickname = $nickname limit 10", parameters)
	if err != nil {
		log.Fatal(err)
	}

	for _, value := range result {
		log.Println("Messages: ", value.Messages)
		for _, s := range value.Series {
			log.Println("Name: ", s.Name)
			log.Println("Columns: ", s.Columns)
			log.Println("Values: ", s.Values)
		}
	}
}

func query(c client.Client, db string, query string, parameters map[string]interface{}) (results []client.Result, err error) {
	q := client.Query{
		Command:    query,
		Database:   db,
		Parameters: parameters,
	}

	response, err := c.Query(q)
	if err != nil {
		return results, err
	}
	if response.Error() != nil {
		return results, response.Error()
	}
	return response.Results, nil
}

func generateHealthMetrics(species string, nickname string, ch chan *client.Point) {
	for {
		tags := map[string]string{
			"species":  species,
			"nickname": nickname,
		}
		fields := map[string]interface{}{
			"weight":      rand.Intn(500) + 1,
			"temperature": rand.Intn(5) + 36,
		}

		// fmt.Println(tags, fields["weight"], fields["temperature"])

		point, err := client.NewPoint("health", tags, fields, time.Now())
		if err != nil {
			log.Println(err)
		}

		ch <- point
	}
}

func checkStopOSSignals(wg *sync.WaitGroup) *bool {
	Signal := false
	go func(s *bool) {
		wg.Add(1)
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		log.Println("Exit signals received... ")
		*s = true
		wg.Done()
	}(&Signal)
	return &Signal
}
