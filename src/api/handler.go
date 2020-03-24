package dinoapi

import (
	"github.com/jinzhu/gorm"

	// Needed for database initialization
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

// Close closes the database connection
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

// FindAll retrieves all specimens
func (handler *DinoDatabaseHandler) FindAll() []Animal {
	animals := []Animal{}
	handler.connection.Find(&animals)
	return animals
}

// FindByAge retrieves a specimen above a certain age
func (handler *DinoDatabaseHandler) FindByAge(age int) []Animal {
	animals := []Animal{}
	handler.connection.Find(&animals, "age > ?", age)
	return animals
}

// FindByID retrieves a specimen with a certain ID
func (handler *DinoDatabaseHandler) FindByID(id int) Animal {
	a := Animal{}
	handler.connection.First(&a, id)
	return a
}
