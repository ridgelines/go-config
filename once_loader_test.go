package config

import (
	"testing"
)

func TestOnceLoader(t *testing.T) {
	static := NewStatic(map[string]string{"enabled": "true"})
	loader := NewOnceLoader(static)

	settings, err := loader.Load()
	if err != nil {
		t.Error(err)
	}

	if settings["enabled"] != "true" {
		t.Errorf("Enabled was '%s', expected 'true'", settings["enabled"])
	}

	static.Set("enabled", "false")

	settings, err = loader.Load()
	if err != nil {
		t.Error(err)
	}

	if settings["enabled"] != "true" {
		t.Errorf("Enabled was '%s', expected 'true'", settings["enabled"])
	}
}
