package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Initialize database with test data
	databaseName := "dino"
	collectionName := "animals"
	collection := client.Database(databaseName).Collection(collectionName)
	collection.DeleteMany(ctx, bson.M{})

	if err := addNewAnimal(ctx, collection, "Tyrannousaurus rex", "T-Rex", 1, 10); err != nil {
		log.Fatal(err)
	}

	if err = addNewAnimal(ctx, collection, "Velociraptor", "Raptor", 2, 25); err != nil {
		log.Fatal(err)
	}

	if err = addNewAnimal(ctx, collection, "Velociraptor", "Velo", 2, 20); err != nil {
		log.Fatal(err)
	}

	// Retrieve animals above a certain age
	animals := queryByAge(ctx, collection, 10)
	fmt.Println(animals)

	// // Retrieve animal with a certain ID
	// a := queryByID(db, 1)
	// fmt.Println(a)

	// Insert a duplicate animal
	if err = addNewAnimal(ctx, collection, "Velociraptor", "Velo", 2, 20); err != nil {
		log.Println(err)
	}
}

// queryByAge retrieves animals above a certain age
func queryByAge(ctx context.Context, collection *mongo.Collection, age int) []animal {
	var animals []animal
	cur, err := collection.Find(ctx, bson.M{"age": bson.M{"$gt": age}})
	if err != nil {
		return animals
	}

	for cur.Next(ctx) {
		// create a value into which the single document can be decoded
		var elem animal
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		animals = append(animals, elem)
	}

	return animals
}

// func queryByID(db *gorm.DB, id int) animal {
// 	a := animal{}
// 	db.First(&a, id)
// 	return a
// }

func addNewAnimal(ctx context.Context, collection *mongo.Collection,
	species string, nickname string, zone int, age int) error { // (primitive.ObjectID, error) {
	a := animal{
		Species:  species,
		Nickname: nickname,
		Zone:     zone,
		Age:      age}

	_, err := collection.InsertOne(ctx, a)
	// TODO getting the inserted ID causes a panic
	//return insertResult.InsertedID.(primitive.ObjectID), err
	return err
}
