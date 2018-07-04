package search

import (
	"context"
	"github.com/olivere/elastic"
)

type Client struct {
	client *elastic.Client
	url    string

	Code    int
	Version string

	Release *ReleaseService
	Media   *MediaService
}

type Service struct {
	client *elastic.Client
}

type Search struct {
	Start int
	Limit int
}

type SearchResponse struct {
	Total int64
	Count int
}

func New(url string) (*Client, error) {
	ctx := context.Background()
	var err error
	c := &Client{url: url}

	e, err := elastic.NewClient()
	if err != nil {
		return nil, err
	}

	info, code, err := e.Ping(url).Do(ctx)
	if err != nil {
		return nil, err
	}
	//logrus.Debugf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	c.Code = code
	c.Version = info.Version.Number

	c.client = e
	c.Release = &ReleaseService{Service: Service{client: e}}
	c.Media = &MediaService{Service: Service{client: e}}

	return c, nil
}
