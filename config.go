package cli

type ContainerConfig struct {
	ServiceName string `json:"service_name"`
	ContainerID string `json:"container"`
}

type Config struct {
	Containers []ContainerConfig `json:"containers"`
}
