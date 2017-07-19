package server

import (
	"github.com/nathan-osman/caddy-docker/configurator"
)

// Config stores the configuration for the HTTP server.
type Config struct {
	Addr         string
	Username     string
	Password     string
	Configurator *configurator.Configurator
}
