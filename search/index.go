package search

import (
	"context"
	"strings"

	"github.com/olivere/elastic/v6"

	runic "github.com/dashotv/runic/client"
)

func (c *Client) IndexMedia(m *Media) (*elastic.IndexResponse, error) {
	m.Type = strings.ToLower(m.Type)
	return c.client.Index().
		Index(c.MediaIndex).
		Type("medium").
		Id(m.ID).
		BodyJson(m).
		Do(context.Background())
}
func (c *Client) DeleteMedia(id string) error {
	_, err := c.client.Delete().
		Index(c.MediaIndex).
		Type("medium").
		Id(id).
		Do(context.Background())
	return err
}

func (c *Client) IndexRelease(r *Release) (*elastic.IndexResponse, error) {
	return c.client.Index().
		Index(c.ReleaseIndex).
		Type("_doc").
		Id(r.ID).
		BodyJson(r).
		Do(context.Background())
}

func (c *Client) IndexRunic(r *runic.Release) (*elastic.IndexResponse, error) {
	return c.client.Index().
		Index(c.RunicIndex).
		Type("_doc").
		Id(r.ID.Hex()).
		BodyJson(r).
		Do(context.Background())
}
