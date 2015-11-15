package config

import (
	"os"
	"testing"
)

func TestEnvironmentLoad(t *testing.T) {
	os.Setenv("GLOBAL_TIMEOUT", "30")
	os.Setenv("GLOBAL_FREQUENCY", "0.5")
	os.Setenv("LOCAL_TIME_ZONE", "PST")
	os.Setenv("LOCAL_ENABLED", "true")

	mappings := map[string]string{
		"GLOBAL_TIMEOUT":   "global.timeout",
		"GLOBAL_FREQUENCY": "global.frequency",
		"LOCAL_TIME_ZONE":  "local.time_zone",
		"LOCAL_ENABLED":    "local.enabled",
	}

	p := NewEnvironment(mappings)

	expectedSettings := map[string]string{
		"global.timeout":   "30",
		"global.frequency": "0.5",
		"local.time_zone":  "PST",
		"local.enabled":    "true",
	}

	actualSettings, err := p.Load()
	if err != nil {
		t.Error(err)
	}

	for key, expected := range expectedSettings {
		actual, ok := actualSettings[key]

		if !ok {
			t.Errorf("Key '%s' not in settings", key)
		}

		if actual != expected {
			t.Errorf("Setting '%s' was '%s', expected '%s'", key, actual, expected)
		}
	}
}
