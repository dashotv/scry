package nzbs

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/dashotv/scry/nzbgeek"
)

var client *nzbgeek.Client

func Routes(c *nzbgeek.Client, e *gin.Engine) {
	client = c

	r := e.Group("/nzbs")
	r.GET("/tv", TvSearch)
}

func TvSearch(c *gin.Context) {
	options := &nzbgeek.TvSearchOptions{}
	options.RageID = c.Query("rageid")
	options.Episode = c.Query("episode")
	options.Season = c.Query("season")
	options.TvdbID = c.Query("tvdbid")

	logrus.Debugf("options: %#v", options)

	response, err := client.TvSearch(options)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	logrus.Debugf("response: %#v", response)

	c.JSON(http.StatusOK, response.Channel.Item)
}
