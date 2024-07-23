package search

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/index"

	runic "github.com/dashotv/runic/client"
)

func (c *Client) DeleteIndex(index string) error {
	_, err := c.client.Indices.Delete(index).Do(context.Background())
	return err
}

func (c *Client) IndexMedia(m *Media) (*index.Response, error) {
	return c.Media.Index(m)
}
func (c *Client) DeleteMedia(id string) error {
	return c.Media.Delete(id)
}

func (c *Client) IndexRunic(r *runic.Release) (*index.Response, error) {
	return c.Runic.Index(r)
}
func (c *Client) DeleteRunic(id string) error {
	return c.Runic.Delete(id)
}
