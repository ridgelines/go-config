package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type JSONFile struct {
	path string
}

func NewJSONFile(path string) *JSONFile {
	return &JSONFile{
		path: path,
	}
}

func (this *JSONFile) Load() (map[string]string, error) {
	encodedJSON, err := ioutil.ReadFile(this.path)
	if err != nil {
		return nil, err
	}

	decodedJSON := map[string]interface{}{}
	if err := json.Unmarshal(encodedJSON, &decodedJSON); err != nil {
		return nil, err
	}

	settings, err := FlattenJSON(decodedJSON, "")
	if err != nil {
		return nil, err
	}

	return settings, nil

}

func FlattenJSON(input map[string]interface{}, namespace string) (map[string]string, error) {
	flattened := map[string]string{}

	for key, value := range input {
		var token string
		if namespace == "" {
			token = key
		} else {
			token = fmt.Sprintf("%s.%s", namespace, key)
		}

		if child, ok := value.(map[string]interface{}); ok {
			settings, err := FlattenJSON(child, token)
			if err != nil {
				return nil, err
			}

			for k, v := range settings {
				flattened[k] = v
			}
		} else {
			flattened[token] = fmt.Sprintf("%v", value)
		}
	}

	return flattened, nil
}
