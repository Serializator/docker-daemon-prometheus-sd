package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"gopkg.in/yaml.v2"
)

type Config struct {
	// TODO: add a configuration to configure what ports to check for (e.g 9200 and not 9300 in case of Elasticsearch)
	Output []struct {
		File     string        `yaml:"file"`
		Criteria MatchCriteria `yaml:"criteria"`
	} `yaml:"output"`
}

type MatchCriteria struct {
	Labels map[string]string `yaml:"labels"`
}

func (criteria MatchCriteria) Match(container *types.Container) bool {
	for label, value := range criteria.Labels {
		fmt.Printf("%s: %s\n", label, value)

		bytes, err := json.Marshal(container.Labels)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(bytes))

		if container.Labels[label] != value {
			return false
		}
	}

	return true
}

func (criteria MatchCriteria) ApplyToFilter(filter filters.Args) filters.Args {
	for label, value := range criteria.Labels {
		filter.Add("label", label+"="+value)
	}

	return filter
}

func NewConfig(path string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	yamlDecoder := yaml.NewDecoder(file)

	if err := yamlDecoder.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
