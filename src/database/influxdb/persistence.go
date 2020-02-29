package main

import (
	"fmt"
	"log"
	"math/rand"
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

	// Create batch points
	batch, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Specimens in the reserve
	nicknameTags := []string{"T-Rex", "Raptor", "Velo"}
	speciesTags := []string{"Tyrannousaurus rex", "Velociraptor", "Velociraptor"}

	// Generate data points every 0.1 second
	rand.Seed(time.Now().Unix())
	for {
		// Populate data point for a random specimen with random values
		i := rand.Intn(len(nicknameTags))
		point, err := generateHealthMetrics(speciesTags[i], nicknameTags[i])
		if err != nil {
			log.Println(err)
			continue
		}

		batch.AddPoint(point)
		// Inefficient, no batch writes in practice
		if err = c.Write(batch); err != nil {
			log.Fatal(err)
		}

		time.Sleep(100 * time.Millisecond)
	}

	// // Retrieve animals above a certain age
	// animals := queryByAge(ctx, collection, 10)
	// fmt.Println(animals)

	// // // Retrieve animal with a certain ID
	// // a := queryByID(db, 1)
	// // fmt.Println(a)

	// // Insert a duplicate animal
	// if err = addNewAnimal(ctx, collection, "Velociraptor", "Velo", 2, 20); err != nil {
	// 	log.Println(err)
	// }
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

func generateHealthMetrics(species string, nickname string) (*client.Point, error) {
	tags := map[string]string{
		"species":  species,
		"nickname": nickname,
	}
	fields := map[string]interface{}{
		"weight":      rand.Intn(500) + 1,
		"temperature": rand.Intn(5) + 36,
	}
	fmt.Println(tags, fields["weight"], fields["temperature"])

	return client.NewPoint("health", tags, fields, time.Now())
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
