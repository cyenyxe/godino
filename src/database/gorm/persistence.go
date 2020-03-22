package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Animal represents a dinosaur specimen
// Struct attributes must be capitalized for GORM to detect them
type Animal struct {
	gorm.Model
	Species  string `gorm:"size:64"`
	Nickname string `gorm:"size:64;unique;not null"`
	Zone     int
	Age      int
}

// DinoDatabaseHandler is a handler for a database containing dinosaur information
type DinoDatabaseHandler struct {
	connection *gorm.DB
}

func main() {
	// Connect to the database
	// Without 'parseTime' set, the mapping breaks when it gets to the dates
	handler, err := NewDinoDatabaseHandler("dinoadmin:dinoadmin@/dino?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer handler.Close()

	// Retrieve animals above a certain age
	animals := handler.QueryByAge(10)
	fmt.Println(animals)

	// Retrieve animal with a certain ID
	a := handler.QueryByID(1)
	fmt.Println(a)

	// Insert a duplicate animal
	if _, err = handler.AddNewAnimal("Velociraptor", "Velo", 2, 20); err != nil {
		log.Println(err)
	}
}

// NewDinoDatabaseHandler creates a new handler for a database containing dinosaur information
func NewDinoDatabaseHandler(url string) (*DinoDatabaseHandler, error) {
	// Connect to the database
	db, err := gorm.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	handler := &DinoDatabaseHandler{
		connection: db,
	}

	// Initialize database with test data
	db.DropTableIfExists(&Animal{})
	db.AutoMigrate(&Animal{})

	if _, err = handler.AddNewAnimal("Tyrannousaurus rex", "T-Rex", 1, 10); err != nil {
		return handler, err
	}

	if _, err = handler.AddNewAnimal("Velociraptor", "Raptor", 2, 25); err != nil {
		return handler, err
	}

	if _, err = handler.AddNewAnimal("Velociraptor", "Velo", 2, 20); err != nil {
		return handler, err
	}

	return handler, nil
}

func (handler *DinoDatabaseHandler) Close() {
	handler.connection.Close()
}

// AddNewAnimal inserts a specimen in the dinosaur database
func (handler *DinoDatabaseHandler) AddNewAnimal(species string, nickname string, zone int, age int) (uint, error) {
	a := Animal{
		Species:  species,
		Nickname: nickname,
		Zone:     zone,
		Age:      age}

	err := handler.connection.Save(&a).Error
	return a.ID, err
}

// QueryByAge retrieves a specimen above a certain age
func (handler *DinoDatabaseHandler) QueryByAge(age int) []Animal {
	animals := []Animal{}
	handler.connection.Find(&animals, "age > ?", age)
	return animals
}

// QueryByID retrieves a specimen with a certain ID
func (handler *DinoDatabaseHandler) QueryByID(id int) Animal {
	a := Animal{}
	handler.connection.First(&a, id)
	return a
}
