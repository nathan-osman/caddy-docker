package configurator

import (
	"sync"
	"time"

	"github.com/mholt/caddy"
	"github.com/nathan-osman/caddy-docker/container"
	"github.com/sirupsen/logrus"

	_ "github.com/mholt/caddy/caddyhttp"
)

// Configurator generates Caddy configurations from running containers and
// runs the server.
type Configurator struct {
	mutex      sync.Mutex
	inst       *caddy.Instance
	containers map[string]*container.Container
	log        *logrus.Entry
	stopCh     chan bool
}

// run receives container events and processes them.
func (c *Configurator) run(events <-chan *container.Container) {
	defer close(c.stopCh)
	defer c.log.Info("configurator stopped")
	c.log.Info("starting configurator")
	var t <-chan time.Time
	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			func() {
				c.mutex.Lock()
				defer c.mutex.Unlock()
				if e.Running {
					c.containers[e.ID] = e
				} else {
					delete(c.containers, e.ID)
				}
			}()
			t = time.After(2 * time.Second)
		case <-t:
			if err := c.generate(); err != nil {
				c.log.Error(err)
			}
			t = nil
		}
	}
}

// New creates a new configurator from the specified configuration.
func New(cfg *Config) *Configurator {
	c := &Configurator{
		log:        logrus.WithField("context", "configurator"),
		containers: make(map[string]*container.Container),
		stopCh:     make(chan bool),
	}
	go c.run(cfg.Events)
	return c
}

// Containers returns a list of containers.
func (c *Configurator) Containers() []*container.Container {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var (
		i = 0
		l = make([]*container.Container, len(c.containers))
	)
	for _, v := range c.containers {
		l[i] = v
		i += 1
	}
	return l
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
