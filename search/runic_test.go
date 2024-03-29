package search

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunicSearch_Find(t *testing.T) {
	c, err := New(elasticURL, false)
	require.NoError(t, err)

	s := c.Release.NewSearch()
	s.Type = "tv"
	s.Name = "cowboy bebop"
	//s.Resolution = 720
	//s.Verified = true

	r, err := s.Find()
	require.NoError(t, err)

	fmt.Printf("found: %d/%d\n", r.Count, r.Total)
	for _, r := range r.Releases {
		fmt.Printf("%5t %5s %s %02dx%02d\n", r.Verified, r.Resolution, r.Name, r.Season, r.Episode)
	}
}
