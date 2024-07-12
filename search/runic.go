package search

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"

	runic "github.com/dashotv/runic/client"
)

const RUNIC_PAGE_SIZE = 25

type RunicService struct {
	Service
}

func (s *RunicService) Index(r *runic.Release) (*elastic.IndexResponse, error) {
	return s.client.Index().
		Index(fmt.Sprintf("%s_%s_%s", s.index, s.env, timeToDateBucket(r.PublishedAt))).
		Type("_doc").
		Id(r.ID.Hex()).
		BodyJson(r).
		Do(context.Background())
}
func (s *RunicService) Delete(id string) error {
	_, err := s.client.Delete().
		Index(s.index + "_" + s.env + "_*").
		Type("_doc").
		Id(id).
		Do(context.Background())
	return err
}

func (s *RunicService) NewSearch() *RunicSearch {
	return &RunicSearch{
		client:     s.client,
		Season:     -1,
		Episode:    -1,
		Resolution: -1,
		Search:     &Search{Start: 0, Limit: RUNIC_PAGE_SIZE, Index: s.index + "_" + s.env + "_*"},
	}
}

type RunicSearch struct {
	Type        string `bson:"type" json:"type"`
	Source      string `bson:"source" json:"source"`
	Title       string `bson:"title" json:"title"`
	Year        int    `bson:"year" json:"year"`
	Description string `bson:"description" json:"description"`
	Size        int64  `bson:"size" json:"size"`
	View        string `bson:"view" json:"view"`
	Download    string `bson:"download" json:"download"`
	Infohash    string `bson:"infohash" json:"infohash"`
	Season      int    `bson:"season" json:"season"`
	Episode     int    `bson:"episode" json:"episode"`
	Volume      int    `bson:"volume" json:"volume"`
	Group       string `bson:"group" json:"group"`
	Website     string `bson:"website" json:"website"`
	Verified    bool   `bson:"verified" json:"verified"`
	Widescreen  bool   `bson:"widescreen" json:"widescreen"`
	Unrated     bool   `bson:"unrated" json:"unrated"`
	Uncensored  bool   `bson:"uncensored" json:"uncensored"`
	Bluray      bool   `bson:"bluray" json:"bluray"`
	ThreeD      bool   `bson:"threeD" json:"threeD"`
	Resolution  int    `bson:"resolution" json:"resolution"`
	Encodings   string `bson:"encoding" json:"encoding"`
	Quality     string `bson:"quality" json:"quality"`
	Downloader  string `bson:"downloader" json:"downloader"`
	Checksum    string `bson:"checksum" json:"checksum"`
	Exact       bool

	client *elastic.Client
	*Search
}

type RunicSearchResponse struct {
	*SearchResponse
	Releases []*runic.Release
}

func (s *RunicSearch) Find() (*RunicSearchResponse, error) {
	var q elastic.Query

	r := &RunicSearchResponse{SearchResponse: &SearchResponse{}}
	ctx := context.Background()

	search := s.client.Search().Index(s.Index)
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

	r.Total = sr.Hits.TotalHits.Value
	r.Count = len(sr.Hits.Hits)

	rels, err := s.processResponse(sr)
	if err != nil {
		return r, err
	}
	r.Releases = rels

	return r, nil
}

func (s *RunicSearch) processResponse(res *elastic.SearchResult) ([]*runic.Release, error) {
	var rels []*runic.Release

	if res == nil || res.TotalHits() == 0 {
		return rels, nil
	}

	for _, hit := range res.Hits.Hits {
		rel := &runic.Release{}
		if err := json.Unmarshal(hit.Source, rel); err != nil {
			return nil, err
		}

		rels = append(rels, rel)
	}

	return rels, nil
}

func (s *RunicSearch) Query() (*elastic.QueryStringQuery, string) {
	list := []string{}

	if s.Title != "" {
		if s.Exact {
			list = append(list, fmt.Sprintf("%s:\"%s\"", "title", s.Title))
		} else {
			words := strings.Split(s.Title, " ")
			list = append(list, fmt.Sprintf("%s:(%s)", "title", strings.Join(words, " AND ")))
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

	if s.Website != "" {
		list = append(list, fmt.Sprintf("%s:\"%s\"", "website", s.Website))
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
	if s.Downloader != "" {
		list = append(list, fmt.Sprintf("%s:%s", "downloader", s.Downloader))
	}

	str := strings.Join(list, " AND ")
	logrus.Debugf("    search: %s", str)
	return elastic.NewQueryStringQuery(str), str
}

func (s *RunicSearch) IsZero() bool {
	if s.Title != "" {
		return false
	}

	if s.Source != "" {
		return false
	}
	if s.Type != "" {
		return false
	}

	if s.Website != "" {
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
	if s.Downloader != "" {
		return false
	}

	return true
}
