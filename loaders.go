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

	this.once.Do(func() {
		this.settings, err = this.provider.Load()
	})

	return this.settings, err
}

type CachedLoader struct {
	invalidated bool
	provider    Provider
	settings    map[string]string
}

func NewCachedLoader(provider Provider) *CachedLoader {
	return &CachedLoader{
		invalidated: true,
		provider:    provider,
		settings:    map[string]string{},
	}
}

func (this *CachedLoader) Load() (map[string]string, error) {
	if this.invalidated {
		settings, err := this.provider.Load()
		if err != nil {
			return nil, err
		}

		this.settings = settings
		this.invalidated = false
	}

	return this.settings, nil
}

func (this *CachedLoader) Invalidate() {
	this.invalidated = true
}
