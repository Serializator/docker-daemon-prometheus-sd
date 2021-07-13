package main

import (
	"encoding/json"
	"fmt"

	"io.serializator/docker-daemon-prometheus-sd/config"
)

// TODO: read the configuration to determine what File-Based Service Discovery configuration to generate
// TODO: listen for events from Docker to add or remove containers from the File-Based Service Discovery configuration

func main() {
	config, err := config.Read()
	if err != nil {
		panic(err)
	}

	for _, probe := range config.Probes {
		fmt.Printf("Search for containers ... (%v)\n", probe.Name)

		containers, err := probe.List()
		if err != nil {
			fmt.Printf("%v\n", err.Error())
			continue
		}

		data, err := json.Marshal(containers)
		if err != nil {
			fmt.Printf("%v\n", err.Error())
			continue
		}

		fmt.Println(string(data))
	}
}
