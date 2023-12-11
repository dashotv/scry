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
type healthFunc func(app *Application) error

type Application struct {
	Client  *search.Client
	Config  *Config
	Default *gin.RouterGroup
	Engine  *gin.Engine
	ES      *elasticsearch.Client
	Events  *Events
	Log     *zap.SugaredLogger
	Nzbgeek *nzbgeek.Client
	Router  *gin.RouterGroup
	Tmdb    *tmdb.Client
	Tvdb    *tvdb.Client
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
		setupEvents,
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

	go a.Events.Start()

	a.Log.Info("starting scry...")
	if err := a.Engine.Run(fmt.Sprintf(":%d", a.Config.Port)); err != nil {
		return errors.Wrap(err, "starting router")
	}

	return nil
}

func (a *Application) Health() (map[string]bool, error) {
	list := map[string]healthFunc{
		// "elasticsearch": healthElasticsearch,
		// "nzbgeek":       healthNzbgeek,
		// "tvdb":          healthTvdb,
		// "tmdb":          healthTmdb,
		"events": healthEvents,
	}

	resp := make(map[string]bool)
	for n, f := range list {
		err := f(a)
		resp[n] = err == nil
	}

	return resp, nil
}
