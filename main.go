package main

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

var cli *client.Client
var config *Config
var sds []*SD

func init() {
	var err error
	cli, err = client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	// TODO: consider both "config.yml" and "config.yaml"
	config, err = NewConfig("config.yml")
	if err != nil {
		panic(err)
	}

	sds = make([]*SD, len(config.Output))
	for index, output := range config.Output {
		sds[index] = &SD{
			File:     output.File,
			Criteria: output.Criteria,
			SDEntry:  make(map[string]SDEntry),
		}
	}
}

func main() {
	finished := make(chan bool)

	for _, sd := range sds {
		containers, err := getContainersByCriteria(sd.Criteria)
		if err != nil {
			panic(err)
		}

		for _, container := range containers {
			sd.AddOrUpdateEntry(container)
		}

		sd.NewWriter().Write()
	}

	go listen()
	<-finished
}

func listen() {
	filter := filters.NewArgs()

	// https://docs.docker.com/engine/reference/commandline/events/#object-types
	// TODO: also update the output if a network change occurs
	// TODO: are "die", "stop" and "kill" the same or should we listen for all three?
	// TODO: what does the "update" event do and can it be used to listen for network related updates

	filter.Add("type", "container")
	filter.Add("event", "start")
	filter.Add("event", "die")

	msgChan, errChan := cli.Events(context.Background(), types.EventsOptions{
		Filters: filter,
	})

	for {
		select {
		case err := <-errChan:
			panic(err)
		case msg := <-msgChan:
			if msg.Status == "start" {
				container, err := getContainerById(msg.ID)
				if err != nil {
					panic(err)
				}

				for _, sd := range sds {
					if !sd.Criteria.Match(container) {
						continue
					}

					sd.AddOrUpdateEntry(*container)
					sd.NewWriter().Write()
				}
			}

			if msg.Status == "die" {
				container, err := getContainerById(msg.ID)
				if err != nil {
					panic(err)
				}

				for _, sd := range sds {
					if !sd.Criteria.Match(container) {
						continue
					}

					sd.RemoveEntry(*container)
					sd.NewWriter().Write()
				}
			}
		}
	}
}
