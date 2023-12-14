package app

import "github.com/dashotv/tvdb"

func init() {
	initializers = append(initializers, setupTvdb)
}

func setupTvdb(app *Application) error {
	t, err := tvdb.Login(app.Config.TvdbKey)
	if err != nil {
		return err
	}

	app.Tvdb = t
	return nil
}
