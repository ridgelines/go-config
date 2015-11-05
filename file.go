package config

type FileProvider struct {
	Path   string
	tokens map[string]string
}

func NewFileProvider(path string) *FileProvider {
	return &FileProvider{
		Path:   path,
		tokens: map[string]string{},
	}
}
