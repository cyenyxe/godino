package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type animal struct {
	id       int
	species  string
	nickname string
	zone     int
	age      int
}

func main() {
	// Connect to the database
	db, err := sql.Open("mysql", "dinoadmin:dinoadmin@/dino")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Retrieve animals above a certain age
	animals, err := queryByAge(db, 10)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(animals)
}

// queryByAge retrieves animals above a certain age
func queryByAge(db *sql.DB, age int) ([]animal, error) {
	rows, err := db.Query("select * from dino.animals where age > ?", age)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	animals := []animal{}
	for rows.Next() {
		a := animal{}
		err := rows.Scan(&a.id, &a.species, &a.nickname, &a.zone, &a.age)
		if err != nil {
			log.Println(err)
			continue
		}

		animals = append(animals, a)
	}

	err = rows.Err()

	return animals, err
}
