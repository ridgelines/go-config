package config

import (
	"sync"
)

type OnceLoader struct {
	once     sync.Once
	provider Provider
	settings map[string]string
}

func NewOnceLoader(provider Provider) *OnceLoader {
	return &OnceLoader{
		once:     sync.Once{},
		provider: provider,
		settings: map[string]string{},
	}
}

func (this *OnceLoader) Load() (map[string]string, error) {
	var err error

	this.once.Do(
		func() {
			this.settings, err = this.provider.Load()
		},
	)

	settings := map[string]string{}

	for key, value := range this.settings {
		settings[key] = value
	}

	return settings, err
}
