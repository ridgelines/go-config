package config

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

func NewTOMLFileProvider(path string) *Provider {
	return &Provider{
		Get: TOMLGetter(path),
	}
}

func TOMLGetter(path string) Getter {
	return func(key string) (string, error) {
		encodedTOML, err := ioutil.ReadFile(path)
		if err != nil {
			return "", err
		}

		decodedTOML := map[string]interface{}{}
		if _, err := toml.Decode(string(encodedTOML), &decodedTOML); err != nil {
			return "", err
		}

		settings, err := flattenJSON(decodedTOML, "")
		if err != nil {
			return "", err
		}

		return settings[key], nil
	}
}
