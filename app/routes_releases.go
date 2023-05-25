package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dashotv/golem/web"
	"github.com/dashotv/scry/search"
)

func ReleasesIndex(c *gin.Context) {
	App().Log.Debugf("    params: %#v", c.Params)
	s, err := CreateReleasesSearch(c)
	if err != nil {
		c.Error(err)
		return
	}

	res, err := s.Find()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func CreateReleasesSearch(c *gin.Context) (*search.ReleaseSearch, error) {
	s := App().Client.Release.NewSearch()

	s.Start, _ = web.QueryDefaultInteger(c, "start", 0)
	s.Limit, _ = web.QueryDefaultInteger(c, "limit", 25)

	s.Source = c.Query("source")
	s.Type = c.Query("type")
	s.Name = c.Query("text")
	s.Year, _ = web.QueryDefaultInteger(c, "year", -1)
	s.Author = c.Query("author")
	s.Group = c.Query("group")

	s.Season, _ = web.QueryDefaultInteger(c, "season", -1)
	s.Episode, _ = web.QueryDefaultInteger(c, "episode", -1)
	s.Resolution, _ = web.QueryDefaultInteger(c, "resolution", -1)

	s.Verified = c.DefaultQuery("verified", "false") == "true"
	s.Uncensored = c.Query("uncensored") == "true"
	s.Bluray = c.Query("bluray") == "true"
	s.Exact = c.Query("exact") == "true"

	App().Log.Debugf("    create: %#v", s)
	return s, nil
}
