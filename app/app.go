package app

import (
	"fmt"

	"github.com/dashotv/tmdb"
	"github.com/dashotv/tvdb"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/dashotv/scry/nzbgeek"
	"github.com/dashotv/scry/search"
)

type setupFunc func(app *Application) error

type Application struct {
	Config  *Config
	ES      *elasticsearch.Client
	Log     *zap.SugaredLogger
	Engine  *gin.Engine
	Router  *gin.RouterGroup
	Default *gin.RouterGroup
	Client  *search.Client
	Nzbgeek *nzbgeek.Client
	Tvdb    *tvdb.Client
	Tmdb    *tmdb.Client
}

func New() (*Application, error) {
	app := &Application{}

	list := []setupFunc{
		setupConfig,
		setupLogger,
		setupRoutes,
		setupTmdb,
		setupTvdb,
		setupElasticsearch,
		setupClient,
		setupNzbgeek,
	}
	for _, f := range list {
		if err := f(app); err != nil {
			return nil, err
		}
	}

	return app, nil
}

func (a *Application) Start() error {
	a.Routes()

	a.Log.Info("starting scry...")
	if err := a.Engine.Run(fmt.Sprintf(":%d", a.Config.Port)); err != nil {
		return errors.Wrap(err, "starting router")
	}

	return nil
}
