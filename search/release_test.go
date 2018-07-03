package search

import (
	"testing"
	"fmt"
)

func TestTorrentSearch_Find(t *testing.T) {
	c, err := New("http://127.0.0.1:9200")
	if err != nil {
		t.Error(err)
	}

	s := c.Release.NewSearch()
	s.Type = "movies"
	s.Name = "notebook"
	s.Resolution = 720
	s.Verified = true

	r, err := s.Find()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("found: %d/%d\n", r.Count, r.Total)
	for _, r := range r.Releases {
		fmt.Printf("%5t %5s %s\n", r.Verified, r.Resolution, r.Name)
	}
}
