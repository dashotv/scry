package releases

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dashotv/scry/search"
	"github.com/dashotv/scry/server/config"
	"github.com/dashotv/scry/server/util"
)

var client *search.Client

func Routes(cfg *config.Config, e *gin.Engine) error {
	var err error

	fmt.Printf("connecting to elasticsearch: %s\n", cfg.URL)
	client, err = search.New(cfg.URL)
	if err != nil {
		return err
	}

	r := e.Group("/releases")
	r.GET("/", Search)

	return nil
}

func Search(c *gin.Context) {
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

	s.Source = c.Query("source")
	s.Type = c.Query("type")
	s.Name = c.Query("text")
	s.Author = c.Query("author")
	s.Group = c.Query("group")

	s.Season, _ = util.QueryInteger(c, "season")
	s.Episode, _ = util.QueryInteger(c, "episode")
	s.Resolution, _ = util.QueryInteger(c, "resolution")

	s.Verified = c.DefaultQuery("verified", "false") == "true"
	s.Uncensored = c.Query("uncensored") == "true"
	s.Bluray = c.Query("bluray") == "true"
	s.Exact = c.Query("exact") == "true"

	fmt.Printf("    create: %#v\n", s)
	return s, nil
}
