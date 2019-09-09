package main

import (
	"fmt"
	"context"
	"errors"
	"github.com/micro/go-micro"

	// Import the generated protobuf code
	pb "github.com/tresko/shippy-service-vessel/proto/vessel"
)

type repository interface {
	FindAvailableVesel(specification *pb.Specification) (*pb.Vessel, error)
}

// Repository - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later on.
type Repository struct {
	vessels []*pb.Vessel
}
// FindAvailableView - checks a specification against a map of vessels,
// if capacity and max weight are below a vessels capacity and max weight,
// then return that vessel.
func (repo *Repository) FindAvailableVesel(spec *pb.Specification) (*pb.Vessel, error) {
	for _, vessel := range repo.vessels {
		if spec.Capacity <= vessel.Capacity && spec.MaxWeight <= vessel.MaxWeight {
			return vessel, nil
		}
	}
	return nil, errors.New("No vessel found by that spec")
}

// Our grpc service handler
type VesselService struct {
	repo repository
}

func (s *VesselService) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {

	// Find the next available vessel
	vessel, err := s.repo.FindAvailableVesel(req)
	if err != nil {
		return err
	}

	// Set the vessel as part of the response message type
	res.Vessel = vessel
	return nil
}

func main() {
	vessels := []*pb.Vessel{
		&pb.Vessel{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}

	repo := &Repository{vessels}

	// Create a new service. Optionally include some options here.
	srv := micro.NewService(
		// This name must match the package name given in your protobuf definition
		micro.Name("shippy.vessel.service"),
	)

	// Init will parse the command line flags.
	srv.Init()

	// Register handler
	pb.RegisterVesselServiceHandler(srv.Server(), &VesselService{repo})

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
