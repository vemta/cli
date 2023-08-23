package cli

import (
	"errors"
	"os/exec"
	"path/filepath"
)

var Services = []VemtaService{
	{
		Name:         "MVC",
		Repository:   "https://github.com/vemta/mvc",
		FolderName:   "mvc",
		DockerPrefix: "mvc",
		Containers:   []string{"vemta-mvc-application-1", "vemta-mvc-mysql-1"},
	},
	{
		Name:         "API",
		Repository:   "https://github.com/vemta/api",
		FolderName:   "api",
		DockerPrefix: "api",
	},
	{
		Name:         "Payment Gateway",
		Repository:   "https://github.com/vemta/payment",
		FolderName:   "payment",
		DockerPrefix: "payment",
	},
}

var MustHaveSoftwares = []string{
	"git", "docker", "docker-compose",
}

type VemtaService struct {
	Name         string
	Repository   string
	FolderName   string
	Containers   []string
	DockerPrefix string
}

func (s *VemtaService) Sync(workingDir string) error {
	repoDir := filepath.Join(workingDir, s.FolderName)
	if err := s.getResetCommand(repoDir).Run(); err != nil {
		return err
	}
	if err := s.getPullCommand(repoDir).Run(); err != nil {
		return err
	}
	return nil
}

func (s *VemtaService) Build(workingDir string) error {
	repoDir := filepath.Join(workingDir, s.FolderName)
	if err := s.getDockerBuildComand(repoDir).Run(); err != nil {
		return err
	}
	return nil
}

func (s *VemtaService) getDockerBuildComand(where string) *exec.Cmd {
	cmd := exec.Command("docker-compose", "up", "-d", "--no-start")
	cmd.Dir = where
	return cmd
}

func (s *VemtaService) Clone(workingDir string) error {
	return s.getCloneCommand(workingDir).Run()
}

func (s *VemtaService) getResetCommand(where string) *exec.Cmd {
	return exec.Command("git", "-C", where, "reset", "--hard", "HEAD")
}

func (s *VemtaService) getPullCommand(where string) *exec.Cmd {
	return exec.Command("git", "-C", where, "pull", "oritin", "master")
}

func (s *VemtaService) getCloneCommand(where string) *exec.Cmd {
	return exec.Command("git", "-C", where, "clone", s.Repository)
}

func GetServiceByName(name string) (*VemtaService, error) {
	for _, service := range Services {
		if service.Name == name {
			return &service, nil
		}
	}
	return nil, errors.New("Service not found.")
}
