package main

import (
	"godino/dinogrpc"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	listener, err := net.Listen("tcp", ":8282")
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}

	grpclog.Println("Listening on port 8282")

	var opts []grpc.ServerOption
	server := grpc.NewServer(opts...)
	dinoServer := dinogrpc.NewServer()

	dinogrpc.RegisterDinoServiceServer(server, dinoServer)

	if err = server.Serve(listener); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
