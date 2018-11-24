package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/vllry/k2/api/0.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type server struct{}

func (s *server) CreateContainer(ctx context.Context, query *api.CreateContainerRequest) (*api.CreateContainerResult, error) {
	fmt.Println("grpc")
	err := createContainerFromImage(query.Image, query.ImageTag)

	return &api.CreateContainerResult{Success: true}, err
}

func createContainerFromImage(image string, tag string) error {
	ctx := context.Background()

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containerRef, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: fmt.Sprintf("%s:%s", image, tag),
			Tty:   true,
		},
		nil,
		nil,
		"",
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
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	api.RegisterKube2Server(s, &server{})
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
