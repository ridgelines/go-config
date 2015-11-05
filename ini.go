package config

import (
	"fmt"
	"github.com/go-ini/ini"
	"strings"
)

type INIFile struct {
	*FileProvider
}

func NewINIFile(path string) *INIFile {
	return &INIFile{
		FileProvider: NewFileProvider(path),
	}
}

func (this *INIFile) Load() error {
	file, err := ini.Load(this.Path)
	if err != nil {
		return err
	}

	for _, section := range file.Sections() {
		for _, key := range section.Keys() {
			token := fmt.Sprintf("%s.%s", section.Name(), key.Name())
			this.tokens[strings.ToLower(token)] = key.String()
		}
	}

	return nil
}

func (this *INIFile) GetTokens() map[string]string {
	return this.tokens
}
