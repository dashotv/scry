package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/dashotv/scry/server/media"
	"github.com/dashotv/scry/server/releases"

	"github.com/dashotv/scry/search"
)

type Server struct {
	URL    string
	Port   int
	Mode   string
	Debug  bool
	Client *search.Client
}

func (s *Server) Start() error {
	var err error

	if s.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetFormatter(&logrus.TextFormatter{})

	if s.Mode == "release" {
		gin.SetMode(s.Mode)
	}

	logrus.Infof("connecting to elasticsearch: %s", s.URL)
	s.Client, err = search.New(s.URL)
	if err != nil {
		return err
	}

	router := gin.Default()
	router.GET("/", homeIndex)

	releases.Routes(s.Client, router)
	media.Routes(s.Client, router)

	if err := router.Run(fmt.Sprintf(":%d", s.Port)); err != nil {
		return err
	}

	return nil
}

//func Start(url string, port int, mode string) error {
//	cfg := config.New(url, port, mode)
//	logrus.SetLevel(logrus.InfoLevel)
//	logrus.SetFormatter(&logrus.TextFormatter{})
//
//	if mode == "release" {
//		gin.SetMode(mode)
//	}
//
//	router := gin.Default()
//	router.GET("/", homeIndex)
//
//	err := releases.Routes(cfg, router)
//	if err != nil {
//		logrus.Fatalf("error: %s", err)
//	}
//
//	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
//		return err
//	}
//
//	return nil
//}

func homeIndex(c *gin.Context) {
	c.String(http.StatusOK, "home")
}
