package server

import (
	"fmt"
	"net/http"

	"github.com/dashotv/scry/server/config"
	"github.com/dashotv/scry/server/releases"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Start(url string, port int, mode string) error {
	cfg := config.New(url, port, mode)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{})

	if mode == "release" {
		gin.SetMode(mode)
	}

	router := gin.Default()
	router.GET("/", homeIndex)

	err := releases.Routes(cfg, router)
	if err != nil {
		logrus.Fatalf("error: %s", err)
	}

	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		return err
	}

	return nil
}

func homeIndex(c *gin.Context) {
	c.String(http.StatusOK, "home")
}
