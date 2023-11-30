package app

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

func setupRoutes(app *Application) error {
	if app.Config.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	logger := app.Log.Named("routes")

	app.Engine = gin.New()
	app.Engine.Use(
		ginzap.Ginzap(logger.Desugar(), time.RFC3339, true),
		ginzap.RecoveryWithZap(logger.Desugar(), true),
	)
	app.Default = app.Engine.Group("/")
	app.Router = app.Engine.Group("/")

	return nil
}
