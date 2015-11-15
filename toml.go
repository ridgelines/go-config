package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type TOMLFile struct {
	path string
}

func NewTOMLFile(path string) *TOMLFile {
	return &TOMLFile{
		path: path,
	}
}

func (t *TOMLFile) Load() (map[string]string, error) {
	data, err := ioutil.ReadFile(t.path)
	if err != nil {
		return nil, err
	}
	out := make(map[string]interface{})
	_, terr := toml.Decode(string(data), &out)
	if terr != nil {
		return nil, terr
	}
	return FlatJSON(out, "")
}
