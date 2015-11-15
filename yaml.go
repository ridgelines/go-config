package config

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

type YAMLFile struct {
	path string
}

func NewYAMLFile(path string) *YAMLFile {
	return &YAMLFile{
		path: path,
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

	return FlattenJSON(decodedJSON, "")
}
