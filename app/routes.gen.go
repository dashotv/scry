// Code generated by github.com/dashotv/golem. DO NOT EDIT.
package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dashotv/fae"
	"github.com/dashotv/golem/plugins/router"
	"github.com/labstack/echo/v4"
)

func init() {
	initializers = append(initializers, setupRoutes)
	healthchecks["routes"] = checkRoutes
	starters = append(starters, startRoutes)
}

func checkRoutes(app *Application) error {
	// TODO: check routes
	return nil
}

func startRoutes(ctx context.Context, app *Application) error {
	go func() {
		app.Routes()
		app.Log.Info("starting routes...")
		if err := app.Engine.Start(fmt.Sprintf(":%d", app.Config.Port)); err != nil {
			app.Log.Errorf("routes: %s", err)
		}
	}()
	return nil
}

func setupRoutes(app *Application) error {
	logger := app.Log.Named("routes").Desugar()
	e, err := router.New(logger)
	if err != nil {
		return fae.Wrap(err, "router plugin")
	}
	app.Engine = e
	// unauthenticated routes
	app.Default = app.Engine.Group("")
	// authenticated routes (if enabled, otherwise same as default)
	app.Router = app.Engine.Group("")

	// TODO: fix auth
	if app.Config.Auth {
		clerkSecret := app.Config.ClerkSecretKey
		if clerkSecret == "" {
			app.Log.Fatal("CLERK_SECRET_KEY is not set")
		}
		clerkToken := app.Config.ClerkToken
		if clerkToken == "" {
			app.Log.Fatal("CLERK_TOKEN is not set")
		}

		app.Router.Use(router.ClerkAuth(clerkSecret, clerkToken))
	}

	return nil
}

type Setting struct {
	Name  string `json:"name"`
	Value bool   `json:"value"`
}

type SettingsBatch struct {
	IDs   []string `json:"ids"`
	Name  string   `json:"name"`
	Value bool     `json:"value"`
}

type Response struct {
	Error   bool        `json:"error"`
	Message string      `json:"message,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Total   int64       `json:"total,omitempty"`
}

func (a *Application) Routes() {
	a.Default.GET("/", a.indexHandler)
	a.Default.GET("/health", a.healthHandler)

	es := a.Router.Group("/es")
	es.GET("/", a.EsIndexHandler)
	es.GET("/media", a.EsMediaHandler)
	es.GET("/release", a.EsReleaseHandler)

	media := a.Router.Group("/media")
	media.GET("/", a.MediaIndexHandler)

	nzbs := a.Router.Group("/nzbs")
	nzbs.GET("/movie", a.NzbsMovieHandler)
	nzbs.GET("/tv", a.NzbsTvHandler)

	releases := a.Router.Group("/releases")
	releases.GET("/", a.ReleasesIndexHandler)

	runic := a.Router.Group("/runic")
	runic.GET("/", a.RunicIndexHandler)

	search := a.Router.Group("/search")
	search.GET("/", a.SearchIndexHandler)

}

func (a *Application) indexHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, router.H{
		"name": "scry",
		"routes": router.H{
			"es":       "/es",
			"media":    "/media",
			"nzbs":     "/nzbs",
			"releases": "/releases",
			"runic":    "/runic",
			"search":   "/search",
		},
	})
}

func (a *Application) healthHandler(c echo.Context) error {
	health, err := a.Health()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, router.H{"name": "scry", "health": health})
}

// Es (/es)
func (a *Application) EsIndexHandler(c echo.Context) error {
	return a.EsIndex(c)
}
func (a *Application) EsMediaHandler(c echo.Context) error {
	return a.EsMedia(c)
}
func (a *Application) EsReleaseHandler(c echo.Context) error {
	return a.EsRelease(c)
}

// Media (/media)
func (a *Application) MediaIndexHandler(c echo.Context) error {
	start := router.QueryParamIntDefault(c, "start", "0")
	limit := router.QueryParamIntDefault(c, "limit", "25")
	types := router.QueryParamString(c, "types")
	name := router.QueryParamString(c, "name")
	display := router.QueryParamString(c, "display")
	title := router.QueryParamString(c, "title")
	source := router.QueryParamString(c, "source")
	source_id := router.QueryParamString(c, "source_id")
	return a.MediaIndex(c, start, limit, types, name, display, title, source, source_id)
}

// Nzbs (/nzbs)
func (a *Application) NzbsMovieHandler(c echo.Context) error {
	return a.NzbsMovie(c)
}
func (a *Application) NzbsTvHandler(c echo.Context) error {
	return a.NzbsTv(c)
}

// Releases (/releases)
func (a *Application) ReleasesIndexHandler(c echo.Context) error {
	types := router.QueryParamString(c, "types")
	text := router.QueryParamString(c, "text")
	year := router.QueryParamString(c, "year")
	season := router.QueryParamString(c, "season")
	episode := router.QueryParamString(c, "episode")
	group := router.QueryParamString(c, "group")
	author := router.QueryParamString(c, "author")
	resolution := router.QueryParamString(c, "resolution")
	source := router.QueryParamString(c, "source")
	uncensored := router.QueryParamBool(c, "uncensored")
	bluray := router.QueryParamBool(c, "bluray")
	verified := router.QueryParamBool(c, "verified")
	exact := router.QueryParamBool(c, "exact")
	return a.ReleasesIndex(c, types, text, year, season, episode, group, author, resolution, source, uncensored, bluray, verified, exact)
}

// Runic (/runic)
func (a *Application) RunicIndexHandler(c echo.Context) error {
	types := router.QueryParamString(c, "types")
	text := router.QueryParamString(c, "text")
	year := router.QueryParamString(c, "year")
	season := router.QueryParamString(c, "season")
	episode := router.QueryParamString(c, "episode")
	group := router.QueryParamString(c, "group")
	website := router.QueryParamString(c, "website")
	resolution := router.QueryParamString(c, "resolution")
	source := router.QueryParamString(c, "source")
	uncensored := router.QueryParamBool(c, "uncensored")
	bluray := router.QueryParamBool(c, "bluray")
	verified := router.QueryParamBool(c, "verified")
	exact := router.QueryParamBool(c, "exact")
	return a.RunicIndex(c, types, text, year, season, episode, group, website, resolution, source, uncensored, bluray, verified, exact)
}

// Search (/search)
func (a *Application) SearchIndexHandler(c echo.Context) error {
	start := router.QueryParamIntDefault(c, "start", "0")
	limit := router.QueryParamIntDefault(c, "limit", "25")
	types := router.QueryParamString(c, "types")
	q := router.QueryParamString(c, "q")
	name := router.QueryParamString(c, "name")
	return a.SearchIndex(c, start, limit, types, q, name)
}
