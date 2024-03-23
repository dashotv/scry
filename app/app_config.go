package app

import (
	"strings"

	"github.com/caarlos0/env/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"
)

func init() {
	initializers = append(initializers, setupConfig)
}

func setupConfig(app *Application) error {
	app.Config = &Config{}
	if err := env.Parse(app.Config); err != nil {
		return errors.Wrap(err, "parsing config")
	}
	return nil
}

type Config struct {
	Mode   string `env:"MODE" envDefault:"dev"`
	Logger string `env:"LOGGER" envDefault:"dev"`
	Port   int    `env:"PORT" envDefault:"10080"`
	//golem:template:app/config_partial_struct
	// DO NOT EDIT. This section is managed by github.com/dashotv/golem.
	// Router Auth
	Auth           bool   `env:"AUTH" envDefault:"false"`
	ClerkSecretKey string `env:"CLERK_SECRET_KEY"`

	// Events
	NatsURL string `env:"NATS_URL"`

	//golem:template:app/config_partial_struct

	ElasticsearchURL string `env:"ELASTICSEARCH_URL,required"`
	NzbgeekKey       string `env:"NZBGEEK_KEY,required"`
	NzbgeekURL       string `env:"NZBGEEK_URL,required"`
	TvdbKey          string `env:"TVDB_KEY,required"`
	TmdbToken        string `env:"TMDB_TOKEN,required"`
	Production       bool   `env:"PRODUCTION" envDefault:"false"`
}

type Connection struct {
	URI        string `yaml:"uri,omitempty"`
	Database   string `yaml:"database,omitempty"`
	Collection string `yaml:"collection,omitempty"`
}

func (c *Connection) UnmarshalText(text []byte) error {
	vals := strings.Split(string(text), ",")
	c.URI = vals[0]
	c.Database = vals[1]
	c.Collection = vals[2]
	return nil
}

type ConnectionSet map[string]*Connection

func (c *ConnectionSet) UnmarshalText(text []byte) error {
	*c = make(map[string]*Connection)
	for _, conn := range strings.Split(string(text), ";") {
		kv := strings.Split(conn, "=")
		vals := strings.Split(kv[1]+",,", ",")
		(*c)[kv[0]] = &Connection{
			URI:        vals[0],
			Database:   vals[1],
			Collection: vals[2],
		}
	}
	return nil
}

func (c *Config) Validate() error {
	list := []func() error{
		c.validateMode,
		c.validateLogger,
		//golem:template:app/config_partial_validate
		// DO NOT EDIT. This section is managed by github.com/dashotv/golem.
		//golem:template:app/config_partial_validate

	}

	for _, fn := range list {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) validateMode() error {
	switch c.Mode {
	case "dev", "release":
		return nil
	default:
		return errors.New("invalid mode (must be 'dev' or 'release')")
	}
}

func (c *Config) validateLogger() error {
	switch c.Logger {
	case "dev", "release":
		return nil
	default:
		return errors.New("invalid logger (must be 'dev' or 'release')")
	}
}

//golem:template:app/config_partial_connection
// DO NOT EDIT. This section is managed by github.com/dashotv/golem.

//golem:template:app/config_partial_connection
