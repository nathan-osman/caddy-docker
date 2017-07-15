package container

import (
	"strings"

	"github.com/docker/docker/api/types"
)

const (
	labelAddr    = "caddy.addr"
	labelDomains = "caddy.domains"
)

// Container contains information about a running container. This includes both
// information from Docker (its ID, etc.) and information from labels (domains,
// address, etc.).
type Container struct {
	ID      string
	Name    string
	Domains []string
	Addr    string
	Running bool
}

// New reads information from a container's labels and populates a Container
// instance with the data. If the container does not contain the requires
// labels, nil is returned.
func New(json types.ContainerJSON, running bool) *Container {
	domainRaw, ok := json.Config.Labels[labelDomains]
	if !ok {
		return nil
	}
	domains := []string{}
	for _, d := range strings.Split(domainRaw, ",") {
		domains = append(domains, strings.TrimSpace(d))
	}
	addr, ok := json.Config.Labels[labelAddr]
	if !ok {
		return nil
	}
	return &Container{
		ID:      json.ID,
		Name:    json.Name[1:],
		Domains: domains,
		Addr:    addr,
		Running: running,
	}
}
