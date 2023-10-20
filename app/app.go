package app

import (
	"sync"

	"github.com/dashotv/tmdb"
	"github.com/dashotv/tvdb"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"github.com/dashotv/scry/nzbgeek"
	"github.com/dashotv/scry/search"
)

var once sync.Once
var instance *Application

func App() *Application {
	once.Do(func() {
		instance = initialize()
	})
	return instance
}

type Application struct {
	Config *Config
	Router *gin.Engine
	// Cache  *redis.Client
	Log *logrus.Entry
	// Add additional clients and connections
	Client  *search.Client
	Nzbgeek *nzbgeek.Client
	Tvdb    *tvdb.Tvdb
	Tmdb    *tmdb.Tmdb
}

func logger() *logrus.Entry {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&prefixed.TextFormatter{DisableTimestamp: false, FullTimestamp: true})
	return logrus.WithField("prefix", "app")
}

func initialize() *Application {
	cfg := ConfigInstance()
	log := logger()

	if cfg.Mode == "dev" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if cfg.Mode == "release" {
		gin.SetMode(cfg.Mode)
	}

	router := gin.New()
	router.Use(ginlogrus.Logger(log), gin.Recovery())

	log.Infof("connecting to elasticsearch: %s", cfg.Elasticsearch.URL)
	client, err := search.New(cfg.Elasticsearch.URL)
	if err != nil {
		log.Fatalf("failed to connect to Elasticsearch: %s", err)
	}

	log.Infof("setting up nzbgeek...")
	nzbg := nzbgeek.NewClient(cfg.Nzbgeek.URL, cfg.Nzbgeek.Key)

	// TODO: add this to config
	// cache := redis.NewClient(&redis.Options{
	//	Addr: "localhost:6379",
	//	DB:   15, // use default DB
	// })

	// Add additional clients and connections
	tvdbClient := tvdb.New(cfg.Tvdb.URL)
	_, err = tvdbClient.Login(cfg.Tvdb.Key)
	if err != nil {
		log.Fatalf("failed to connect to TVDB: %s", err)
	}

	tmdbClient, err := tmdb.New(cfg.Tmdb.URL, cfg.Tmdb.Token)
	if err != nil {
		log.Fatalf("failed to connect to TMDB: %s", err)
	}

	return &Application{
		Config: cfg,
		Router: router,
		// Cache:    cache,
		Log:     log,
		Client:  client,
		Nzbgeek: nzbg,
		Tvdb:    tvdbClient,
		Tmdb:    tmdbClient,
	}
}
