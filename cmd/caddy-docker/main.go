package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mholt/caddy"
	"github.com/nathan-osman/caddy-docker/configurator"
	"github.com/nathan-osman/caddy-docker/container"
	"github.com/nathan-osman/caddy-docker/docker"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func init() {
	caddy.AppName = "caddy-docker"
	caddy.AppVersion = "0.1"
}

func main() {
	app := cli.NewApp()
	app.Name = "caddy-docker"
	app.Usage = "sync Caddy config with running Docker containers"
	app.Flags = []cli.Flag{
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
	}
	app.Action = func(c *cli.Context) {
		log := logrus.WithField("context", "main")
		if c.Bool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		var (
			events = make(chan *container.Container)
			conf   = configurator.New(&configurator.Config{
				Events: events,
			})
		)
		defer conf.Close()
		docker, err := docker.New(&docker.Config{
			Host:   "unix:///var/run/docker.sock",
			Events: events,
		})
		if err != nil {
			log.Error(err)
			return
		}
		defer docker.Close()
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
	}
	app.Run(os.Args)
}
