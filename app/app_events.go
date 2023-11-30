package app

import (
	"github.com/dashotv/mercury"
	"go.uber.org/zap"

	"github.com/dashotv/scry/search"
)

type EventsChannel string
type EventsTopic string

type Events struct {
	Client   *search.Client
	Merc     *mercury.Mercury
	Log      *zap.SugaredLogger
	Media    chan *search.Media
	Releases chan *search.Release
}

func setupEvents(app *Application) error {
	m, err := mercury.New("scry", app.Config.NatsURL)
	if err != nil {
		return err
	}

	e := &Events{
		Client:   app.Client,
		Merc:     m,
		Log:      app.Log.Named("events"),
		Media:    make(chan *search.Media),
		Releases: make(chan *search.Release),
	}

	if err := e.Merc.Receiver("tower.index.media", e.Media); err != nil {
		return err
	}
	if err := e.Merc.Receiver("tower.index.releases", e.Releases); err != nil {
		return err
	}

	app.Events = e
	return nil
}

func (e *Events) Start() error {
	e.Log.Infof("starting events...")

	for {
		select {
		case m := <-e.Media:
			resp, err := e.Client.IndexMedia(m)
			if err != nil {
				e.Log.Errorf("index media failed: %s", err)
				e.Log.Debugf("response: %#v", resp)
			}
		case m := <-e.Releases:
			resp, err := e.Client.IndexRelease(m)
			if err != nil {
				e.Log.Errorf("index release failed: %s", err)
				e.Log.Debugf("response: %#v", resp)
			}
		}
	}
}
