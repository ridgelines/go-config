package config

import (
	"fmt"
	"strconv"
	"sync"
)

type Config struct {
	sync.RWMutex
	Providers []Provider
	Validate  func(map[string]string) error
	settings  map[string]string
}

func NewConfig(providers []Provider) *Config {
	return &Config{
		Providers: providers,
		settings:  map[string]string{},
	}
}

func (this *Config) Load() error {
	this.Lock()
	defer this.Unlock()
	this.settings = map[string]string{}
	for _, provider := range this.Providers {
		settings, err := provider.Load()
		if err != nil {
			return err
		}

		for key, val := range settings {
			this.settings[key] = val
		}
	}

	if this.Validate != nil {
		if err := this.Validate(this.settings); err != nil {
			return err
		}
	}

	return nil
}

func (this *Config) lookup(key string) (string, error) {
	if err := this.Load(); err != nil {
		return "", err
	}
	this.Lock()
	val, ok := this.settings[key]
	this.Unlock()
	if ok {
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
		return "", fmt.Errorf("Required setting '%s' not set", key)
	}

	return val, nil
}

func (this *Config) StringOr(key, alt string) (string, error) {
	val, err := this.lookup(key)
	if err != nil {
		return "", err
	}

	if val == "" {
		return alt, nil
	}

	return val, nil
}

func (this *Config) Int(key string) (int, error) {
	val, err := this.lookup(key)
	if err != nil {
		return 0, err
	}

	if val == "" {
		return 0, fmt.Errorf("Required setting '%s' not set", key)
	}

	return strconv.Atoi(val)
}

func (this *Config) IntOr(key string, alt int) (int, error) {
	val, err := this.lookup(key)
	if err != nil {
		return 0, err
	}

	if val == "" {
		return alt, nil
	}

	return strconv.Atoi(val)
}

func (this *Config) Float(key string) (float64, error) {
	val, err := this.lookup(key)
	if err != nil {
		return 0, err
	}

	if val == "" {
		return 0, fmt.Errorf("Required setting '%s' not set", key)
	}

	return strconv.ParseFloat(val, 64)
}

func (this *Config) FloatOr(key string, alt float64) (float64, error) {
	val, err := this.lookup(key)
	if err != nil {
		return 0, err
	}

	if val == "" {
		return alt, nil
	}

	return strconv.ParseFloat(val, 64)
}

func (this *Config) Bool(key string) (bool, error) {
	val, err := this.lookup(key)
	if err != nil {
		return false, err
	}

	if val == "" {
		return false, fmt.Errorf("Required setting '%s' not set", key)
	}

	return strconv.ParseBool(val)
}

func (this *Config) BoolOr(key string, alt bool) (bool, error) {
	val, err := this.lookup(key)
	if err != nil {
		return false, err
	}

	if val == "" {
		return alt, nil
	}

	return strconv.ParseBool(val)
}

func (this *Config) Settings() (map[string]string, error) {
	if err := this.Load(); err != nil {
		return nil, err
	}

	return this.settings, nil
}
