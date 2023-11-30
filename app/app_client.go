package app

import "github.com/dashotv/scry/search"

func setupClient(app *Application) error {
	client, err := search.New(app.Config.ElasticsearchURL)
	if err != nil {
		return err
	}

	app.Client = client
	return nil
}
