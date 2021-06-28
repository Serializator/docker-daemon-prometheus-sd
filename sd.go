package main

import (
	"os"
	"strconv"

	"github.com/docker/docker/api/types"
	"gopkg.in/yaml.v2"
)

type SD struct {
	File     string
	Criteria MatchCriteria
	SDEntry  map[string]SDEntry
}

type SDEntry struct {
	Labels  map[string]string `yaml:"labels"`
	Targets []string          `yaml:"targets"`
}

func (sd *SD) AddOrUpdateEntry(container types.Container) {
	targets := make([]string, 0)

	for _, port := range container.Ports {
		publicPort := strconv.FormatUint(uint64(port.PublicPort), 10)

		// only add it as a target if the port is actually exposed
		if port.IP != "" {
			targets = append(targets, port.IP+":"+publicPort)
		}
	}

	if entry, ok := sd.SDEntry[container.ID]; ok {
		entry.Targets = targets
	} else {
		sd.SDEntry[container.ID] = SDEntry{
			Labels:  make(map[string]string),
			Targets: targets,
		}
	}
}

func (sd *SD) RemoveEntry(container types.Container) {
	delete(sd.SDEntry, container.ID)
}

func (sd *SD) NewWriter() *SDWriter {
	return &SDWriter{sd}
}

type SDWriter struct {
	sd *SD
}

func (sdWriter SDWriter) Write() error {
	ioWriter, err := os.Create(sdWriter.sd.File + ".yaml")
	if err != nil {
		return err
	}
	defer ioWriter.Close()

	// TODO: use the "len(sdWriter.sd.SDEntry)" as pre-defined length for the array
	sdEntries := make([]SDEntry, 0)
	for _, sdEntry := range sdWriter.sd.SDEntry {
		sdEntries = append(sdEntries, sdEntry)
	}

	yamlEncoder := yaml.NewEncoder(ioWriter)
	defer yamlEncoder.Close()
	return yamlEncoder.Encode(sdEntries)
}
