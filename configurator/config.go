package configurator

import (
	"github.com/nathan-osman/caddy-docker/container"
)

type Config struct {
	Events <-chan *container.Container
}
