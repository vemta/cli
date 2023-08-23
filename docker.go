package cli

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"

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

func LaunchContainer(ctx context.Context) []*Container {
	return nil
}

func GetContainers() *[]Container {
	cmd := exec.Command("docker", "container ls -a --filter \"name=vemta-\" --format \"{{.ID}} {{.Names}} {{.Image}} {{.Status}}\"")
	return parseContainers(cmd)
}

func GetContainersOfService(service *VemtaService) *[]Container {
	cmdStr := fmt.Sprintf("container ls -a --filter \"name=vemta-%s\" --format \"{{.ID}} {{.Names}} {{.Image}} {{.Status}}\"", service.)
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
