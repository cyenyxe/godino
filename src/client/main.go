package main

import (
	"context"
	"godino/dinogrpc"
	"io"
	"log"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
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
		log.Panic(err)
	}
	grpclog.Printf("Animal 1 : %v\n", animal) // Prints fields from animal

	animal, err = client.GetAnimal(context.Background(), &dinogrpc.Request{Nickname: "Trololo"})
	if err != nil {
		log.Panic(err)
	}
	grpclog.Printf("Animal 2 : %v\n", animal) // Prints an empty (non-existing) animal

	animals, err := client.GetAllAnimals(context.Background(), &dinogrpc.Request{})
	if err != nil {
		log.Panic(err)
	}

	for {
		animal, err := animals.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			grpclog.Fatal(err)
		}
		grpclog.Printf("Animal: %v\n", animal)
	}

}
