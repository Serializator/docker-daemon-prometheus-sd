package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func getContainerById(id string) (*types.Container, error) {
	filter := filters.NewArgs()
	filter.Add("id", id)

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filter,
	})

	if err != nil {
		return nil, err
	}

	if len(containers) < 1 {
		return nil, fmt.Errorf("unable to find container by ID \"%s\"", id)
	}

	return &containers[0], nil
}

func getContainersByCriteria(criteria MatchCriteria) ([]types.Container, error) {
	filter := filters.NewArgs()
	criteria.ApplyToFilter(filter)

	return cli.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filter,
	})
}
