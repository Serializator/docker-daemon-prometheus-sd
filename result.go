package main

type Result struct {
	Labels  map[string]string `yaml:"labels"`
	Targets []string          `yaml:"targets"`
}
