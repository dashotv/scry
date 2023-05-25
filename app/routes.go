// This file is autogenerated by Golem
// Do NOT make modifications, they will be lost
package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Routes() {
	s.Router.GET("/", homeHandler)

	media := s.Router.Group("/media")
	media.GET("/", mediaIndexHandler)

}

func homeHandler(c *gin.Context) {
	Index(c)
}

func Index(c *gin.Context) {
	c.String(http.StatusOK, "home")
}

// /media
func mediaIndexHandler(c *gin.Context) {

	MediaIndex(c)
}
