package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dashotv/scry/nzbgeek"
)

func (a *Application) NzbsTv(c *gin.Context) {
	options := &nzbgeek.TvSearchOptions{}
	options.RageID = c.Query("rageid")
	options.Episode = c.Query("episode")
	options.Season = c.Query("season")
	options.TvdbID = c.Query("tvdbid")

	a.Log.Debugf("options: %#v", options)

	response, err := a.Nzbgeek.TvSearch(options)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, response.Channel.Item)
}

func (a *Application) NzbsMovie(c *gin.Context) {
	options := &nzbgeek.MovieSearchOptions{}
	options.ImdbID = c.Query("imdbid")

	a.Log.Debugf("options: %#v", options)

	response, err := a.Nzbgeek.MovieSearch(options)
	if err != nil {
		a.Log.Errorf("nzbgeek movie search: %s", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, response.Channel.Item)
}
