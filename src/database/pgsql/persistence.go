package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
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
	db, err := sql.Open("postgres", "user=dinoadmin password=dinoadmin dbname=dino sslmode=disable")
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

	// Retrieve animal with a certain ID
	a, err := queryByID(db, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(a)

	// Insert a new animal
	newID, err := addNewAnimal(db, "Carnotaurus", "Carnitas", 3, 30)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(newID)

	if err = queryByAgeMultiple(db, 10, 20, 50); err != nil {
		log.Fatal(err)
	}
}

func queryByAgeMultiple(db *sql.DB, ages ...int) error {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("select * from dino.animals where age > $1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, age := range ages {
		rows, err := stmt.Query(age)
		if err != nil {
			log.Fatal(err)
		}
		animals, err := handleRows(rows)
		if err == nil {
			fmt.Println(animals)
		}
	}

	return tx.Commit()
}

// queryByAge retrieves animals above a certain age
func queryByAge(db *sql.DB, age int) ([]animal, error) {
	rows, err := db.Query("select * from dino.animals where age > $1", age)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	return handleRows(rows)
}

func handleRows(rows *sql.Rows) ([]animal, error) {
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

	err := rows.Err()

	return animals, err
}

func queryByID(db *sql.DB, id int) (animal, error) {
	row := db.QueryRow("select * from dino.animals where id = $1", id)

	a := animal{}
	err := row.Scan(&a.id, &a.species, &a.nickname, &a.zone, &a.age)

	return a, err
}

func addNewAnimal(db *sql.DB, species string, nickname string, zone int, age int) (int64, error) {
	// result.LastInsertId() is not supported
	var id int64
	db.QueryRow("insert into dino.animals (species, nickname, zone, age) values ($1, $2, $3, $4) returning id",
		species, nickname, zone, age).Scan(&id)

	return id, nil
}
