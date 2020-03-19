package main

import (
	"context"
	"fmt"
	"godino/dinogrpc"
	"log"

	grpc "google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8282", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := dinogrpc.NewDinoServiceClient(conn)

	animal, err := client.GetAnimal(context.Background(), &dinogrpc.Request{Nickname: "Velo"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Animal 1 : %v\n", animal) // Prints fields from animal

	animal, err = client.GetAnimal(context.Background(), &dinogrpc.Request{Nickname: "Trololo"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Animal 2 : %v\n", animal) // Prints an empty (non-existing) animal
}
