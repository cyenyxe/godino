package main

import (
	"encoding/json"
	api "godino/api"
	"log"
	"os"
)

type configuration struct {
	Host string
	Port int
}

func main() {
	// Open and read JSON file
	file, err := os.Open("configuration.json")
	if err != nil {
		log.Fatal(err)
	}
	config := new(configuration)
	json.NewDecoder(file).Decode(config)

	// Run API server
	err = api.RunWebPortalAPI(config.Host, config.Port)
	if err != nil {
		log.Fatal(err)
	}

}
