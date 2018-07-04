package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/dashotv/scry/server/releases"
	"github.com/dashotv/scry/server/config"
)

func Start(url string, port int, mode string) error {
	cfg := config.New(url, port, mode)

	if mode == "release" {
		gin.SetMode(mode)
	}

	router := gin.Default()
	router.GET("/", homeIndex)

	releases.Routes(cfg, router)

	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		return err
	}

	return nil
}

func homeIndex(c *gin.Context) {
	c.String(http.StatusOK, "home")
}
