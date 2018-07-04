package config

type Config struct {
	URL string
	Port int
	Mode string
}

func New(url string, port int, mode string) *Config {
	return &Config{URL: url, Port: port, Mode: mode}
}
