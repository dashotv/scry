package search

import (
	"context"

	"github.com/olivere/elastic/v6"
)

type Client struct {
	client *elastic.Client
	url    string

	Production bool
	Code       int
	Version    string

	Media      *MediaService
	MediaIndex string
	Runic      *RunicService
	RunicIndex string
}

type Service struct {
	client *elastic.Client
	env    string
	index  string
}

type Search struct {
	Start int
	Limit int
	Index string
}

type SearchResponse struct {
	Search string
	Total  int64
	Count  int
}

func New(url string, production bool) (*Client, error) {
	var err error
	c := &Client{url: url, Production: production}
	env := "dev"
	if production {
		env = "prod"
	}

	e, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	info, code, err := e.Ping(url).Do(ctx)
	if err != nil {
		return nil, err
	}
	//logrus.Debugf("Elasticsearch returned with code %d and version %s", code, info.Version.Number)
	c.Code = code
	c.Version = info.Version.Number

	c.client = e
	c.Runic = &RunicService{Service: Service{client: e, env: env, index: "runic"}}
	c.Media = &MediaService{Service: Service{client: e, env: env, index: "media"}}

	return c, nil
}
