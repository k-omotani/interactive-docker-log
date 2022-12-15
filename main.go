package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/manifoldco/promptui"
)

func main() {

	ctx := context.Background()
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	dockerClient.NegotiateAPIVersion(ctx)

	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	if len(containers) == 0 {
		fmt.Print("no containers")
		return
	}

	var containerNames []string
	for _, container := range containers {
		containerNames = append(containerNames, container.Names[0])
	}

	prompt := promptui.Select{
		Label: "Select image",
		Items: containerNames,
	}

	_, result, err := prompt.Run()

	if err != nil {
		panic(err)
	}
	fmt.Printf("You choose %q\n", result)
	logs, err := dockerClient.ContainerLogs(ctx, result, types.ContainerLogsOptions{
		Follow:     true,
		ShowStdout: true,
	})
	if err != nil {
		panic(err)
	}
	defer logs.Close()
	_, err = io.Copy(os.Stdout, logs)
	if err != nil {
		panic(err)
	}
}
