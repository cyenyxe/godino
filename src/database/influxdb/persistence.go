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
	numSpecimens := len(nicknameTags)

	channels := make([](chan *client.Point), numSpecimens)
	done := make(chan bool)

	for i := 0; i < len(nicknameTags); i++ {
		channels[i] = make(chan *client.Point)
		// Generate random data points for each specimen
		go generateHealthMetrics(speciesTags[i], nicknameTags[i], channels[i], done)
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
	wg := sync.WaitGroup{}
	detectSignal := checkStopOSSignals(&wg)

	for openChannels := numSpecimens; openChannels > 0 && !*detectSignal; {
		select {
		case p := <-channels[0]:
			batch.AddPoint(p)
		case p := <-channels[1]:
			batch.AddPoint(p)
		case p := <-channels[2]:
			batch.AddPoint(p)
		case <-done:
			openChannels--
		}

		if len(batch.Points()) >= 5000 || openChannels == 0 {
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
	close(done)
	// // Retrieve animals above a certain age
	// animals := queryByAge(ctx, collection, 10)
	// fmt.Println(animals)

	// // Retrieve animal with a certain ID
	// a := queryByID(db, 1)
	// fmt.Println(a)
}

func query(c client.Client, db string, query string) (results []client.Result, err error) {
	q := client.Query{
		Command:  query,
		Database: db,
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

func generateHealthMetrics(species string, nickname string, ch chan *client.Point, done chan bool) { //} (*client.Point, error) {
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

	done <- true
	close(ch)
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

// queryByAge retrieves animals above a certain age
// func queryByAge(ctx context.Context, collection *mongo.Collection, age int) []animal {
// 	var animals []animal
// 	cur, err := collection.Find(ctx, bson.M{"age": bson.M{"$gt": age}})
// 	if err != nil {
// 		return animals
// 	}

// 	for cur.Next(ctx) {
// 		// create a value into which the single document can be decoded
// 		var elem animal
// 		err := cur.Decode(&elem)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		animals = append(animals, elem)
// 	}

// 	return animals
// }

// func queryByID(db *gorm.DB, id int) animal {
// 	a := animal{}
// 	db.First(&a, id)
// 	return a
// }

// func addNewAnimal(ctx context.Context, collection *mongo.Collection,
// 	species string, nickname string, zone int, age int) error { // (primitive.ObjectID, error) {
// 	a := animal{
// 		Species:  species,
// 		Nickname: nickname,
// 		Zone:     zone,
// 		Age:      age}

// 	_, err := collection.InsertOne(ctx, a)
// 	// TODO getting the inserted ID causes a panic
// 	//return insertResult.InsertedID.(primitive.ObjectID), err
// 	return err
// }
