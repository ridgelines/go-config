package config

type Config struct {
	Providers []*Provider
}

func NewConfig(providers ...*Provider) *Config {
	return &Config{
		Providers: providers,
	}
}

func (c *Config) Get(key string) (string, error) {
	return "", nil
}

/*
func Google(query string) (results []Result) {
    c := make(chan Result)
    go func() { c <- Web(query) } ()
    go func() { c <- Image(query) } ()
    go func() { c <- Video(query) } ()

    for i := 0; i < 3; i++ {
        result := <-c
        results = append(results, result)
    }
    return
}
*/
