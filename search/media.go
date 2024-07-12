package search

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/olivere/elastic/v6"
	"github.com/sirupsen/logrus"
)

const MEDIA_PAGE_SIZE = 25

type MediaService struct {
	Service
}

func (s *MediaService) Index(m *Media) (*elastic.IndexResponse, error) {
	m.Type = strings.ToLower(m.Type)
	return s.client.Index().
		Index(s.index + "_" + s.env + "_*").
		Type("medium").
		Id(m.ID).
		BodyJson(m).
		Do(context.Background())
}
func (s *MediaService) Delete(id string) error {
	_, err := s.client.Delete().
		Index(s.index + "_" + s.env + "_*").
		Type("medium").
		Id(id).
		Do(context.Background())
	return err
}

func (s *MediaService) NewSearch() *MediaSearch {
	return &MediaSearch{
		client: s.client,
		Search: &Search{Start: 0, Limit: MEDIA_PAGE_SIZE, Index: s.index},
	}
}

type Media struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Kind        string    `json:"kind"`
	Source      string    `json:"source"`
	SourceID    string    `json:"source_id"`
	SearchName  string    `json:"search_name"`
	Display     string    `json:"display"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	ReleaseDate string    `json:"release_date"`
	Background  string    `json:"background"`
	Cover       string    `json:"cover"`
	Season      int       `json:"season_number"`
	Episode     int       `json:"episode_number"`
	Absolute    int       `json:"absolute_number"`
	Skipped     bool      `json:"skipped"`
	Downloaded  bool      `json:"downloaded"`
	Completed   bool      `json:"completed"`
	Created     time.Time `json:"created_at"`
	Updated     time.Time `json:"updated_at"`
}

type MediaSearch struct {
	//ID      string `json:"id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	Display    string `json:"display"`
	Title      string `json:"title"`
	Source     string `json:"source"`
	SourceID   string `json:"source_id"`
	Season     int    `json:"season"`
	Episode    int    `json:"episode"`
	Absolute   int    `json:"absolute"`
	Skipped    bool   `json:"skipped"`
	Downloaded bool   `json:"downloaded"`
	Completed  bool   `json:"completed"`

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

	search := s.client.Search().Index(s.Index)
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

	if s.Source != "" {
		list = append(list, fmt.Sprintf("%s:\"%s\"", "source", s.Source))
	}

	if s.SourceID != "" {
		list = append(list, fmt.Sprintf("%s:\"%s\"", "source_id", s.SourceID))
	}

	if s.Season > 0 {
		list = append(list, fmt.Sprintf("%s:%d", "season_number", s.Season))
	}
	if s.Episode > 0 {
		list = append(list, fmt.Sprintf("%s:%d", "episode_number", s.Episode))
	}
	if s.Absolute > 0 {
		list = append(list, fmt.Sprintf("%s:%d", "absolute_number", s.Absolute))
	}
	if s.Skipped {
		list = append(list, fmt.Sprintf("%s:%t", "skipped", s.Skipped))
	}
	if s.Downloaded {
		list = append(list, fmt.Sprintf("%s:%t", "downloaded", s.Downloaded))
	}
	if s.Completed {
		list = append(list, fmt.Sprintf("%s:%t", "completed", s.Completed))
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
