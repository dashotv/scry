package app

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/dashotv/scry/search"
)

func (a *Application) ReleasesIndex(c echo.Context) error {
	s, err := a.CreateReleasesSearch(c)
	if err != nil {
		return err
	}

	res, err := s.Find()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (a *Application) CreateReleasesSearch(c echo.Context) (*search.ReleaseSearch, error) {
	s := a.Client.Release.NewSearch()

	s.Start, _ = QueryDefaultInteger(c, "start", 0)
	s.Limit, _ = QueryDefaultInteger(c, "limit", 25)

	s.Source = c.QueryParam("source")
	s.Type = c.QueryParam("type")
	s.Name = c.QueryParam("text")
	s.Year, _ = QueryDefaultInteger(c, "year", -1)
	s.Author = c.QueryParam("author")
	s.Group = c.QueryParam("group")

	s.Season, _ = QueryDefaultInteger(c, "season", -1)
	s.Episode, _ = QueryDefaultInteger(c, "episode", -1)
	s.Resolution, _ = QueryDefaultInteger(c, "resolution", -1)

	s.Verified = QueryBool(c, "verified")
	s.Uncensored = c.QueryParam("uncensored") == "true"
	s.Bluray = c.QueryParam("bluray") == "true"
	s.Exact = c.QueryParam("exact") == "true"

	a.Log.Debugf("create: %#v", s)
	return s, nil
}
