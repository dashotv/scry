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

	App().Log.Infof("searchAll wait %s", name)
	wg.Wait()

	App().Log.Infof("searchAll return %s", name)
	return responses
}

type Response struct {
	Results []*Result
	Error   error
}

type Result struct {
	ID    string
	Title string
	Type  string
	Date  string
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
			ID:    v.ID,
			Title: v.Title,
			Type:  v.Type,
			Date:  v.ReleaseDate,
		})
	}

	return out, nil
}

func searchTvdb(q string) ([]*Result, error) {
	out := []*Result{}

	t := "series"
	var l float32 = 10
	r, err := App().Tvdb.GetSearchResults(&tvdb.GetSearchResultsParams{Query: &q, Type: &t, Limit: &l})
	if err != nil {
		return nil, err
	}

	if r.JSON200 == nil || r.JSON200.Data == nil || len(*r.JSON200.Data) == 0 {
		return out, nil
	}

	for _, v := range *r.JSON200.Data {
		a := &Result{
			ID:    fmt.Sprintf("%d", v.Id),
			Title: *v.Name,
			Type:  "series",
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

	r, err := App().Tmdb.SearchMovie(&tmdb.SearchMovieParams{Query: q})
	if err != nil {
		return nil, err
	}

	if r.JSON200 == nil || r.JSON200.Results == nil || len(*r.JSON200.Results) == 0 {
		return out, nil
	}

	for _, v := range *r.JSON200.Results {
		out = append(out, &Result{
			ID:    fmt.Sprintf("%d", v.Id),
			Title: *v.Title,
			Type:  "movie",
			Date:  *v.ReleaseDate,
		})
	}

	return out[:10], nil
}
