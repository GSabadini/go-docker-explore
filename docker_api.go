package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type dockerAPI struct {
	client *client.Client
	ctx    context.Context
}

func newDockerAPI(client *client.Client, ctx context.Context) dockerAPI {
	return dockerAPI{client: client, ctx: ctx}
}

func (d dockerAPI) RemoveContainersExited() {
	args := filters.NewArgs()
	args.Add("status", "exited")

	containers, err := d.client.ContainerList(d.ctx, types.ContainerListOptions{
		Filters: args,
	})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		err := d.client.ContainerRemove(d.ctx, container.ID[:10], types.ContainerRemoveOptions{})
		if err != nil {
			panic(err)
		}
		fmt.Printf("[Container] Name: %s - ID: %s - Status %s removido com sucesso\n",
			container.ID[:10],
			container.Names,
			container.Status,
		)
	}
}

func (d dockerAPI) RemoveImagesDangling() {
	args := filters.NewArgs()
	args.Add("dangling", "true")

	images, err := d.client.ImageList(d.ctx, types.ImageListOptions{
		Filters: args,
	})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		del, err := d.client.ImageRemove(d.ctx, image.ID[7:19], types.ImageRemoveOptions{})
		if err != nil {
			panic(err)
		}

		fmt.Println(del)
		fmt.Printf("[Imagem] ID: %s - Size %s removida com sucesso\n", image.ID[7:19], ByteCountSI(image.Size))
	}
}

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func (d dockerAPI) PrintLogsContainer(containerID string) {
	options := types.ContainerLogsOptions{ShowStdout: true}
	// Replace this ID with a container that really exists
	out, err := d.client.ContainerLogs(d.ctx, containerID, options)
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)
}

func (d dockerAPI) StopAllContainers() {
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

func (d dockerAPI) RunContainerBackground(imageName string) {
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

func (d dockerAPI) PullImage(imageName string) {
	out, err := d.client.ImagePull(d.ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	defer out.Close()

	io.Copy(os.Stdout, out)
}

func (d dockerAPI) RunContainer(imageName string) {
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

func (d dockerAPI) ListContainers() {
	containers, err := d.client.ContainerList(d.ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s %s\n", container.ID[:10], container.Image, container.Names)
	}
}

func (d dockerAPI) ListImages() {
	images, err := d.client.ImageList(d.ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		fmt.Printf("[Imagem] ID: %s - Size %s\n", image.ID[7:19], ByteCountSI(image.Size))
	}
}
