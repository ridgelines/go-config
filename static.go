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
	return this.settings, nil
}
