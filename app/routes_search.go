package app

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"

	"github.com/dashotv/scry/search"
	"github.com/dashotv/tmdb"
	"github.com/dashotv/tvdb"
)

func (a *Application) SearchIndex(c echo.Context, start, limit int, types, q, name string) error {
	responses := a.searchAll(limit, q, name, types)
	return c.JSON(http.StatusOK, &Response{Error: false, Result: responses})
}

func (a *Application) searchAll(limit int, q, name, types string) *SearchAllResponse {
	wg := sync.WaitGroup{}
	wg.Add(3)

	responses := &SearchAllResponse{}

	go func() {
		defer wg.Done()

		s := a.Client.Media.NewSearch()
		s.Start = 0
		s.Limit = limit
		s.Type = types
		s.Name = name

		r, err := a.searchMedia(s)
		e := ""
		if err != nil {
			e = err.Error()
		}
		responses.Media = &SearchResponse{Results: r, Error: e}
	}()

	go func() {
		defer wg.Done()
		r, err := a.searchTmdb(q, limit)
		e := ""
		if err != nil {
			e = err.Error()
		}
		responses.Tmdb = &SearchResponse{Results: r, Error: e}
	}()

	go func() {
		defer wg.Done()
		r, err := a.searchTvdb(q, limit)
		e := ""
		if err != nil {
			e = err.Error()
		}
		responses.Tvdb = &SearchResponse{Results: r, Error: e}
	}()

	wg.Wait()

	if responses.Media.Error != "" {
		a.Log.Errorf("searchAll media error: %s", responses.Media.Error)
	}
	if responses.Tmdb.Error != "" {
		a.Log.Errorf("searchAll tmdb error: %s", responses.Tmdb.Error)
	}
	if responses.Tvdb.Error != "" {
		a.Log.Errorf("searchAll tvdb error: %s", responses.Tvdb.Error)
	}
	return responses
}

func (a *Application) searchMedia(s *search.MediaSearch) ([]*SearchResult, error) {
	out := []*SearchResult{}

	r, err := s.Find()
	if err != nil {
		return nil, err
	}

	if len(r.Media) == 0 {
		return out, nil
	}

	for _, v := range r.Media {
		out = append(out, &SearchResult{
			ID:          v.ID,
			Title:       v.Title,
			Description: v.Description,
			Type:        v.Type,
			Kind:        v.Kind,
			Date:        v.ReleaseDate,
			Source:      "media",
			Image:       v.Cover,
			Completed:   v.Completed,
		})
	}

	return out, nil
}

func (a *Application) searchTvdb(q string, limit int) ([]*SearchResult, error) {
	out := []*SearchResult{}
	if q == "" {
		return out, nil
	}

	req := tvdb.GetSearchResultsRequest{
		Query:    &q,
		Type:     tvdb.String("series"),
		Limit:    tvdb.Float64(float64(limit)),
		Language: tvdb.String("eng"),
	}
	r, err := a.Tvdb.GetSearchResults(req)
	if err != nil {
		return nil, err
	}

	for _, v := range r.Data {
		a := &SearchResult{
			ID:          tvdb.StringValue(v.TvdbID),
			Title:       tvdb.StringValue(v.Name),
			Description: tvdb.StringValue(v.Overview),
			Type:        "series",
			Kind:        "tv",
			Source:      "tvdb",
			Image:       tvdb.StringValue(v.Thumbnail),
		}
		if v.FirstAirTime != nil {
			a.Date = *v.FirstAirTime
		}
		if v.PrimaryLanguage != nil && *v.PrimaryLanguage != "eng" {
			if v.Translations["eng"] != "" {
				a.Title = v.Translations["eng"]
			}
			if v.Overviews["eng"] != "" {
				a.Description = v.Overviews["eng"]
			}
		}
		out = append(out, a)
	}

	return out, nil
}

func (a *Application) searchTmdb(q string, limit int) ([]*SearchResult, error) {
	out := []*SearchResult{}
	if q == "" {
		return out, nil
	}

	p := tmdb.SearchMovieRequest{
		Query:    q,
		Language: tmdb.String("en-US"),
	}
	r, err := a.Tmdb.SearchMovie(p)
	if err != nil {
		return nil, err
	}

	for _, v := range r.Results {
		img := tmdb.StringValue(v.PosterPath)
		if img != "" {
			img = "https://image.tmdb.org/t/p/original" + img
		}
		out = append(out, &SearchResult{
			ID:          fmt.Sprintf("%d", tmdb.Int64Value(v.ID)),
			Title:       tmdb.StringValue(v.Title),
			Description: tmdb.StringValue(v.Overview),
			Type:        "movie",
			Kind:        "movies",
			Date:        tmdb.StringValue(v.ReleaseDate),
			Source:      "tmdb",
			Image:       img,
		})
	}

	if len(out) > limit {
		return out[:limit], nil
	}
	return out, nil
}
