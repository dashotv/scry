package app

import "github.com/dashotv/scry/nzbgeek"

func init() {
	initializers = append(initializers, setupNzbgeek)
}

func setupNzbgeek(app *Application) error {
	app.Nzbgeek = nzbgeek.NewClient(app.Config.NzbgeekURL, app.Config.NzbgeekKey)
	return nil
}
