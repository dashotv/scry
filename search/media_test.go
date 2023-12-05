package search

import (
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

var elasticURL = os.Getenv("ELASTICSEARCH_URL")

func TestMediaSearch_Find(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	c, err := New(elasticURL)
	if err != nil {
		t.Error(err)
		return
	}

	s := c.Media.NewSearch()
	s.Type = "series"
	s.Name = "my hero academia"

	r, err := s.Find()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("found: %d/%d\n", r.Count, r.Total)
	for _, f := range r.Media {
		fmt.Printf("%10s %s\n", f.Type, f.Name)
	}
}
