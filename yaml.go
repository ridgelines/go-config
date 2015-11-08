package config

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

type YAMLFile struct {
	*JSONFile
}

func NewYAMLFile(path string) *YAMLFile {
	return &YAMLFile{
		JSONFile: NewJSONFile(path),
	}
}

func (this *YAMLFile) Load() (map[string]string, error) {
	encodedYAML, err := ioutil.ReadFile(this.path)
	if err != nil {
		return nil, err
	}

	encodedJSON, err := yaml.YAMLToJSON(encodedYAML)
	if err != nil {
		return nil, err
	}

	decodedJSON := map[string]interface{}{}
	if err := json.Unmarshal(encodedJSON, &decodedJSON); err != nil {
		return nil, err
	}

	settings, err := this.flatten(decodedJSON, "")
	if err != nil {
		return nil, err
	}

	return settings, nil
}
