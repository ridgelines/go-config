package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func NewJSONFileProvider(path string) *Provider {
	return &Provider{
		Get: JSONGetter(path),
	}
}

func JSONGetter(path string) Getter {
	return func(key string) (string, error) {
		encodedJSON, err := ioutil.ReadFile(path)
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

func flattenJSON(input map[string]interface{}, namespace string) (map[string]string, error) {
	flattened := map[string]string{}

	for key, value := range input {
		token := key
		if namespace != "" {
			token = fmt.Sprintf("%s.%s", namespace, key)
		}

		child, ok := value.(map[string]interface{})
		if !ok {
			flattened[token] = fmt.Sprintf("%v", value)
			continue
		}

		settings, err := flattenJSON(child, token)
		if err != nil {
			return nil, err
		}

		for k, v := range settings {
			flattened[k] = v
		}
	}

	return flattened, nil
}
