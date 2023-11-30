package search

import (
	"context"

	"github.com/olivere/elastic"
)

func (c *Client) IndexMedia(m *Media) (*elastic.IndexResponse, error) {
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
		Type("torrent").
		Id(r.ID).
		BodyJson(r).
		Do(context.Background())
}
