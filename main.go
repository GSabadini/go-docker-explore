package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	var docker = DockerAPI{
		client: cli,
		ctx:    ctx,
	}

	//docker.RunContainer()
	//docker.RunContainerBackground()
	//docker.StopContainer()
	//docker.PrintLogsContainer()
	//docker.ListContainers()
	//docker.ListContainers("alpine")
	docker.ListImages()
}

type DockerAPI struct {
	client *client.Client
	ctx    context.Context
}

func (d DockerAPI) PullImage(imageName string) {
	out, err := d.client.ImagePull(d.ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	defer out.Close()

	io.Copy(os.Stdout, out)
}

func (d DockerAPI) ListImages() {
	images, err := d.client.ImageList(d.ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		fmt.Println(image.Containers)
		fmt.Println(image.Created)
		fmt.Println(image.Size)
		fmt.Println(image.ParentID)
	}
}

func (d DockerAPI) PrintLogsContainer(containerID string) {
	options := types.ContainerLogsOptions{ShowStdout: true}
	// Replace this ID with a container that really exists
	out, err := d.client.ContainerLogs(d.ctx, containerID, options)
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)
}

func (d DockerAPI) StopContainer() {
	containers, err := d.client.ContainerList(d.ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Print("Stopping container ", container.ID[:10], "... ")
		if err := d.client.ContainerStop(d.ctx, container.ID, nil); err != nil {
			panic(err)
		}
		fmt.Println("Success")
	}
}

func (d DockerAPI) RunContainerBackground(imageName string) {
	//imageName := "docker.io/library/alpine"
	out, err := d.client.ImagePull(d.ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)

	resp, err := d.client.ContainerCreate(d.ctx, &container.Config{
		Image: imageName,
	}, nil, nil, "test")
	if err != nil {
		panic(err)
	}

	if err := d.client.ContainerStart(d.ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)
}


func (d DockerAPI) RunContainer(imageName string) {
	reader, err := d.client.ImagePull(d.ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := d.client.ContainerCreate(d.ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
		Tty:   true,
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := d.client.ContainerStart(d.ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	_, err = d.client.ContainerWait(d.ctx, resp.ID)
	if err != nil {
		panic(err)
	}

	out, err := d.client.ContainerLogs(d.ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}

func (d DockerAPI) ListContainers() {
	containers, err := d.client.ContainerList(d.ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s %s\n", container.ID[:10], container.Image, container.Names)
	}
}