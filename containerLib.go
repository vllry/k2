package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dockerclient "github.com/docker/docker/client"
)

func addContainer(db *dbConn, containerName string, image string, tag string) error {
	// TODO split into consts or builder funcs.
	containerStatusPath := "/containers/" + containerName + "/status"
	containerIdPath := "/containers/" + containerName + "/id"

	fmt.Println("Create WorkloadContainerSpec")
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
