package app

import "github.com/dashotv/scry/search"

func onReleases(app *Application, msg *search.Release) error {
	_, err := app.Client.IndexRelease(msg)
	if err != nil {
		return err
	}
	return nil
}
