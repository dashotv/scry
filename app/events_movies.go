package app

import "github.com/dashotv/scry/search"

func onMovies(app *Application, msg *search.Media) error {
	_, err := app.Client.IndexMedia(msg)
	if err != nil {
		return err
	}
	return nil
}
