package dinogrpc

import (
	"context"
)

// Server serves dinosaur data via gRPC
type Server struct {
	Dinosaurs map[string]Animal
}

// NewServer creates a new gRPC server for dinosaur data
func NewServer() *Server {
	dinos := map[string]Animal{}
	dinos["Raptor"] = Animal{
		Species:  "Velociraptor",
		Nickname: "Raptor",
		Age:      40,
		Zone:     2,
	}
	dinos["Velo"] = Animal{
		Species:  "Velociraptor",
		Nickname: "Velo",
		Age:      30,
		Zone:     2,
	}
	dinos["T-Rex"] = Animal{
		Species:  "Tyrannosaurus Rex",
		Nickname: "T-Rex",
		Age:      25,
		Zone:     1,
	}
	return &Server{dinos}
}

// GetAnimal retrieves the animal idenfied by the nickname in the request
func (server *Server) GetAnimal(ctx context.Context, in *Request) (*Animal, error) {
	dinosaur := server.Dinosaurs[in.GetNickname()]
	return &dinosaur, nil
}

// GetAllAnimals retrieves all the animals and streams their data to the client
func (server *Server) GetAllAnimals(in *Request, stream DinoService_GetAllAnimalsServer) error {
	for _, animal := range server.Dinosaurs {
		if err := stream.Send(&animal); err != nil {
			return err
		}
	}

	return nil
}
