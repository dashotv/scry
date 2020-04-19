package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/dashotv/scry/nzbgeek"
	"github.com/dashotv/scry/search"
	"github.com/dashotv/scry/server/config"
	"github.com/dashotv/scry/server/media"
	"github.com/dashotv/scry/server/nzbs"
	"github.com/dashotv/scry/server/releases"
)

type Server struct {
	Config  *config.Config
	Client  *search.Client
	Nzbgeek *nzbgeek.Client
}

func (s *Server) Start() error {
	var err error

	if s.Config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetFormatter(&logrus.TextFormatter{})

	if s.Config.Mode == "release" {
		gin.SetMode(s.Config.Mode)
	}

	logrus.Infof("connecting to elasticsearch: %s", s.Config.Elasticsearch.URL)
	s.Client, err = search.New(s.Config.Elasticsearch.URL)
	if err != nil {
		return err
	}

	logrus.Infof("setting up nzbgeek...")
	s.Nzbgeek = nzbgeek.NewClient(s.Config.Nzbgeek.URL, s.Config.Nzbgeek.Key)

	router := gin.Default()
	router.GET("/", homeIndex)

	releases.Routes(s.Client, router)
	media.Routes(s.Client, router)
	nzbs.Routes(s.Nzbgeek, router)

	if err := router.Run(fmt.Sprintf(":%d", s.Config.Port)); err != nil {
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
