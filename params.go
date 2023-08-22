package cli

import "errors"

type VemtaService struct {
	Name       string
	Repository string
	FolderName string
}

var Services = []VemtaService{
	{
		Name:       "MVC",
		Repository: "https://github.com/vemta/mvc",
		FolderName: "mvc",
	},
	{
		Name:       "API",
		Repository: "https://github.com/vemta/api",
		FolderName: "api",
	},
	{
		Name:       "Payment Gateway",
		Repository: "https://github.com/vemta/payment",
		FolderName: "payment",
	},
}

var MustHaveSoftwares = []string{
	"git", "docker", "docker-compose",
}

func GetServiceByName(name string) (*VemtaService, error) {
	for _, service := range Services {
		if service.Name == name {
			return &service, nil
		}
	}
	return nil, errors.New("Service not found.")
}
