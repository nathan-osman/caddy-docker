package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
)

// processContainer inspects the specified container and sends on the event
// channel if it includes the appropriate labels.
func (m *Monitor) processContainer(ctx context.Context, id string, running bool) error {
	json, err := m.client.ContainerInspect(ctx, id)
	if err != nil {
		return err
	}
	if c := containerFromJSON(json, running); c != nil {
		m.events <- c
	}
	return nil
}

// processContainers examines all of the containers that are already running in
// the Docker daemon and processes each of them.
func (m *Monitor) processContainers(ctx context.Context) error {
	containers, err := m.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return err
	}
	for _, c := range containers {
		m.log.Debugf("found container %s", c.ID)
		if err := m.processContainer(ctx, c.ID, true); err != nil {
			return err
		}
	}
	return nil
}

// processMessage examines the provided message and acts on it accordingly.
func (m *Monitor) processMessage(ctx context.Context, msg events.Message) error {
	switch msg.Action {
	case "start":
		m.log.Debugf("container %s started", msg.ID)
		return m.processContainer(ctx, msg.ID, true)
	case "die":
		m.log.Debugf("container %s died", msg.ID)
		return m.processContainer(ctx, msg.ID, false)
	}
	return nil
}
