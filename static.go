package config

func NewStaticProvider(settings map[string]string) *Provider {
	return &Provider{
		Get: StaticGetter(settings),
	}
}

func StaticGetter(settings map[string]string) Getter {
	return func(key string) (string, error) {
		return settings[key], nil
	}
}
