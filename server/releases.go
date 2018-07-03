package server

import (
	"net/http"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/dashotv/scry/search"
)

func releasesSearch(c *gin.Context) {
	s, err := releaseSearcher(c)
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

func releaseSearcher(c *gin.Context) (*search.ReleaseSearch, error) {
	s := client.Release.NewSearch()

	s.Source = c.Query("source")
	s.Type = c.Query("type")
	s.Name = c.Query("text")
	s.Author = c.Query("author")
	s.Group = c.Query("group")

	s.Season, _ = queryInteger(c, "season")
	s.Episode, _ = queryInteger(c, "episode")
	s.Resolution, _ = queryInteger(c, "resolution")

	s.Verified = c.DefaultQuery("verified", "true") == "true"
	s.Uncensored = c.Query("uncensored") == "true"
	s.Bluray = c.Query("bluray") == "true"
	s.Exact = c.Query("exact") == "true"

	return s, nil
}

func queryInteger(c *gin.Context, name string) (int, error) {
	v := c.Query(name)
	if v == "" {
		return -1, fmt.Errorf("not set")
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return -1, err
	}

	return n, nil
}
