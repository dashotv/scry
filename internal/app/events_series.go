package app

func onSeries(app *Application, event *EventSeries) error {
	msg := event.Series
	if event.Event == "deleted" {
		return app.Client.DeleteMedia(msg.ID)
	}

	msg.Name = msg.Title
	if msg.Display != "" {
		msg.Name = msg.Display
	}

	_, err := app.Client.IndexMedia(msg)
	if err != nil {
		return err
	}
	return nil
}
