package main

import (
	"godino/models"
	"log"
	"net"

	"google.golang.org/protobuf/proto"
)

func main() {
	a := &models.Animal{
		Id:       1,
		Species:  "Velociraptor",
		Nickname: "Raptor",
		Zone:     3,
		Age:      21,
	}

	data, err := proto.Marshal(a)
	if err != nil {
		log.Fatal(err)
	}

	sendData(data)
}

func sendData(data []byte) {
	c, err := net.Dial("tcp", "127.0.0.1:8282")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	c.Write(data)
}
