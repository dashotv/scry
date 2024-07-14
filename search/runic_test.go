package search

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRunicSearch_Find(t *testing.T) {
	l, err := zap.NewDevelopment()
	require.NoError(t, err)

	c, err := New(elasticURLs, l.Sugar(), false)
	require.NoError(t, err)

	s := c.Runic.NewSearch()
	s.Type = "tv"
	s.Title = "cowboy bebop"
	//s.Resolution = 720
	//s.Verified = true

	r, err := s.Find()
	require.NoError(t, err)

	fmt.Printf("found: %d/%d\n", r.Count, r.Total)
	for _, r := range r.Releases {
		fmt.Printf("%5t %5s %s %02dx%02d\n", r.Verified, r.Resolution, r.Title, r.Season, r.Episode)
	}
}
