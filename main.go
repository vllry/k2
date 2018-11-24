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
	fmt.Println("CreateContainer call")
	containerName := randStringBytes(5)
	s.db.set("/containers/status" + containerName, "PENDING")
	err := createContainerFromImage(containerName, query.Image, query.ImageTag)
	if err != nil {
		return nil, err
	}
	s.db.set("/containers/status" + containerName, "STARTED")

	return &api.CreateContainerResult{Success: true}, err
}

func createContainerFromImage(name string, image string, tag string) error {
	ctx := context.Background()

	cli, err := dockerclient.NewEnvClient()
	if err != nil {
		return err
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
		return err
	}

	if err := cli.ContainerStart(ctx, containerRef.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return err
	}

	for _, containerRef := range containers {
		fmt.Printf("%s %s\n", containerRef.ID[:10], containerRef.Image)
	}

	return nil
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
