package config

type Static struct {
	settings map[string]string
}

func NewStatic(settings map[string]string) *Static {
	return &Static{
		settings: settings,
	}
}

func (this *Static) Load() (map[string]string, error) {
	settings := map[string]string{}

	for key, value := range this.settings {
		settings[key] = value
	}

	return settings, nil
}

func (this *Static) Set(key, val string) {
	this.settings[key] = val
}
