package app

import (
	"github.com/pkg/errors"

	"github.com/dashotv/scry/search"
)

func init() {
	initializers = append(initializers, setupClient)
}

func setupClient(app *Application) error {
	client, err := search.New(app.Config.ElasticsearchURL)
	if err != nil {
		return errors.Wrap(err, "failed to create search client")
	}

	app.Client = client
	return nil
}
