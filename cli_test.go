package config

import (
	"testing"

	"github.com/urfave/cli"
)

func TestCLILoad(t *testing.T) {
	var executed bool

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name: "timeout",
		},
		&cli.Float64Flag{
			Name: "frequency",
		},
		&cli.StringFlag{
			Name: "time_zone",
		},
		&cli.BoolFlag{
			Name: "enabled",
		},
	}

	app.Action = func(c *cli.Context) error {
		executed = true
		cliProvider := NewCLI(c, false)

		expectedSettings := map[string]string{
			"timeout":   "30",
			"frequency": "0.5",
			"time_zone": "PST",
			"enabled":   "true",
		}

		actualSettings, err := cliProvider.Load()
		if err != nil {
			t.Error(err)
			return err
		}

		for key, expected := range expectedSettings {
			actual, ok := actualSettings[key]

			if !ok {
				t.Errorf("Key '%s' not in settings", key)
			}

			if actual != expected {
				t.Errorf("Setting '%s' was '%s', expected '%s'", key, actual, expected)
			}
		}
		return err
	}

	app.Run(
		[]string{
			"main.go",
			"--timeout",
			"30",
			"--frequency",
			"0.5",
			"--time_zone",
			"PST",
			"--enabled",
		},
	)

	if !executed {
		t.Fail()
	}
}
