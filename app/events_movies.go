package app

func onMovies(app *Application, event *EventMovies) error {
	msg := event.Movie
	if event.Event == "deleted" {
		return app.Client.DeleteMedia(msg.ID)
	}

	msg.Name = msg.Title
	if msg.Display != "" {
		msg.Name = msg.Display
	}
	msg.Type = "movie"

	_, err := app.Client.IndexMedia(msg)
	if err != nil {
		return err
	}
	return nil
}
