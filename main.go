package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dockerclient "github.com/docker/docker/client"
	"github.com/vllry/k2/api/0.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"math/rand"
	"net"
)

type server struct {
	db *dbConn
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (s *server) CreateContainer(ctx context.Context, query *api.CreateContainerRequest) (*api.CreateContainerResult, error) {
	//containerName := randStringBytes(5)
	containerName := query.Image  // Use a deterministic name for now.
	containerStatusPath := "/containers/" + containerName + "/status"
	containerIdPath := "/containers/" + containerName + "/id"

	fmt.Println("Trying to start " + containerName)
	_, foundStatus, err := s.db.get(containerStatusPath)
	if err != nil {
		fmt.Println("failed to get status")
		return nil, err
	}
	if !foundStatus {  // Start if not already started.
		addContainer(s.db, containerName, query.Image, query.ImageTag)
		return &api.CreateContainerResult{Success: true}, err
	} else {  // Check that logged container is actually/still running.
		fmt.Println("Check on previously-running container")
		containerId, foundContainerId, err := s.db.get(containerIdPath)
		if err != nil {
			fmt.Println("Error getting container details from database", err)
			return nil, err
		}

		if foundContainerId {
			cli := dockerclient.Client{}
			info, err := cli.ContainerInspect(context.Background(), containerId)
			if err != nil {
				fmt.Println("Error getting Docker details", err)
				return nil, err
			}
			fmt.Println(info)
		} else {  // Indicates bug - container was half tracked.
			err = addContainer(s.db, containerName, query.Image, query.ImageTag)
			if err != nil {
				return &api.CreateContainerResult{Success: false}, err
			}
			return &api.CreateContainerResult{Success: true}, nil
		}
	}
	// Does anything reach this?
	return &api.CreateContainerResult{Success:false}, nil
}

func addContainer(db *dbConn, containerName string, image string, tag string) error {
	// TODO split into consts or builder funcs.
	containerStatusPath := "/containers/" + containerName + "/status"
	containerIdPath := "/containers/" + containerName + "/id"

	fmt.Println("Create container")
	db.set(containerStatusPath, "PENDING")
	containerId, err := createContainerFromImage(containerName, image, tag)
	if err != nil {
		fmt.Println("Failed to create")
		return err
	}
	db.set(containerIdPath, containerId)
	db.set(containerStatusPath, "STARTED")
	return nil
}

func createContainerFromImage(name string, image string, tag string) (string, error) {
	ctx := context.Background()

	cli, err := dockerclient.NewEnvClient()
	if err != nil {
		return "", err
	}

	containerRef, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: fmt.Sprintf("%s:%s", image, tag),
			Tty:   true,
		},
		nil,
		nil,
		name,
	)
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, containerRef.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	//containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	//if err != nil {
	//	return err
	//}

	//for _, containerRef := range containers {
	//	fmt.Printf("%s %s\n", containerRef.ID[:10], containerRef.Image)
	//}

	return containerRef.ID, nil
}

func main() {
	db, err := newDb()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to db.")

	listener, err := net.Listen("tcp", ":50051")
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
