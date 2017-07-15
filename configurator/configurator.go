package configurator

import (
	"sync"
	"time"

	"github.com/mholt/caddy"
	"github.com/nathan-osman/caddy-docker/container"
	"github.com/sirupsen/logrus"

	// TODO: is this still needed?
	_ "github.com/mholt/caddy/caddyhttp"
)

// Configurator generates Caddy configurations from running containers and
// runs the server.
type Configurator struct {
	mutex  sync.Mutex
	inst   *caddy.Instance
	log    *logrus.Entry
	stopCh chan bool
}

// run receives container events and processes them.
func (c *Configurator) run(events <-chan *container.Container) {
	defer close(c.stopCh)
	defer c.log.Info("configurator stopped")
	c.log.Info("starting configurator")
	var (
		t <-chan time.Time
		m = make(map[string]*container.Container)
	)
	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			if e.Running {
				m[e.ID] = e
			} else {
				delete(m, e.ID)
			}
			t = time.After(2 * time.Second)
		case <-t:
			if err := c.generate(m); err != nil {
				c.log.Error(err)
			}
			t = nil
		}
	}
}

// New creates a new configurator from the specified configuration.
func New(cfg *Config) *Configurator {
	c := &Configurator{
		log:    logrus.WithField("context", "configurator"),
		stopCh: make(chan bool),
	}
	go c.run(cfg.Events)
	return c
}

// Close shuts down the configurator.
func (c *Configurator) Close() {
	<-c.stopCh
	func() {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		if c.inst != nil {
			c.log.Info("stopping server")
			c.inst.Stop()
		}
	}()
}
