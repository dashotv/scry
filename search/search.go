package search

import (
	"context"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
)

type Client struct {
	client *elasticsearch.TypedClient
	url    []string

	Production bool
	Code       int
	Version    string

	Media      *MediaService
	MediaIndex string
	Runic      *RunicService
	RunicIndex string
}

type Service struct {
	client *elasticsearch.TypedClient
	env    string
	index  string
	log    *zap.SugaredLogger
}

type Search struct {
	Start int
	Limit int
	Index string
	log   *zap.SugaredLogger
}

type SearchResponse struct {
	Search string
	Total  int64
	Count  int
}

func New(urls []string, log *zap.SugaredLogger, production bool) (*Client, error) {
	var err error
	c := &Client{url: urls, Production: production}
	env := "dev"
	if production {
		env = "prod"
	}

	e, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: urls,
	})
	if err != nil {
		return nil, err
	}
	ok, err := e.Ping().Do(context.Background())
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("Elasticsearch is not available")
	}

	c.client = e
	c.Runic = &RunicService{Service: Service{client: e, env: env, index: "runic", log: log.Named("runic")}}
	c.Media = &MediaService{Service: Service{client: e, env: env, index: "media", log: log.Named("media")}}

	return c, nil
}
