package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type Container struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Image    string `json:"image"`
	Launched bool   `json:"-"`
}

type DockerClient struct {
	Client *client.Client
}

var Docker DockerClient

func StopContainer(ctx context.Context, c *Container) error {
	stats, err := Docker.Client.ContainerInspect(ctx, c.Id)
	if err != nil {
		return err
	}

	if stats.State.Dead || stats.State.OOMKilled {
		return errors.New("Container already stopped.")
	}

	timeout := new(int)
	*timeout = 60

	return Docker.Client.ContainerStop(ctx, c.Id, container.StopOptions{Timeout: timeout})
}

func LaunchContainer(ctx context.Context, container *Container) error {
	running, err := IsContainerRunning(ctx, container)

	if err != nil {
		return err
	}

	if running {
		return nil
	}

	if err := Docker.Client.ContainerStart(ctx, container.Id, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

func GetCreatedContainers(ctx context.Context) (*[]Container, error) {
	foundContainers, err := Docker.Client.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: "vemta-"}),
	})

	containers := make([]Container, 0)

	if err != nil {
		return nil, err
	}

	for _, container := range foundContainers {
		containers = append(containers, Container{
			Id:       container.ID,
			Name:     container.Names[0],
			Image:    container.Image,
			Launched: container.State == "running",
		})
	}

	return &containers, nil

}

func IsContainerRunning(ctx context.Context, container *Container) (bool, error) {
	stats, err := Docker.Client.ContainerInspect(ctx, container.Id)
	if err != nil {
		return false, err
	}
	return stats.State.Running || stats.State.Paused, nil
}

func GetContainers() *[]Container {
	cmd := exec.Command("docker", "container ls -a --filter \"name=vemta-\" --format \"{{.ID}} {{.Names}} {{.Image}} {{.Status}}\"")
	return parseContainers(cmd)
}

func GetContainersOfService(service *VemtaService) *[]Container {
	cmdStr := fmt.Sprintf("container ls -a --filter \"name=vemta-%s\" --format \"{{.ID}} {{.Names}} {{.Image}} {{.Status}}\"", service.DockerPrefix)
	cmd := exec.Command("docker", cmdStr)
	return parseContainers(cmd)
}

func parseContainers(cmd *exec.Cmd) *[]Container {
	containers := make([]Container, 0)
	output, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	scanner := bufio.NewScanner(output)

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	for scanner.Scan() {
		line := scanner.Text()
		params := strings.Split(line, " ")

		containers = append(containers, Container{
			Id:       params[0],
			Name:     params[1],
			Image:    params[2],
			Launched: params[3] == "Up",
		})
	}

	if err := cmd.Wait(); err != nil {
		panic(err)
	}

	return &containers
}
