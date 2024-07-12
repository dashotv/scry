package search

import (
	"context"

	"github.com/olivere/elastic/v7"

	runic "github.com/dashotv/runic/client"
)

func (c *Client) DeleteIndex(index string) error {
	_, err := c.client.DeleteIndex(index).Do(context.Background())
	return err
}

func (c *Client) IndexMedia(m *Media) (*elastic.IndexResponse, error) {
	return c.Media.Index(m)
}
func (c *Client) DeleteMedia(id string) error {
	return c.Media.Delete(id)
}

func (c *Client) IndexRunic(r *runic.Release) (*elastic.IndexResponse, error) {
	return c.Runic.Index(r)
}
func (c *Client) DeleteRunic(id string) error {
	return c.Runic.Delete(id)
}
