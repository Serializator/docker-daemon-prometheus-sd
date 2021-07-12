package config

import (
	"os"

	"github.com/go-yaml/yaml"
	"io.serializator/docker-daemon-prometheus-sd/probe"
)

type Config struct {
	Probes []probe.Probe `json:"probes"`
}

func Read() (*Config, error) {
	file, err := os.Open("config.yaml")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	config := &Config{}
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

func Write(config Config) error {
	file, err := os.Create("config.yaml")
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := yaml.NewEncoder(file)
	if err = encoder.Encode(config); err != nil {
		return err
	}

	return nil
}
