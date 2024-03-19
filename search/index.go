package search

import (
	"context"
	"strings"

	"github.com/olivere/elastic/v6"

	runic "github.com/dashotv/runic/app"
)

func (c *Client) IndexMedia(m *Media) (*elastic.IndexResponse, error) {
	m.Type = strings.ToLower(m.Type)
	return c.client.Index().
		Index("media").
		Type("medium").
		Id(m.ID).
		BodyJson(m).
		Do(context.Background())
}

func (c *Client) IndexRelease(r *Release) (*elastic.IndexResponse, error) {
	return c.client.Index().
		Index("torrents").
		Type("_doc").
		Id(r.ID).
		BodyJson(r).
		Do(context.Background())
}

func (c *Client) IndexRunic(r *runic.Release) (*elastic.IndexResponse, error) {
	return c.client.Index().
		Index("runic").
		Type("_doc").
		Id(r.ID.Hex()).
		BodyJson(r).
		Do(context.Background())
}
