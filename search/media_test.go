package search

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var elasticURLs []string

func init() {
	godotenv.Load("../.env")
	elasticURLs = strings.Split(os.Getenv("ELASTICSEARCH_URL"), ",")
}

func TestMediaSearch_Find(t *testing.T) {
	l, err := zap.NewDevelopment()
	require.NoError(t, err)

	c, err := New(elasticURLs, l.Sugar(), false)
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
