package docker

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

// Monitor connects to the Docker daemon and watches for containers being
// started and stopped. The monitor sends on the Event channel when containers
// with the appropriate labels are started and stopped.
type Monitor struct {
	Events <-chan *Container
	events chan<- *Container
	client *client.Client
	log    *logrus.Entry
	stop   context.CancelFunc
	stopCh chan bool
}

// run processes events from the Docker daemon until stopped.
func (m *Monitor) run(ctx context.Context) {
	defer close(m.stopCh)
	defer close(m.events)
	defer m.log.Info("event loop stopped")
	f := filters.NewArgs()
	f.Add("event", "start")
	f.Add("event", "die")
	options := types.EventsOptions{Filters: f}
	for {
		err := func() error {
			m.log.Info("processing existing containers")
			if err := m.processContainers(ctx); err != nil {
				return err
			}
			m.log.Info("starting event loop")
			msgChan, errChan := m.client.Events(ctx, options)
			for {
				select {
				case msg := <-msgChan:
					if err := m.processMessage(ctx, msg); err != nil {
						return err
					}
				case err := <-errChan:
					return err
				}
			}
		}()
		if err == context.Canceled {
			return
		}
		m.log.Error(err)
		m.log.Info("reconnecting in 30 seconds")
		select {
		case <-time.After(30 * time.Second):
		case <-ctx.Done():
			return
		}
	}
}

// New creates a new monitor from the specified config. Use the Close method to
// shut down the monitor.
func New(cfg *Config) (*Monitor, error) {
	c, err := client.NewClient(cfg.Host, "1.24", nil, nil)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	var (
		events = make(chan *Container)
		m      = &Monitor{
			Events: events,
			events: events,
			client: c,
			log:    logrus.WithField("context", "docker"),
			stop:   cancel,
			stopCh: make(chan bool),
		}
	)
	go m.run(ctx)
	return m, nil
}

// Close shuts down the monitor.
func (m *Monitor) Close() {
	m.stop()
	<-m.stopCh
}
