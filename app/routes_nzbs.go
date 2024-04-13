package app

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/dashotv/scry/nzbgeek"
)

func (a *Application) NzbsTv(c echo.Context, tvdbid string, season, episode int) error {
	options := &nzbgeek.TvSearchOptions{}
	options.TvdbID = tvdbid
	if season >= 0 {
		options.Season = season
	}
	if episode >= 0 {
		options.Episode = episode
	}
	// options.RageID = c.QueryParam("rageid") // TODO: support rageid?

	a.Log.Debugf("options: %#v", options)

	response, err := a.Nzbgeek.TvSearch(options)
	if err != nil {
		a.Log.Errorf("nzbgeek tv search: %s", err)
		return c.JSON(http.StatusBadRequest, &Response{Error: true, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: response.Channel.Item})
}

func (a *Application) NzbsMovie(c echo.Context, imdbid, tmdbid string) error {
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
