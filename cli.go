package config

import (
	"github.com/urfave/cli"
)

func NewCLIProvider(c *cli.Context) *Provider {
	return &Provider{
		Get: CLIGetter(c),
	}
}

func CLIGetter(c *cli.Context) Getter {
	return func(key string) (string, error) {
		return c.String(key), nil
	}
}
