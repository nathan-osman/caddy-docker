package configurator

import (
	"github.com/nathan-osman/caddy-docker/container"
)

// Config stores Caddy configuration.
type Config struct {
	Events <-chan *container.Container
}
