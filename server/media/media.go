package media

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

	r := e.Group("/media")
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

func CreateSearch(c *gin.Context) (*search.MediaSearch, error) {
	s := client.Media.NewSearch()

	s.Start, _ = util.QueryDefaultInteger(c, "start", 0)
	s.Limit, _ = util.QueryDefaultInteger(c, "limit", 25)

	s.Type = c.Query("type")
	s.Name = c.Query("name")
	s.Display = c.Query("display")
	s.Title = c.Query("title")

	logrus.Debugf("    create: %#v", s)
	return s, nil
}
