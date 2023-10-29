package app

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/dashotv/tmdb"
	"github.com/dashotv/tvdb"
	"github.com/gin-gonic/gin"
)

func SearchIndex(c *gin.Context) {
	responses := searchAll(c)

	App().Log.Infof("searchIndex: %#v", responses)

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
		responses.Media = &Response{Results: r, Error: err}
	}()

	go func() {
		defer wg.Done()
		r, err := searchTmdb(name)
		responses.Tmdb = &Response{Results: r, Error: err}
	}()

	go func() {
		defer wg.Done()
		r, err := searchTvdb(name)
		responses.Tvdb = &Response{Results: r, Error: err}
	}()

	wg.Wait()

	if responses.Media.Error != nil {
		App().Log.Errorf("searchAll media error: %s", responses.Media.Error)
	}
	if responses.Tmdb.Error != nil {
		App().Log.Errorf("searchAll tmdb error: %s", responses.Tmdb.Error)
	}
	if responses.Tvdb.Error != nil {
		App().Log.Errorf("searchAll tvdb error: %s", responses.Tvdb.Error)
	}
	return responses
}

type Response struct {
	Results []*Result
	Error   error
}

type Result struct {
	ID     string
	Title  string
	Type   string
	Date   string
	Source string
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
			ID:     v.ID,
			Title:  v.Title,
			Type:   v.Type,
			Date:   v.ReleaseDate,
			Source: "media",
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
		Query: &q,
		Type:  tvdb.String("series"),
		Limit: tvdb.Int64(10),
	}
	r, err := App().Tvdb.GetSearchResults(req)
	if err != nil {
		return nil, err
	}

	for _, v := range r.Data {
		s := strings.Split(*v.ID, "-")
		a := &Result{
			ID:     s[1],
			Title:  *v.Name,
			Type:   "series",
			Source: "tvdb",
		}
		if v.FirstAirTime != nil {
			a.Date = *v.FirstAirTime
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
		Query: q,
	}
	r, err := App().Tmdb.SearchMovie(p)
	if err != nil {
		return nil, err
	}

	for _, v := range r.Results {
		out = append(out, &Result{
			ID:     fmt.Sprintf("%d", tmdb.Int64Value(v.ID)),
			Title:  tmdb.StringValue(v.Title),
			Type:   "movie",
			Date:   tmdb.StringValue(v.ReleaseDate),
			Source: "tmdb",
		})
	}

	if len(out) > 10 {
		return out[:10], nil
	}
	return out, nil
}
