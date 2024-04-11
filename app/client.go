package app

import (
	"github.com/pkg/errors"

	"github.com/dashotv/scry/search"
)

func init() {
	initializers = append(initializers, setupClient)
}

func setupClient(app *Application) error {
	app.Log.Infof("search: %t", app.Config.Production)
	client, err := search.New(app.Config.ElasticsearchURL, app.Config.Production)
	if err != nil {
		return errors.Wrap(err, "failed to create search client")
	}

	app.Client = client
	return nil
}
