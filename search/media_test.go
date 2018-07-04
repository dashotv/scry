package search

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestMediaSearch_Find(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	c, err := New("http://127.0.0.1:9200")
	if err != nil {
		t.Error(err)
	}

	s := c.Media.NewSearch()
	s.Type = "series"
	s.Name = "my hero academia"

	r, err := s.Find()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("found: %d/%d\n", r.Count, r.Total)
	for _, f := range r.Media {
		fmt.Printf("%10s %s\n", f.Type, f.Name)
	}
}
