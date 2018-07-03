package config

type Config struct {
	URL string
}

func New(url string) *Config {
	return &Config{URL: url}
}
