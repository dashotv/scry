package server

import (
	"github.com/gin-gonic/gin"
	"github.com/dashotv/scry/server/releases"
	"net/http"
	"github.com/dashotv/scry/server/config"
)

func Start(url string) error {
	cfg := config.New(url)

	router := gin.Default()
	router.GET("/", homeIndex)

	releases.Routes(cfg, router)

	if err := router.Run(":8080"); err != nil {
		return err
	}

	return nil
}

func homeIndex(c *gin.Context) {
	c.String(http.StatusOK, "home")
}
