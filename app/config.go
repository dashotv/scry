package app

import (
	"sync"
)

var configOnce sync.Once
var configInstance *Config

type Config struct {
	Mode          string                 `yaml:"mode"`
	Port          int                    `yaml:"port"`
	Cron          bool                   `yaml:"cron"`
	Connections   map[string]*Connection `yaml:"connections"`
	Elasticsearch struct {
		URL string `yaml:"url"`
	} `yaml:"elasticsearch"`
	Nzbgeek struct {
		URL string `yaml:"url"`
		Key string `yaml:"key"`
	} `yaml:"nzbgeek"`
	Tvdb struct {
		URL string `yaml:"url"`
		Key string `yaml:"key"`
	} `yaml:"tvdb"`
	Tmdb struct {
		URL   string `yaml:"url"`
		Token string `yaml:"token"`
	} `yaml:"tmdb"`
}

type Connection struct {
	URI        string `yaml:"uri,omitempty"`
	Database   string `yaml:"database,omitempty"`
	Collection string `yaml:"collection,omitempty"`
}

func ConfigInstance() *Config {
	configOnce.Do(func() {
		configInstance = &Config{}
	})
	return configInstance
}

func (c *Config) Validate() error {
	// No DB
	//if err := c.validateDefaultConnection(); err != nil {
	//	return err
	//}
	// TODO: add more validations?
	return nil
}

// No DB
// func (c *Config) validateDefaultConnection() error {
// 	if len(c.Connections) == 0 {
// 		return errors.New("you must specify a default connection")
// 	}
//
// 	var def *Connection
// 	for n, c := range c.Connections {
// 		if n == "default" || n == "Default" {
// 			def = c
// 			break
// 		}
// 	}
//
// 	if def == nil {
// 		return errors.New("no 'default' found in connections list")
// 	}
// 	if def.Database == "" {
// 		return errors.New("default connection must specify database")
// 	}
// 	if def.URI == "" {
// 		return errors.New("default connection must specify URI")
// 	}
//
// 	return nil
// }
