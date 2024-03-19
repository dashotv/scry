package search

import (
	"context"

	"github.com/olivere/elastic/v6"
)

type Client struct {
	client *elastic.Client
	url    string

	Code    int
	Version string

	Release *ReleaseService
	Media   *MediaService
	Runic   *RunicService
}

type Service struct {
	client *elastic.Client
}

type Search struct {
	Start int
	Limit int
}

type SearchResponse struct {
	Search string
	Total  int64
	Count  int
}

func New(url string) (*Client, error) {
	ctx := context.Background()
	var err error
	c := &Client{url: url}

	e, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}

	info, code, err := e.Ping(url).Do(ctx)
	if err != nil {
		return nil, err
	}
	//logrus.Debugf("Elasticsearch returned with code %d and version %s", code, info.Version.Number)
	c.Code = code
	c.Version = info.Version.Number

	c.client = e
	c.Release = &ReleaseService{Service: Service{client: e}}
	c.Runic = &RunicService{Service: Service{client: e}}
	c.Media = &MediaService{Service: Service{client: e}}

	return c, nil
}
