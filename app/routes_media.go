package app

import (
	"net/http"

	"github.com/dashotv/golem/web"
	"github.com/gin-gonic/gin"

	"github.com/dashotv/scry/search"
)

func (a *Application) MediaIndex(c *gin.Context) {
	a.Log.Debugf("params: %#v", c.Params)
	s, err := a.CreateSearch(c)
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

func (a *Application) CreateSearch(c *gin.Context) (*search.MediaSearch, error) {
	s := a.Client.Media.NewSearch()

	s.Start, _ = web.QueryDefaultInteger(c, "start", 0)
	s.Limit, _ = web.QueryDefaultInteger(c, "limit", 25)

	s.Type = c.Query("type")
	s.Name = c.Query("name")
	s.Display = c.Query("display")
	s.Title = c.Query("title")
	s.Source = c.Query("source")
	s.SourceID = c.Query("source_id")

	a.Log.Debugf("create: %#v", s)
	return s, nil
}
