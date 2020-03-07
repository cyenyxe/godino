package main

import (
	"fmt"
	"godino/models"
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/protobuf/proto"
)

func main() {
	l, err := net.Listen("tcp", ":8282")
	if err != nil {
		log.Fatal(err)
	}

	// TODO Stop condition
	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		defer l.Close()

		go func(c net.Conn) {
			defer c.Close()
			data, err := ioutil.ReadAll(c)
			if err != nil {
				return
			}

			a := &models.Animal{}
			err = proto.Unmarshal(data, a)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println(a)
		}(c)
	}
}
