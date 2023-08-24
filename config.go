package cli

var Configuration Config

type Config struct {
	BackendNetwork string `mapstructure:"backend_network"`
}
