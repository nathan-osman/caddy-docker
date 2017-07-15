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
			c.generate(m)
			t = nil
		}
	}
}

// New creates a new configurator from the specified configuration.
func New(cfg *Config) (*Configurator, error) {
	cdyfile, err := caddy.LoadCaddyfile("http")
	if err != nil {
		return nil, err
	}
	inst, err := caddy.Start(cdyfile)
	if err != nil {
		return nil, err
	}
	c := &Configurator{
		inst:   inst,
		log:    logrus.WithField("context", "configurator"),
		stopCh: make(chan bool),
	}
	go c.run(cfg.Events)
	return c, nil
}

// Close shuts down the configurator.
func (c *Configurator) Close() {
	func() {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		c.inst.Stop()
	}()
	<-c.stopCh
}
