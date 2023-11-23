package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dashotv/scry/nzbgeek"
)

func NzbsTv(c *gin.Context) {
	options := &nzbgeek.TvSearchOptions{}
	options.RageID = c.Query("rageid")
	options.Episode = c.Query("episode")
	options.Season = c.Query("season")
	options.TvdbID = c.Query("tvdbid")

	App().Log.Debugf("options: %#v", options)

	response, err := App().Nzbgeek.TvSearch(options)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	App().Log.Debugf("response: %#v", response)

	c.JSON(http.StatusOK, response.Channel.Item)
}

func NzbsMovie(c *gin.Context) {
	options := &nzbgeek.MovieSearchOptions{}
	options.ImdbID = c.Query("imdbid")

	App().Log.Debugf("options: %#v", options)

	response, err := App().Nzbgeek.MovieSearch(options)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	App().Log.Debugf("response: %#v", response)

	c.JSON(http.StatusOK, response.Channel.Item)
}
