package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dashotv/scry/search"
)

func (a *Application) EsIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (a *Application) EsMedia(c *gin.Context) {
	m := &search.Media{}
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := a.Client.IndexMedia(m)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": resp})
}

func (a *Application) EsRelease(c *gin.Context) {
	r := &search.Release{}
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := a.Client.IndexRelease(r)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": resp})
}
