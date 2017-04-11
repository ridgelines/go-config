package config

import (
	"fmt"
	"github.com/codegangsta/cli"
)

type CLI struct {
	useDefaults bool
	context     *cli.Context
}

func NewCLI(context *cli.Context, useDefaults bool) *CLI {
	return &CLI{
		useDefaults: useDefaults,
		context:     context,
	}
}

func (this *CLI) Load() (map[string]string, error) {
	settings := map[string]string{}

	for _, flag := range this.context.GlobalFlagNames() {
		val := fmt.Sprintf("%v", this.context.GlobalGeneric(flag))

		if this.context.GlobalIsSet(flag) {
			settings[flag] = val
		} else if this.useDefaults && val != "" {
			settings[flag] = val
		}
	}

	for _, flag := range this.context.FlagNames() {
		val := fmt.Sprintf("%v", this.context.Generic(flag))

		if this.context.IsSet(flag) {
			settings[flag] = val
		} else if this.useDefaults && val != "" {
			settings[flag] = val
		}
	}

	return settings, nil
}
