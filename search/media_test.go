package search

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var elasticURL string

func init() {
	godotenv.Load("../.env")
	elasticURL = os.Getenv("ELASTICSEARCH_URL")
}

func TestMediaSearch_Find(t *testing.T) {
	c, err := New(elasticURL)
	require.NoError(t, err)

	s := c.Media.NewSearch()
	s.Type = "series"
	s.Name = "my hero academia"

	r, err := s.Find()
	require.NoError(t, err)

	fmt.Printf("found: %d/%d\n", r.Count, r.Total)
	for _, f := range r.Media {
		fmt.Printf("%10s %s\n", f.Type, f.Name)
	}
}
