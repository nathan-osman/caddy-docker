package docker

import (
	"context"
)

// Restart attempts to restart the specified container.
func (m *Monitor) Restart(ctx context.Context, id string) error {
	return m.client.ContainerRestart(ctx, id, nil)
}
