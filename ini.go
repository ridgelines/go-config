package config

import (
	"fmt"
	"github.com/go-ini/ini"
)

type INIFile struct {
	path string
}

func NewINIFile(path string) *INIFile {
	return &INIFile{
		path: path,
	}
}

func (this *INIFile) Load() (map[string]string, error) {
	settings := map[string]string{}

	file, err := ini.Load(this.path)
	if err != nil {
		return nil, err
	}

	for _, section := range file.Sections() {
		for _, key := range section.Keys() {
			token := fmt.Sprintf("%s.%s", section.Name(), key.Name())
			settings[token] = key.String()
		}
	}

	return settings, nil
}
