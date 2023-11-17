package app

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/dashotv/tmdb"
	"github.com/dashotv/tvdb"
	"github.com/gin-gonic/gin"
)

func SearchIndex(c *gin.Context) {
	responses := searchAll(c)

	c.JSON(http.StatusOK, responses)
}

type searchAllResponse struct {
	Media *Response
	Tmdb  *Response
	Tvdb  *Response
}

func searchAll(c *gin.Context) *searchAllResponse {
	wg := sync.WaitGroup{}
	wg.Add(3)

	name := c.Query("q")
	responses := &searchAllResponse{}

	go func() {
		defer wg.Done()
		r, err := searchMedia(c)
		e := ""
		if err != nil {
			e = err.Error()
		}
		responses.Media = &Response{Results: r, Error: e}
	}()

	go func() {
		defer wg.Done()
		r, err := searchTmdb(name)
		e := ""
		if err != nil {
			e = err.Error()
		}
		responses.Tmdb = &Response{Results: r, Error: e}
	}()

	go func() {
		defer wg.Done()
		r, err := searchTvdb(name)
		e := ""
		if err != nil {
			e = err.Error()
		}
		responses.Tvdb = &Response{Results: r, Error: e}
	}()

	wg.Wait()

	if responses.Media.Error != "" {
		App().Log.Errorf("searchAll media error: %s", responses.Media.Error)
	}
	if responses.Tmdb.Error != "" {
		App().Log.Errorf("searchAll tmdb error: %s", responses.Tmdb.Error)
	}
	if responses.Tvdb.Error != "" {
		App().Log.Errorf("searchAll tvdb error: %s", responses.Tvdb.Error)
	}
	return responses
}

type Response struct {
	Results []*Result
	Error   string
}

type Result struct {
	ID          string
	Title       string
	Description string
	Type        string
	Kind        string
	Date        string
	Source      string
	Image       string
}

func searchMedia(c *gin.Context) ([]*Result, error) {
	out := []*Result{}

	s, err := CreateSearch(c)
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
		out = append(out, &Result{
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

func searchTvdb(q string) ([]*Result, error) {
	out := []*Result{}
	if q == "" {
		return out, nil
	}

	req := tvdb.GetSearchResultsRequest{
		Query:    &q,
		Type:     tvdb.String("series"),
		Limit:    tvdb.Int64(10),
		Language: tvdb.String("eng"),
	}
	r, err := App().Tvdb.GetSearchResults(req)
	if err != nil {
		return nil, err
	}

	for _, v := range r.Data {
		a := &Result{
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

func searchTmdb(q string) ([]*Result, error) {
	out := []*Result{}
	if q == "" {
		return out, nil
	}

	p := tmdb.SearchMovieRequest{
		Query:    q,
		Language: tmdb.String("en-US"),
	}
	r, err := App().Tmdb.SearchMovie(p)
	if err != nil {
		return nil, err
	}

	for _, v := range r.Results {
		img := tmdb.StringValue(v.PosterPath)
		if img != "" {
			img = "https://image.tmdb.org/t/p/original" + img
		}
		out = append(out, &Result{
			ID:          fmt.Sprintf("%d", tmdb.Int64Value(v.ID)),
			Title:       tmdb.StringValue(v.Title),
			Description: tmdb.StringValue(v.Overview),
			Type:        "movie",
			Kind:        "movie",
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
