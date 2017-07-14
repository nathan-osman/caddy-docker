package docker

import (
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

// Monitor connects to the Docker daemon and watches for containers being
// started and stopped.
type Monitor struct {
	client *client.Client
	log    *logrus.Entry
}

// New creates a new monitor from the specified config.
func New(cfg *Config) (*Monitor, error) {
	c, err := client.NewClient(cfg.Host, "1.24", nil, nil)
	if err != nil {
		return nil, err
	}
	m := &Monitor{
		client: c,
		log:    logrus.WithField("context", "docker"),
	}
	return m, nil
}
