package config

type Provider interface {
	Load() (map[string]string, error)
}
