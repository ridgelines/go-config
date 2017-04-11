package config

type Resolver struct {
	provider Provider
	mappings map[string]string
}

func NewResolver(provider Provider, mappings map[string]string) *Resolver {
	return &Resolver{
		provider: provider,
		mappings: mappings,
	}
}

func (this *Resolver) Load() (map[string]string, error) {
	settings, err := this.provider.Load()
	if err != nil {
		return nil, err
	}

	resolved := map[string]string{}
	for key, val := range settings {
		if dest, ok := this.mappings[key]; ok {
			resolved[dest] = val
		} else {
			resolved[key] = val
		}
	}

	return resolved, nil
}
