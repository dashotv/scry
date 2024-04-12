package app

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"

	"github.com/dashotv/tmdb"
	"github.com/dashotv/tvdb"
)

func (a *Application) SearchIndex(c echo.Context) error {
	responses := a.searchAll(c)
	return c.JSON(http.StatusOK, responses)
}

func (a *Application) searchAll(c echo.Context) *SearchAllResponse {
	wg := sync.WaitGroup{}
	wg.Add(3)

	name := c.QueryParam("q")
	responses := &SearchAllResponse{}

	go func() {
		defer wg.Done()
		r, err := a.searchMedia(c)
		e := ""
		if err != nil {
			e = err.Error()
		}
		responses.Media = &SearchResponse{Results: r, Error: e}
	}()

	go func() {
		defer wg.Done()
		r, err := a.searchTmdb(name)
		e := ""
		if err != nil {
			e = err.Error()
		}
		responses.Tmdb = &SearchResponse{Results: r, Error: e}
	}()

	go func() {
		defer wg.Done()
		r, err := a.searchTvdb(name)
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

func (a *Application) searchMedia(c echo.Context) ([]*SearchResult, error) {
	out := []*SearchResult{}

	s, err := a.CreateSearch(c)
	if err != nil {
		return nil, err
	}

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
		})
	}

	return out, nil
}

func (a *Application) searchTvdb(q string) ([]*SearchResult, error) {
	out := []*SearchResult{}
	if q == "" {
		return out, nil
	}

	req := tvdb.GetSearchResultsRequest{
		Query:    &q,
		Type:     tvdb.String("series"),
		Limit:    tvdb.Int64(10),
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

func (a *Application) searchTmdb(q string) ([]*SearchResult, error) {
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

	if len(out) > 10 {
		return out[:10], nil
	}
	return out, nil
}
