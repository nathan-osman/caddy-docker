package configurator

import (
	"github.com/nathan-osman/caddy-docker/container"
)

func (c *Configurator) generate(m map[string]*container.Container) {
	c.log.Info("generating new configuration")
}
