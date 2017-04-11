package config

import (
	"os"
)

type Environment struct {
	mappings map[string]string
}

func NewEnvironment(mappings map[string]string) *Environment {
	return &Environment{
		mappings: mappings,
	}
}

func (this *Environment) Load() (map[string]string, error) {
	settings := map[string]string{}

	for env, key := range this.mappings {
		if val := os.Getenv(env); val != "" {
			settings[key] = val
		}
	}

	return settings, nil
}
