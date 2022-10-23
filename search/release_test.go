package search

import (
	"fmt"
	"testing"
)

func TestReleaseSearch_Find(t *testing.T) {
	c, err := New(elasticURL)
	if err != nil {
		t.Error(err)
		return
	}

	s := c.Release.NewSearch()
	s.Type = "anime"
	s.Name = "my hero academia"
	s.Exact = true
	//s.Resolution = 720
	//s.Verified = true

	r, err := s.Find()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("found: %d/%d\n", r.Count, r.Total)
	for _, r := range r.Releases {
		fmt.Printf("%5t %5s %s\n", r.Verified, r.Resolution, r.Name)
	}
}
