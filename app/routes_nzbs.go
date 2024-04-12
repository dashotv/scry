package app

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/dashotv/scry/nzbgeek"
)

func (a *Application) NzbsTv(c echo.Context) error {
	options := &nzbgeek.TvSearchOptions{}
	options.RageID = c.QueryParam("rageid")
	options.Episode = c.QueryParam("episode")
	options.Season = c.QueryParam("season")
	options.TvdbID = c.QueryParam("tvdbid")

	a.Log.Debugf("options: %#v", options)

	response, err := a.Nzbgeek.TvSearch(options)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &Response{Error: true, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: response.Channel.Item})
}

func (a *Application) NzbsMovie(c echo.Context) error {
	options := &nzbgeek.MovieSearchOptions{}
	options.ImdbID = c.QueryParam("imdbid")

	a.Log.Debugf("options: %#v", options)

	response, err := a.Nzbgeek.MovieSearch(options)
	if err != nil {
		a.Log.Errorf("nzbgeek movie search: %s", err)
		return c.JSON(http.StatusBadRequest, &Response{Error: true, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: response.Channel.Item})
}
