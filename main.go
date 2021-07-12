package main

import (
	"fmt"

	"github.com/go-yaml/yaml"
	"io.serializator/docker-daemon-prometheus-sd/config"
)

// TODO: read the configuration to determine what File-Based Service Discovery configuration to generate
// TODO: listen for events from Docker to add or remove containers from the File-Based Service Discovery configuration

func main() {
	config, err := config.Read()
	if err != nil {
		panic(err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", string(data))
}
