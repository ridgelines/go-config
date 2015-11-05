package config

import (
	"os"
)

type Environment struct {
	TokenMap map[string]string
	tokens   map[string]string
}

func NewEnvironment(tokenMap map[string]string) *Environment {
	return &Environment{
		TokenMap: tokenMap,
		tokens:   map[string]string{},
	}
}

func (this *Environment) Load() error {
	this.tokens = map[string]string{}

	for token, env := range this.TokenMap {
		this.tokens[token] = os.Getenv(env)
	}

	return nil
}

func (this *Environment) GetTokens() map[string]string {
	return this.tokens
}
