package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func homeIndex(c *gin.Context) {
	c.String(http.StatusOK, "home")
}
