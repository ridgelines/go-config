package config

import (
	"testing"
)

func TestResolver(t *testing.T) {
	static := NewStatic(map[string]string{"items.server.timeout": "30"})

	mappings := map[string]string{
		"items.server.timeout": "server.timeout",
	}

	resolver := NewResolver(static, mappings)

	settings, err := resolver.Load()
	if err != nil {
		t.Error(err)
	}

	if settings["server.timeout"] != "30" {
		t.Errorf("Timeout was '%s', expected '30'", settings["server.timeout"])
	}
}
