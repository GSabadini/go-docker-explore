package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"

	"github.com/GSabadini/go-docker-explore/cmd"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	var docker = newDockerAPI(cli, ctx)

	fmt.Println(docker)

	var cmdListImages = &cobra.Command{
		Use:   "list-images [string to remove]",
		Short: "List images",
		Long:  `list all images`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			docker.ListImages()
		},
	}

	cmd.AddCommand(cmdListImages)
	cmd.Execute()

	//docker.RunContainer()
	//docker.RunContainerBackground()
	//docker.StopContainer()
	//docker.PrintLogsContainer()
	//docker.ListContainers()
	//docker.ListContainers("alpine")
	//docker.ListImages()
	//docker.RemoveImagesDangling()
	//docker.RemoveContainersExited()
}
