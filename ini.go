package config

import (
	"fmt"
	"github.com/go-ini/ini"
)

func NewINIFileProvider(path string) *Provider {
	return &Provider{
		Get: INIGetter(path),
	}
}

func INIGetter(path string) Getter {
	return func(key string) (string, error) {
		decodedINI, err := ini.Load(path)
		if err != nil {
			return "", err
		}

		settings := map[string]string{}
		for _, section := range decodedINI.Sections() {
			for _, key := range section.Keys() {
				token := fmt.Sprintf("%s.%s", section.Name(), key.Name())
				settings[token] = key.String()
			}
		}

		return settings[key], nil
	}
}
