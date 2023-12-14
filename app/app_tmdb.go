package app

import (
	"github.com/dashotv/tmdb"
)

func init() {
	initializers = append(initializers, setupTmdb)
}

func setupTmdb(app *Application) error {
	app.Tmdb = tmdb.New(app.Config.TmdbToken)
	return nil
}
