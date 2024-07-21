package app

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/dashotv/scry/search"
)

func (a *Application) MediaIndex(c echo.Context, start, limit int, types, name, display, title, source, source_id string, season, episode, absolute int, downloaded, completed, skipped string) error {
	s := a.Client.Media.NewSearch()
	s.Start = start
	s.Limit = limit
	s.Type = types
	s.Name = name
	s.Display = display
	s.Title = title
	s.Source = source
	s.SourceID = source_id
	if season > 0 {
		s.Season = season
	}
	if episode > 0 {
		s.Episode = episode
	}
	if absolute > 0 {
		s.Absolute = absolute
	}
	switch downloaded {
	case "true":
		s.Downloaded = true
	case "false":
		s.Downloaded = false
	}
	switch completed {
	case "true":
		s.Completed = true
	case "false":
		s.Completed = false
	}
	switch skipped {
	case "true":
		s.Skipped = true
	case "false":
		s.Skipped = false
	}

	res, err := s.Find()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (a *Application) CreateSearch(c echo.Context) (*search.MediaSearch, error) {
	s := a.Client.Media.NewSearch()

	s.Start, _ = QueryDefaultInteger(c, "start", 0)
	s.Limit, _ = QueryDefaultInteger(c, "limit", 25)

	s.Type = c.QueryParam("type")
	s.Name = c.QueryParam("name")
	s.Display = c.QueryParam("display")
	s.Title = c.QueryParam("title")
	s.Source = c.QueryParam("source")
	s.SourceID = c.QueryParam("source_id")

	a.Log.Debugf("create: %#v", s)
	return s, nil
}
