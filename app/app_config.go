package app

import (
	"github.com/caarlos0/env/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"
)

func setupConfig(app *Application) error {
	app.Config = &Config{}
	if err := env.Parse(app.Config); err != nil {
		return errors.Wrap(err, "parsing config")
	}
	return nil
}

type Config struct {
	Mode             string `env:"MODE" envDefault:"dev"`
	Logger           string `env:"LOGGER" envDefault:"dev"`
	Port             int    `env:"PORT" envDefault:"10080"`
	Cron             bool   `env:"CRON" envDefault:"true"`
	ElasticsearchURL string `env:"ELASTICSEARCH_URL,required"`
	NzbgeekKey       string `env:"NZBGEEK_KEY,required"`
	NzbgeekURL       string `env:"NZBGEEK_URL,required"`
	TvdbKey          string `env:"TVDB_KEY,required"`
	TmdbToken        string `env:"TMDB_TOKEN,required"`
	NatsURL          string `env:"NATS_URL,required"`
}
