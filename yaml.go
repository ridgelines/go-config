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

func (this *YAMLFile) Load() error {
	encodedYAML, err := ioutil.ReadFile(this.Path)
	if err != nil {
		return err
	}

	encodedJSON, err := yaml.YAMLToJSON(encodedYAML)
	if err != nil {
		return err
	}

	decodedJSON := map[string]interface{}{}
	if err := json.Unmarshal(encodedJSON, &decodedJSON); err != nil {
		return err
	}

	tokens, err := this.flatten(decodedJSON, "")
	if err != nil {
		return err
	}

	this.tokens = tokens
	return nil
}
