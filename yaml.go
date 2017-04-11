package config

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

func NewYAMLFileProvider(path string) *Provider {
	return &Provider{
		Get: YAMLGetter(path),
	}
}

func YAMLGetter(path string) Getter {
	return func(key string) (string, error) {
		encodedYAML, err := ioutil.ReadFile(path)
		if err != nil {
			return "", err
		}

		encodedJSON, err := yaml.YAMLToJSON(encodedYAML)
		if err != nil {
			return "", err
		}

		decodedJSON := map[string]interface{}{}
		if err := json.Unmarshal(encodedJSON, &decodedJSON); err != nil {
			return "", err
		}

		settings, err := flattenJSON(decodedJSON, "")
		if err != nil {
			return "", err
		}

		return settings[key], nil
	}
}
