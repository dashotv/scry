package app

func onEpisodes(app *Application, event *EventEpisodes) error {
	m := event.Episode
	if event.Event == "deleted" {
		return app.Client.DeleteMedia(m.ID)
	}

	m.Name = m.Title
	if m.Display != "" {
		m.Name = m.Display
	}

	_, err := app.Client.IndexMedia(m)
	if err != nil {
		return err
	}
	return nil
}
