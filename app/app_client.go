package app

import "github.com/dashotv/scry/search"

func init() {
	initializers = append(initializers, setupClient)
}

func setupClient(app *Application) error {
	client, err := search.New(app.Config.ElasticsearchURL)
	if err != nil {
		return err
	}

	app.Client = client
	return nil
}
