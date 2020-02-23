package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Struct attributes must be capitalized for GORM to detect them
type animal struct {
	//gorm.Model
	ID       uint
	Species  string `gorm:"size:64"`
	Nickname string `gorm:"size:64;unique;not null"`
	Zone     int
	Age      int
}

func main() {
	// Connect to the database
	db, err := gorm.Open("mysql", "dinoadmin:dinoadmin@/dino")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize database with test data
	db.DropTableIfExists(&animal{})
	db.AutoMigrate(&animal{})

	if _, err = addNewAnimal(db, "Tyrannousaurus rex", "T-Rex", 1, 10); err != nil {
		log.Fatal(err)
	}

	if _, err = addNewAnimal(db, "Velociraptor", "Raptor", 2, 25); err != nil {
		log.Fatal(err)
	}

	if _, err = addNewAnimal(db, "Velociraptor", "Velo", 2, 20); err != nil {
		log.Fatal(err)
	}

	// Retrieve animals above a certain age
	animals := queryByAge(db, 10)
	fmt.Println(animals)

	// Retrieve animal with a certain ID
	a := queryByID(db, 1)
	fmt.Println(a)

	// Insert a duplicate animal
	if _, err = addNewAnimal(db, "Velociraptor", "Velo", 2, 20); err != nil {
		log.Println(err)
	}
}

// queryByAge retrieves animals above a certain age
func queryByAge(db *gorm.DB, age int) []animal {
	animals := []animal{}
	db.Find(&animals, "age > ?", age)
	return animals
}

func queryByID(db *gorm.DB, id int) animal {
	a := animal{}
	db.First(&a, id)
	return a
}

func addNewAnimal(db *gorm.DB, species string, nickname string, zone int, age int) (uint, error) {
	a := animal{
		Species:  species,
		Nickname: nickname,
		Zone:     zone,
		Age:      age}

	err := db.Save(&a).Error
	return a.ID, err
}
