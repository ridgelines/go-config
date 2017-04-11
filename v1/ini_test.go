package config

import (
	"testing"
)

func TestINILoad(t *testing.T) {
	p := NewINIFile("test/config.ini")

	actualSettings, err := p.Load()
	if err != nil {
		t.Error(err)
	}

	expectedSettings := map[string]string{
		"global.timeout":   "30",
		"global.frequency": "0.5",
		"local.time_zone":  "PST",
		"local.enabled":    "true",
		"ABC.def":          "ghi",
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
