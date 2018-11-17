package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"time"
)

type server struct{}

func (s *server) CreateContainer(ctx context.Context, query *CreateContainerRequest) (*CreateContainerResult, error) {
	fmt.Println("grpc")
	err := createContainerFromImage(query.Image, query.ImageTag)

	return &CreateContainerResult{Success: true}, err
}

func createContainerFromImage(image string, tag string) error {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts()
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
	if os.Args[1] == "client" {
		// Set up a connection to the server.
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := NewKube2Client(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		r, err := c.CreateContainer(ctx, &CreateContainerRequest{Image:"nginx", ImageTag:"latest"})
		if err != nil {
			log.Fatalf("could not launch: %v", err)
		}
		log.Printf("Status: %b", r.Success)
	} else {
		listener, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		RegisterKube2Server(s, &server{})
		reflection.Register(s)
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}
}
