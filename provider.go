package config

type Provider interface {
	Load() error
	GetTokens() map[string]string
}
