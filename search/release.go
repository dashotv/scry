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

const RELEASE_SEARCH_INDEX = "torrents" // TODO: Fix
const RELEASE_PAGE_SIZE = 25

type ReleaseService struct {
	Service
}

func (s *ReleaseService) NewSearch() *ReleaseSearch {
	return &ReleaseSearch{
		client:     s.client,
		Season:     -1,
		Episode:    -1,
		Resolution: -1,
		Search:     &Search{Start: 0, Limit: RELEASE_PAGE_SIZE},
	}
}

type Release struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Year        int       `json:"year"`
	Display     string    `json:"display"`
	Raw         string    `json:"raw"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Season      int       `json:"season"`
	Episode     int       `json:"episode"`
	Volume      int       `json:"volume"`
	Size        string    `json:"size"`
	Encoding    string    `json:"encoding"`
	Resolution  string    `json:"resolution"`
	Quality     string    `json:"quality"`
	Group       string    `json:"group"`
	Author      string    `json:"author"`
	GroupAuthor string    `json:"groupauthor"`
	Verified    bool      `json:"verified"`
	Bluray      bool      `json:"bluray"`
	NZB         bool      `json:"nzb"`
	Uncensored  bool      `json:"uncensored"`
	Checksum    string    `json:"checksum"`
	View        string    `json:"view"`
	Download    string    `json:"download"`
	Source      string    `json:"source"`
	Type        string    `json:"type"`
	Created     time.Time `json:"created_at"`
	Published   time.Time `json:"published_at"`
}

type ReleaseSearch struct {
	Source     string `json:"source"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	Year       int    `json:"year"`
	Author     string `json:"author"`
	Group      string `json:"group"`
	Season     int    `json:"season"`
	Episode    int    `json:"episode"`
	Resolution int    `json:"resolution"`
	Verified   bool   `json:"verified"`
	Uncensored bool   `json:"uncensored"`
	Bluray     bool   `json:"bluray"`
	NZB        bool   `json:"nzb"`
	Exact      bool

	client *elastic.Client
	*Search
}

type ReleaseSearchResponse struct {
	*SearchResponse
	Releases []*Release
}

func (s *ReleaseSearch) Find() (*ReleaseSearchResponse, error) {
	var q elastic.Query

	r := &ReleaseSearchResponse{SearchResponse: &SearchResponse{}}
	ctx := context.Background()

	search := s.client.Search().Index(RELEASE_SEARCH_INDEX)
	logrus.Debugf("Find(): start=%d limit=%d", s.Start, s.Limit)
	search = search.From(s.Start)
	search = search.Size(s.Limit)
	search = search.Sort("published_at", false)

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

	rels, err := s.processResponse(sr)
	if err != nil {
		return r, err
	}
	r.Releases = rels

	return r, nil
}

func (s *ReleaseSearch) processResponse(res *elastic.SearchResult) ([]*Release, error) {
	var rels []*Release

	if res == nil || res.TotalHits() == 0 {
		return rels, nil
	}

	for _, hit := range res.Hits.Hits {
		rel := &Release{}
		if err := json.Unmarshal(*hit.Source, rel); err != nil {
			return nil, err
		}

		rels = append(rels, rel)
	}

	return rels, nil
}

func (s *ReleaseSearch) Query() (*elastic.QueryStringQuery, string) {
	list := []string{}

	if s.Name != "" {
		if s.Exact {
			list = append(list, fmt.Sprintf("%s:\"%s\"", "name", s.Name))
		} else {
			words := strings.Split(s.Name, " ")
			list = append(list, fmt.Sprintf("%s:(%s)", "name", strings.Join(words, " AND ")))
		}
	}

	if s.Year > 0 {
		list = append(list, fmt.Sprintf("%s:\"%d\"", "year", s.Year))
	}

	if s.Source != "" {
		list = append(list, fmt.Sprintf("%s:\"%s\"", "source", s.Source))
	}
	if s.Type != "" {
		list = append(list, fmt.Sprintf("%s:\"%s\"", "type", s.Type))
	}

	if s.Author != "" {
		list = append(list, fmt.Sprintf("%s:\"%s\"", "author", s.Author))
	}
	if s.Group != "" {
		list = append(list, fmt.Sprintf("%s:\"%s\"", "group", s.Group))
	}

	if s.Season >= 0 {
		list = append(list, fmt.Sprintf("%s:%d", "season", s.Season))
	}
	if s.Episode >= 0 {
		list = append(list, fmt.Sprintf("%s:%d", "episode", s.Episode))
	}

	if s.Resolution >= 0 {
		list = append(list, fmt.Sprintf("%s:%d", "resolution", s.Resolution))
	}

	if s.Verified {
		list = append(list, fmt.Sprintf("%s:%t", "verified", s.Verified))
	}
	if s.Uncensored {
		list = append(list, fmt.Sprintf("%s:%t", "uncensored", s.Uncensored))
	}
	if s.Bluray {
		list = append(list, fmt.Sprintf("%s:%t", "bluray", s.Bluray))
	}
	if s.NZB {
		list = append(list, fmt.Sprintf("%s:%t", "nzb", s.NZB))
	}

	str := strings.Join(list, " AND ")
	logrus.Debugf("    search: %s", str)
	return elastic.NewQueryStringQuery(str), str
}

func (s *ReleaseSearch) IsZero() bool {
	if s.Name != "" {
		return false
	}

	if s.Source != "" {
		return false
	}
	if s.Type != "" {
		return false
	}

	if s.Author != "" {
		return false
	}
	if s.Group != "" {
		return false
	}

	if s.Season >= 0 {
		return false
	}
	if s.Episode >= 0 {
		return false
	}

	if s.Resolution >= 0 {
		return false
	}

	if s.Verified {
		return false
	}
	if s.Uncensored {
		return false
	}
	if s.Bluray {
		return false
	}
	if s.NZB {
		return false
	}

	return true
}

//func (s *ReleaseSearch) Query() *elastic.BoolQuery {
//	query := elastic.NewBoolQuery()
//
//	//logrus.Debugf("search: %#v", s)
//
//	if s.Name != "" {
//		if s.Exact {
//			query = query.Must(elastic.NewTermQuery("name", "\""+s.Name+"\""))
//		} else {
//			query = query.Must(elastic.NewTermQuery("name", s.Name))
//		}
//	}
//
//	if s.Source != "" {
//		query = query.Must(elastic.NewTermQuery("source", s.Source))
//	}
//	if s.Type != "" {
//		query = query.Must(elastic.NewTermQuery("type", s.Type))
//	}
//
//	if s.Author != "" {
//		query = query.Must(elastic.NewTermQuery("author", "\""+s.Author+"\""))
//	}
//	if s.Group != "" {
//		query = query.Must(elastic.NewTermQuery("group", "\""+s.Group+"\""))
//	}
//
//	if s.Season >= 0 {
//		query = query.Must(elastic.NewTermQuery("season", s.Season))
//	}
//	if s.Episode >= 0 {
//		query = query.Must(elastic.NewTermQuery("episode", s.Episode))
//	}
//
//	if s.Resolution >= 0 {
//		query = query.Must(elastic.NewTermQuery("resolution", s.Resolution))
//	}
//
//	if s.Verified {
//		query = query.Must(elastic.NewTermQuery("verified", s.Verified))
//	}
//	if s.Uncensored {
//		query = query.Must(elastic.NewTermQuery("uncensored", s.Uncensored))
//	}
//	if s.Bluray {
//		query = query.Must(elastic.NewTermQuery("bluray", s.Bluray))
//	}
//
//	return query
//}
