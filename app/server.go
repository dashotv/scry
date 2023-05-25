package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"

	"github.com/dashotv/scry/nzbgeek"
	"github.com/dashotv/scry/search"
)

type Server struct {
	Config  *Config
	Log     *logrus.Entry
	Router  *gin.Engine
	Client  *search.Client
	Nzbgeek *nzbgeek.Client
}

func New() (*Server, error) {
	s := &Server{
		Config: App().Config,
		Log:    App().Log,
		Router: App().Router,
	}

	return s, nil
}

func (s *Server) Start() error {
	s.Log.Info("starting scry...")

	if s.Config.Cron {
		c := cron.New(cron.WithSeconds())

		//if _, err := c.AddFunc("* * * * * *", s.function); err != nil {
		//	return errors.Wrap(err, "adding cron function")
		//}

		go func() {
			s.Log.Info("starting cron...")
			c.Start()
		}()
	}

	s.Routes()

	//s.Jobs configuration

	s.Log.Info("starting web...")
	if err := s.Router.Run(fmt.Sprintf(":%d", s.Config.Port)); err != nil {
		return errors.Wrap(err, "starting router")
	}

	return nil
}
