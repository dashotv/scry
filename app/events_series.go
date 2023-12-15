package app

func onSeries(app *Application, event *EventSeries) error {
	msg := event.Series
	msg.Name = msg.Title
	if msg.Display != "" {
		msg.Name = msg.Display
	}
	app.Log.Warnf("Indexing series %#v", msg)
	_, err := app.Client.IndexMedia(msg)
	if err != nil {
		return err
	}
	return nil
}
