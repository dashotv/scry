package app

import runic "github.com/dashotv/runic/app"

func onRunic(app *Application, msg *runic.Release) error {
	app.Log.Named("runic").Debugf("index: %+v", msg)
	_, err := app.Client.IndexRunic(msg)
	if err != nil {
		return err
	}

	return nil
}
