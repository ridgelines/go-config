package config

type CachedLoader struct {
	invalidated bool
	provider    Provider
	settings    map[string]string
}

func NewCachedLoader(provider Provider) *CachedLoader {
	return &CachedLoader{
		invalidated: true,
		provider:    provider,
		settings:    map[string]string{},
	}
}

func (this *CachedLoader) Load() (map[string]string, error) {
	if this.invalidated {
		settings, err := this.provider.Load()
		if err != nil {
			return nil, err
		}

		this.settings = settings
		this.invalidated = false
	}

	settings := map[string]string{}

	for key, value := range this.settings {
		settings[key] = value
	}

	return settings, nil
}

func (this *CachedLoader) Invalidate() {
	this.invalidated = true
}
