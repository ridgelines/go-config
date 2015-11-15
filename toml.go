package config

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

type TOMLFile struct {
	path string
}

func NewTOMLFile(path string) *TOMLFile {
	return &TOMLFile{
		path: path,
	}
}

func (this *TOMLFile) Load() (map[string]string, error) {
	data, err := ioutil.ReadFile(this.path)
	if err != nil {
		return nil, err
	}

	out := make(map[string]interface{})
	if _, err := toml.Decode(string(data), &out); err != nil {
		return nil, err
	}

	return FlattenJSON(out, "")
}
