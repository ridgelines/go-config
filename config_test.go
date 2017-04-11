package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigGet(t *testing.T) {
	provider := NewStaticProvider(map[string]string{"key": "val"})
	config := NewConfig(provider)

	val, err := config.Get("key")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "val", val)
}

/*

	Want:
		jsonProvider := config.NewJSONFileProvider("config.json", config.LoadOnce())
		environmentProvider := config.NewEnvironmentProvider(map[string]string{
			"NAME": "name",
			"SIZE": "size",
		})

		config := config.New(jsonProvider, environmentProvider, cliProvider)


*/
