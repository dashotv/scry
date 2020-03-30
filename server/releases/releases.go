package releases

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/dashotv/scry/search"
	"github.com/dashotv/scry/server/util"
)

var client *search.Client

func Routes(c *search.Client, e *gin.Engine) {
	client = c

	r := e.Group("/releases")
	r.GET("/", Search)
}

func Search(c *gin.Context) {
	logrus.Debugf("    params: %#v", c.Params)
	s, err := CreateSearch(c)
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

func CreateSearch(c *gin.Context) (*search.ReleaseSearch, error) {
	s := client.Release.NewSearch()

	s.Start, _ = util.QueryDefaultInteger(c, "start", 0)
	s.Limit, _ = util.QueryDefaultInteger(c, "limit", 25)

	s.Source = c.Query("source")
	s.Type = c.Query("type")
	s.Name = c.Query("text")
	s.Year, _ = util.QueryInteger(c, "year")
	s.Author = c.Query("author")
	s.Group = c.Query("group")

	s.Season, _ = util.QueryInteger(c, "season")
	s.Episode, _ = util.QueryInteger(c, "episode")
	s.Resolution, _ = util.QueryInteger(c, "resolution")

	s.Verified = c.DefaultQuery("verified", "false") == "true"
	s.Uncensored = c.Query("uncensored") == "true"
	s.Bluray = c.Query("bluray") == "true"
	s.Exact = c.Query("exact") == "true"

	logrus.Debugf("    create: %#v", s)
	return s, nil
}
