package search

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
)

const MEDIA_SEARCH_INDEX = "media" // TODO: Fix
const MEDIA_PAGE_SIZE = 25

type MediaService struct {
	Service
}

func (s *MediaService) NewSearch() *MediaSearch {
	return &MediaSearch{
		client: s.client,
		Search: &Search{Start: 0, Limit: RELEASE_PAGE_SIZE},
	}
}

type Media struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	SearchName  string    `json:"search_name"`
	Display     string    `json:"display"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	ReleaseDate string    `json:"release_date"`
	Background  string    `json:"background"`
	Cover       string    `json:"cover"`
	Created     time.Time `json:"created_at"`
	Updated     time.Time `json:"updated_at"`
}

type MediaSearch struct {
	//ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Display string `json:"display"`
	Title   string `json:"title"`

	client *elastic.Client
	*Search
}

type MediaSearchResponse struct {
	*SearchResponse
	Media []*Media
}

func (s *MediaSearch) Find() (*MediaSearchResponse, error) {
	var q elastic.Query

	r := &MediaSearchResponse{SearchResponse: &SearchResponse{}}
	ctx := context.Background()

	search := s.client.Search().Index(MEDIA_SEARCH_INDEX)
	logrus.Debugf("Find(): start=%d limit=%d", s.Start, s.Limit)
	search = search.From(s.Start)
	search = search.Size(s.Limit)
	search = search.Sort("created_at", false)

	if s.IsZero() {
		q = elastic.NewMatchAllQuery()
		r.Search = "*"
	} else {
		q, r.Search = s.Query()
	}

	search.Query(q)

	sr, err := search.Do(ctx)
	if err != nil {
		logrus.Errorf("Find(): %s", err)
		if e, ok := err.(*elastic.Error); ok {
			logrus.Errorf("Elastic failed with status %d and error %s.", e.Status, e.Details.Reason)
		}
		return r, err
	}

	r.Total = sr.Hits.TotalHits
	r.Count = len(sr.Hits.Hits)

	ms, err := s.processResponse(sr)
	if err != nil {
		return r, err
	}
	r.Media = ms

	return r, nil
}

func (s *MediaSearch) processResponse(res *elastic.SearchResult) ([]*Media, error) {
	var ms []*Media

	if res == nil || res.TotalHits() == 0 {
		return ms, nil
	}

	for _, hit := range res.Hits.Hits {
		m := &Media{}
		if err := json.Unmarshal(*hit.Source, m); err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}

	return ms, nil
}

func (s *MediaSearch) Query() (*elastic.QueryStringQuery, string) {
	list := []string{}

	if s.Name != "" {
		words := strings.Split(s.Name, " ")
		list = append(list, fmt.Sprintf("%s:(%s)", "name", strings.Join(words, " AND ")))
	}

	if s.Display != "" {
		list = append(list, fmt.Sprintf("%s:\"%s\"", "display", s.Display))
	}

	if s.Type != "" {
		list = append(list, fmt.Sprintf("%s:(%s)", "type", s.Type))
	}

	if s.Title != "" {
		list = append(list, fmt.Sprintf("%s:\"%s\"", "title", s.Title))
	}

	str := strings.Join(list, " AND ")
	logrus.Debugf("    search: %s", str)
	return elastic.NewQueryStringQuery(str), str
}

func (s *MediaSearch) IsZero() bool {
	if s.Name != "" {
		return false
	}

	if s.Display != "" {
		return false
	}

	if s.Type != "" {
		return false
	}

	if s.Title != "" {
		return false
	}

	return true
}
