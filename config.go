package cli

var Configuration Config

type Config struct {
	BackendNetwork string          `json:"backend_network"`
	Services       []ServiceConfig `json:"services"`
}
type ServiceConfig struct {
	ServiceName string      `json:"service_name"`
	Containers  []Container `json:"containers"`
}
