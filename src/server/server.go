package server

import (
	"context"
	"godino/models"

	"google.golang.org/grpc"
)

// DinoGrpcServer serves dinosaur data via gRPC
type DinoGrpcServer struct {
	Dinosaurs map[string]models.Animal
}

// NewDinoGrpcServer creates a new gRPC server for dinosaur data
func NewDinoGrpcServer() *DinoGrpcServer {
	dinos := map[string]models.Animal{}
	dinos["Raptor"] = models.Animal{
		Species:  "Velociraptor",
		Nickname: "Raptor",
		Age:      40,
		Zone:     2,
	}
	dinos["Velo"] = models.Animal{
		Species:  "Velociraptor",
		Nickname: "Velo",
		Age:      30,
		Zone:     2,
	}
	dinos["T-Rex"] = models.Animal{
		Species:  "Tyrannosaurus Rex",
		Nickname: "T-Rex",
		Age:      25,
		Zone:     1,
	}
	return &DinoGrpcServer{dinos}
}

// GetAnimal retrieves the animal idenfied by the nickname in the request
func (server *DinoGrpcServer) GetAnimal(ctx context.Context, in *models.Request, opts ...grpc.CallOption) (*models.Animal, error) {
	dinosaur := server.Dinosaurs[in.GetNickname()]
	return &dinosaur, nil
}

// GetAllAnimals retrieves all the animals and streams their data to the client
func (server *DinoGrpcServer) GetAllAnimals(in *Request, stream DinoService_GetAllAnimalsServer) error {
	for _, animal := range server.Dinosaurs {
		if err := stream.Send(animal); err != nil {
			return err
		}
	}

	return nil
}
