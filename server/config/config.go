package config

type Config struct {
	Port          int
	Mode          string
	Debug         bool
	Elasticsearch struct {
		URL string
	}
	Nzbgeek struct {
		URL string
		Key string
	}
}
