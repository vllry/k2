package main

import (
	"context"
	"fmt"
	dockerclient "github.com/docker/docker/client"
	"github.com/vllry/k2/api/0.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strings"
)

type server struct {
	db *dbConn
}

func (s *server) CreateContainer(ctx context.Context, query *api.CreateContainerRequest) (*api.CreateContainerResult, error) {
	//containerName := randStringBytes(5)
	containerName := query.Image // Use a deterministic name for now.
	containerStatusPath := "/containers/" + containerName + "/status"
	containerIdPath := "/containers/" + containerName + "/id"

	fmt.Println("Trying to start " + containerName)
	_, foundStatus, err := s.db.get(containerStatusPath)
	if err != nil {
		fmt.Println("failed to get status")
		return nil, err
	}
	if !foundStatus { // Start if not already started.
		addContainer(s.db, containerName, query.Image, query.ImageTag)
		return &api.CreateContainerResult{Success: true}, err
	} else { // Check that logged WorkloadContainerSpec is actually/still running.
		fmt.Println("Check on previously-running WorkloadContainerSpec")
		containerId, foundContainerId, err := s.db.get(containerIdPath)
		if err != nil {
			fmt.Println("Error getting WorkloadContainerSpec details from database", err)
			return nil, err
		}

		if foundContainerId {
			cli, err := dockerclient.NewEnvClient()
			_, err = cli.ContainerInspect(context.Background(), containerId)
			if err != nil {
				fmt.Println(err)
				if strings.Contains(err.Error(), "No such WorkloadContainerSpec") {
					fmt.Println("Create missing WorkloadContainerSpec.")
					err = addContainer(s.db, containerName, query.Image, query.ImageTag)
					if err != nil {
						return &api.CreateContainerResult{Success: false}, err
					}
					return &api.CreateContainerResult{Success: true}, nil
				} else {
					fmt.Println("Error getting Docker details", err)
					return nil, err
				}
			}
		} else { // Indicates bug - WorkloadContainerSpec was half tracked.
			err = addContainer(s.db, containerName, query.Image, query.ImageTag)
			if err != nil {
				return &api.CreateContainerResult{Success: false}, err
			}
			return &api.CreateContainerResult{Success: true}, nil
		}
	}
	// Does anything reach this?
	return &api.CreateContainerResult{Success: false}, nil
}

func main() {
	db, err := newDb()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to db.")

	//go scheduleController(db)
	startWorkloadApiServer(db)

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	api.RegisterKube2Server(s, &server{db})
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
