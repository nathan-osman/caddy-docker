package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddytls"
	"github.com/nathan-osman/caddy-docker/configurator"
	"github.com/nathan-osman/caddy-docker/container"
	"github.com/nathan-osman/caddy-docker/docker"
	"github.com/nathan-osman/caddy-docker/server"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "caddy-docker"
	app.Usage = "sync Caddy config with running Docker containers"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "acme-email",
			Usage:  "account email for Let's Encrypt",
			EnvVar: "ACME_EMAIL",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "enable debug output",
			EnvVar: "DEBUG",
		},
		cli.StringFlag{
			Name:   "docker-host",
			Usage:  "Docker engine `URI`",
			EnvVar: "DOCKER_HOST",
			Value:  "unix:///var/run/docker.sock",
		},
		cli.StringFlag{
			Name:   "server-addr",
			Usage:  "address for the HTTP server",
			EnvVar: "SERVER_ADDR",
			Value:  ":8000",
		},
		cli.StringFlag{
			Name:   "server-username",
			Usage:  "username for the HTTP server",
			EnvVar: "SERVER_USERNAME",
		},
		cli.StringFlag{
			Name:   "server-password",
			Usage:  "password for the HTTP server",
			EnvVar: "SERVER_PASSWORD",
		},
	}
	app.Action = func(c *cli.Context) {

		// Set up Caddy
		caddy.AppName = "caddy-docker"
		caddy.AppVersion = "0.1"
		caddytls.Agreed = true
		caddytls.DefaultCAUrl = "https://acme-v01.api.letsencrypt.org/directory"
		caddytls.DefaultEmail = c.String("acme-email")

		// Configure logging
		log := logrus.WithField("context", "main")
		if c.Bool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}

		// Create the configurator
		var (
			events = make(chan *container.Container)
			conf   = configurator.New(&configurator.Config{
				Events: events,
			})
		)
		defer conf.Close()

		// Create the connection to Docker
		docker, err := docker.New(&docker.Config{
			Host:   c.String("docker-host"),
			Events: events,
		})
		if err != nil {
			log.Error(err)
			return
		}
		defer docker.Close()

		// Create the application server
		srv, err := server.New(&server.Config{
			Addr:         c.String("server-addr"),
			Username:     c.String("server-username"),
			Password:     c.String("server-password"),
			Configurator: conf,
			Monitor:      docker,
		})
		if err != nil {
			log.Error(err)
			return
		}
		defer srv.Close()

		// Wait for SIGINT or SIGTERM
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
	}
	app.Run(os.Args)
}
