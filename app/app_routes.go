// Code generated by github.com/dashotv/golem. DO NOT EDIT.
package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.infratographer.com/x/echox/echozap"
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
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(echozap.Middleware(logger))

	app.Engine = e
	// unauthenticated routes
	app.Default = app.Engine.Group("")
	// authenticated routes (if enabled, otherwise same as default)
	app.Router = app.Engine.Group("")

	// if app.Config.Auth {
	// 	clerkSecret := app.Config.ClerkSecretKey
	// 	if clerkSecret == "" {
	// 		app.Log.Fatal("CLERK_SECRET_KEY is not set")
	// 	}
	//
	// 	clerkClient, err := clerk.NewClient(clerkSecret)
	// 	if err != nil {
	// 		app.Log.Fatalf("clerk: %s", err)
	// 	}
	//
	// 	app.Router.Use(requireSession(clerkClient))
	// }

	return nil
}

// Enable Auth and uncomment to use Clerk to manage auth
// also add this import: "github.com/clerkinc/clerk-sdk-go/clerk"
//
// requireSession wraps the clerk.RequireSession middleware
// func requireSession(client clerk.Client) HandlerFunc {
// 	requireActiveSession := clerk.RequireSessionV2(client)
// 	return func(gctx *gin.Context) {
// 		var skip = true
// 		var handler http.HandlerFunc = func(http.ResponseWriter, *http.Request) {
// 			skip = false
// 		}
// 		requireActiveSession(handler).ServeHTTP(gctx.Writer, gctx.Request)
// 		switch {
// 		case skip:
// 			gctx.AbortWithStatusJSON(http.StatusBadRequest, H{"error": "session required"})
// 		default:
// 			gctx.Next()
// 		}
// 	}
// }

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
	return c.JSON(http.StatusOK, H{
		"name": "scry",
		"routes": H{
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
	return c.JSON(http.StatusOK, H{"name": "scry", "health": health})
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
	return a.MediaIndex(c)
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
	return a.ReleasesIndex(c)
}

// Runic (/runic)
func (a *Application) RunicIndexHandler(c echo.Context) error {
	return a.RunicIndex(c)
}

// Search (/search)
func (a *Application) SearchIndexHandler(c echo.Context) error {
	return a.SearchIndex(c)
}
