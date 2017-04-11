package config

type Getter func(string) (string, error)
type Provider struct {
	Get Getter
	// Cache, Translators, middleware, etc.
}
