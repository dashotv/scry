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

	Release      *ReleaseService
	ReleaseIndex string
	Media        *MediaService
	MediaIndex   string
	Runic        *RunicService
	RunicIndex   string
}

type Service struct {
	client *elastic.Client
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
	a := ""
	if !production {
		a = "_development"
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
	c.Release = &ReleaseService{Service: Service{client: e, index: "torrents" + a}}
	c.ReleaseIndex = "torrents" + a
	c.Runic = &RunicService{Service: Service{client: e, index: "runic" + a}}
	c.RunicIndex = "runic" + a
	c.Media = &MediaService{Service: Service{client: e, index: "media" + a}}
	c.MediaIndex = "media" + a

	return c, nil
}
