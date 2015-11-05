package config

type Static struct {
	Tokens map[string]string
}

func NewStatic(tokens map[string]string) *Static {
	return &Static{
		Tokens: tokens,
	}
}

func (this *Static) Load() error {
	return nil
}

func (this *Static) GetTokens() map[string]string {
	return this.Tokens
}
