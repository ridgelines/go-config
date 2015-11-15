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

	settings, err := FlatJSON(decodedJSON, "")
	if err != nil {
		return nil, err
	}

	return settings, nil

}

func (this *JSONFile) flatten(inputJSON map[string]interface{}, namespace string) (map[string]string, error) {
	flattened := map[string]string{}

	for key, value := range inputJSON {
		var token string
		if namespace == "" {
			token = key
		} else {
			token = fmt.Sprintf("%s.%s", namespace, key)
		}

		if child, ok := value.(map[string]interface{}); ok {
			settings, err := this.flatten(child, token)
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

//FlatJSON convert the input into map[string]string and adds namespace.
func FlatJSON(input map[string]interface{}, namespace string) (map[string]string, error) {
	out := make(map[string]string)
	for k, v := range input {
		if namespace != "" {
			k = fmt.Sprintf("%s.%s", namespace, k)
		}
		switch v.(type) {
		case map[string]interface{}:
			val, err := FlatJSON(v.(map[string]interface{}), k)
			if err != nil {
				return nil, err
			}

			for cKey, cVal := range val {
				out[cKey] = cVal
			}
		default:
			out[k] = fmt.Sprintf("%v", v)
		}
	}
	return out, nil
}
