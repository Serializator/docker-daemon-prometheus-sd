package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"gopkg.in/yaml.v2"
)

var cli *client.Client
var config *Config

func init() {
	var err error
	cli, err = client.NewClientWithOpts()

	if err != nil {
		panic(err)
	}

	config, err = NewConfig("config.yml")

	if err != nil {
		panic(err)
	}
}

func main() {
	finished := make(chan bool)

	for _, output := range config.Output {
		containers, err := getContainersByCriteria(output.Criteria)

		if err != nil {
			// TODO: some kind of other error handling here
			panic(err)
		}

		targets := make([]string, 0)

		for _, container := range containers {
			for _, port := range container.Ports {
				publicPort := strconv.FormatUint(uint64(port.PublicPort), 10)

				// only add it as a target if the port is actually exposed
				if port.IP != "" {
					targets = append(targets, port.IP+":"+publicPort)
				}
			}
		}

		// TODO: detect extension, if not "yaml" or "json" then error
		file, err := os.Create(output.File + ".yaml")

		if err != nil {
			panic(err)
		}

		defer file.Close()

		yamlEncoder := yaml.NewEncoder(file)
		yamlEncoder.Encode([1]Result{{
			Labels:  make(map[string]string),
			Targets: targets,
		}})
	}

	go listen()
	<-finished
}

func getContainersByCriteria(criteria MatchCriteria) ([]types.Container, error) {
	filter := filters.NewArgs()
	criteria.ApplyToFilter(filter)

	return cli.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filter,
	})

}

func listen() {
	filter := filters.NewArgs()
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

				for _, output := range config.Output {
					if !output.Criteria.Match(container) {
						continue
					}

					// TODO: close "file"
					file, err := os.Open(output.File + ".yaml")

					if err != nil {
						panic(err)
					}

					result := make([]Result, 1)
					yamlDecoder := yaml.NewDecoder(file)

					if err := yamlDecoder.Decode(&result); err != nil {
						panic(err)
					}

					for _, port := range container.Ports {
						publicPort := strconv.FormatUint(uint64(port.PublicPort), 10)

						// only add it as a target if the port is actually exposed
						if port.IP != "" {
							result[0].Targets = append(result[0].Targets, port.IP+":"+publicPort)
						}
					}

					// TODO: close "file"
					file, err = os.Create(output.File + ".yaml")

					if err != nil {
						panic(err)
					}

					yamlEncoder := yaml.NewEncoder(file)
					yamlEncoder.Encode(result)
				}
			}

			if msg.Status == "die" {
				// TODO: remove the container from the result
			}
		}
	}
}

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
