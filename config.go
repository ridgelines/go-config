package config

import (
	"fmt"
	"strconv"
	"sync"
)

type Config struct {
	Providers []Provider
	Loadf     func(*Config) (map[string]string, error)
	Validatef func(map[string]string) error
	tokens    map[string]string
}

func NewConfig(providers []Provider) *Config {
	return &Config{
		Providers: providers,
		Loadf:     LoadDefault,
		tokens:    map[string]string{},
	}
}

func (this *Config) Load() error {
	if this.Loadf != nil {
		tokens, err := this.Loadf(this)
		if err != nil {
			return err
		}

		this.tokens = tokens
	}

	if this.Validatef != nil {
		if err := this.Validatef(this.tokens); err != nil {
			return err
		}
	}

	return nil
}

func (this *Config) lookup(key string) (string, error) {
	if err := this.Load(); err != nil {
		return "", err
	}

	if val, ok := this.tokens[key]; ok {
		return val, nil
	}

	return "", nil
}

func (this *Config) String(key string) (string, error) {
	val, err := this.lookup(key)
	if err != nil {
		return "", err
	}

	if val == "" {
		return "", fmt.Errorf("Required token '%s' not set", key)
	}

	return val, nil
}

func (this *Config) Int(key string) (int, error) {
	val, err := this.lookup(key)
	if err != nil {
		return 0, err
	}

	if val == "" {
		return 0, fmt.Errorf("Required token '%s' not set", key)
	}

	return strconv.Atoi(val)
}

func (this *Config) Bool(key string) (bool, error) {
	val, err := this.lookup(key)
	if err != nil {
		return false, err
	}

	if val == "" {
		return false, fmt.Errorf("Required token '%s' not set", key)
	}

	return strconv.ParseBool(val)
}

func (this *Config) GetTokens() map[string]string {
	return this.tokens
}

func LoadDefault(c *Config) (map[string]string, error) {
	tokens := map[string]string{}

	for _, provider := range c.Providers {
		if err := provider.Load(); err != nil {
			return nil, err
		}

		for key, val := range provider.GetTokens() {
			tokens[key] = val
		}
	}

	return tokens, nil
}

var once sync.Once

func LoadOnce(c *Config) (map[string]string, error) {
	var err error
	var tokens = c.GetTokens()

	once.Do(
		func() {
			tokens, err = LoadDefault(c)
		},
	)

	return tokens, err
}
