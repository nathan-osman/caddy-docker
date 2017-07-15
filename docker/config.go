package docker

import (
	"github.com/nathan-osman/caddy-docker/container"
)

// Config stores configuration information for connecting to the Docker daemon.
type Config struct {
	Host   string
	Events chan<- *container.Container
}
