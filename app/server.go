package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Server struct {
	App    *App
	Config *Config
	Log    *logrus.Entry
	Router *gin.Engine
}

func NewServer() (*Server, error) {
	app := Instance()
	s := &Server{
		App:    app,
		Config: ConfigInstance(),
		Log:    app.Log.WithField("prefix", "server"),
		Router: app.Router,
	}

	return s, nil
}

func (s *Server) Start() error {
	s.Log.Info("starting scry...")

	s.Routes()

	//s.Jobs configuration

	s.Log.Info("starting web...")
	if err := s.Router.Run(fmt.Sprintf(":%d", s.Config.Port)); err != nil {
		return errors.Wrap(err, "starting router")
	}

	return nil
}
