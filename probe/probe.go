package probe

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type Probe struct {
	Name     string   `json:"name"`
	Format   string   `json:"format"`
	Criteria Criteria `json:"criteria"`
}

func (probe Probe) List() ([]types.Container, error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}

	if containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: probe.Criteria.Apply(filters.NewArgs()),
	}); err == nil {
		return containers, nil
	} else {
		return nil, err
	}
}

type Criteria struct {
	Labels map[string]string `json:"labels"`
}

func (criteria Criteria) Match(container types.Container) bool {
	for label, value := range container.Labels {
		if criteria.Labels[label] != value {
			return false
		}
	}

	return true
}

func (criteria Criteria) Apply(filter filters.Args) filters.Args {
	for label, value := range criteria.Labels {
		filter.Add("label", label+"="+value)
	}

	return filter
}
