package app

func onMovies(app *Application, event *EventMovies) error {
	msg := event.Movie
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
