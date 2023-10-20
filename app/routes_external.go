package app

import (
	"net/http"

	"github.com/dashotv/tmdb"
	"github.com/dashotv/tvdb"
	"github.com/gin-gonic/gin"
)

func TmdbIndex(c *gin.Context) {
	r, err := App().Tmdb.SearchMulti(&tmdb.SearchMultiParams{Query: c.Query("q")})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, v := range *r.JSON200.Results {
		if *v.MediaType != "tv" && *v.MediaType != "movie" {
			continue
		}

	}

	c.JSON(http.StatusOK, gin.H{"tmdb": r.JSON200.Results})
}

func TvdbIndex(c *gin.Context) {
	q := c.Query("q")
	r, err := App().Tvdb.GetSearchResults(&tvdb.GetSearchResultsParams{Query: &q})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tvdb": r.JSON200.Data,
	})
}
