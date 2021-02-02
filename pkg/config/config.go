package config

import (
	"github.com/urfave/cli/v2"
)

// Logs defines the level and color for log configuration.
type Logs struct {
	Level  string
	Pretty bool
	Color  bool
}

// Driver defines the driver configuration.
type Driver struct {
	Href        string
	Insecure    bool
	Username    string
	Password    string
	Org         string
	Datacenters cli.StringSlice
	Endpoint    string
}

// Kubernetes defines the kubernetes configuration.
type Kubernetes struct {
	Nodename  string
	Namespace string
	PodIP     string
}

// Config is a combination of all available configurations.
type Config struct {
	Logs       Logs
	Driver     Driver
	Kubernetes Kubernetes
}

// Load initializes a default configuration struct.
func Load() *Config {
	return &Config{}
}
