package search

import (
	"context"
	"fmt"
	"encoding/json"
	"time"

	"github.com/olivere/elastic"
)

const RELEASE_SEARCH_INDEX = "torrents" // TODO: Fix

type ReleaseService struct {
	Service
}

func (s *ReleaseService) NewSearch() *ReleaseSearch {
	return &ReleaseSearch{client: s.client}
}

type Release struct {
	ID          string
	Name        string
	DisplayName string    `json:"display_name"`
	Raw         string    `json:"raw"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Season      int       `json:"season"`
	Episode     int       `json:"episode"`
	Size        string    `json:"size"`
	Guid        string    `json:"guid"`
	Resolution  string    `json:"resolution"`
	Team        string    `json:"team"`
	Author      string    `json:"author"`
	Verified    bool      `json:"verified"`
	Bluray      bool      `json:"bluray"`
	Uncensored  bool      `json:"uncensored"`
	Checksum    string    `json:"checksum" gorm:"unique_index"`
	Download    string    `json:"download"`
	Source      string    `json:"source"`
	Type        string    `json:"type"`
	Published   time.Time `json:"published"`
}

type ReleaseSearch struct {
	Source string
	Type       string
	Name       string
	Author     string
	Group      string
	Season     int
	Episode    int
	Resolution int
	Verified   bool
	Uncensored bool
	Bluray     bool
	Exact      bool

	client *elastic.Client
}

type ReleaseSearchResponse struct {
	*SearchResponse
	Releases []*Release
}

func (s *ReleaseSearch) Find() (*ReleaseSearchResponse, error) {
	r := &ReleaseSearchResponse{SearchResponse: &SearchResponse{}}
	ctx := context.Background()

	search := s.client.Search().Index(RELEASE_SEARCH_INDEX)
	search = search.From(0)
	search = search.Size(10)
	search.Query(s.Query())

	sr, err := search.Do(ctx)
	if err != nil {
		return r, err
	}

	r.Total = sr.Hits.TotalHits
	r.Count = len(sr.Hits.Hits)

	rels, err := s.processResponse(sr)
	if err != nil {
		return r, err
	}
	r.Releases = rels

	return r, nil
}

func (s *ReleaseSearch) processResponse(res *elastic.SearchResult) ([]*Release, error) {
	if res == nil || res.TotalHits() == 0 {
		return nil, nil
	}

	var rels []*Release
	for _, hit := range res.Hits.Hits {
		rel := &Release{}
		if err := json.Unmarshal(*hit.Source, rel); err != nil {
			return nil, err
		}

		rels = append(rels, rel)
	}

	return rels, nil
}

func (s *ReleaseSearch) Query() *elastic.BoolQuery {
	query := elastic.NewBoolQuery()
	if s.Name != "" {
		if s.Exact {
			query = query.Must(elastic.NewTermQuery("name", "\""+s.Name+"\""))
		} else {
			query = query.Must(elastic.NewTermQuery("name", s.Name))
		}
	}

	if s.Source != "" {
		query = query.Must(elastic.NewTermQuery("source", s.Source))
	}
	if s.Type != "" {
		query = query.Must(elastic.NewTermQuery("type", s.Type))
	}

	if s.Author != "" {
		query = query.Must(elastic.NewTermQuery("author", "\""+s.Author+"\""))
	}
	if s.Group != "" {
		query = query.Must(elastic.NewTermQuery("group", "\""+s.Group+"\""))
	}

	//if s.Season >= 0 {
	//	query = query.Must(elastic.NewTermQuery("season", s.Season))
	//}
	//if s.Episode >= 0 {
	//	query = query.Must(elastic.NewTermQuery("episode", s.Episode))
	//}
	if s.Resolution >= 0 {
		query = query.Must(elastic.NewTermQuery("resolution", s.Resolution))
	}
	if s.Verified {
		query = query.Must(elastic.NewTermQuery("verified", s.Verified))
	}
	if s.Uncensored {
		query = query.Must(elastic.NewTermQuery("uncensored", s.Uncensored))
	}
	if s.Bluray {
		query = query.Must(elastic.NewTermQuery("bluray", s.Bluray))
	}
	fmt.Printf("query: %#v\n", query)
	return query
}
