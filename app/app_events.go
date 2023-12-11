package app

import (
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/dashotv/mercury"

	"github.com/dashotv/scry/search"
)

type EventsChannel string
type EventsTopic string

type Events struct {
	Client   *search.Client
	Merc     *mercury.Mercury
	Log      *zap.SugaredLogger
	Series   chan *search.Media
	Movies   chan *search.Media
	Releases chan *search.Release
}

func healthEvents(app *Application) error {
	switch app.Events.Merc.Status() {
	case nats.CONNECTED:
		return nil
	default:
		return errors.Errorf("nats status: %s", app.Events.Merc.Status())
	}
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
		Series:   make(chan *search.Media),
		Movies:   make(chan *search.Media),
		Releases: make(chan *search.Release),
	}

	if err := e.Merc.Receiver("tower.index.series", e.Series); err != nil {
		return err
	}
	if err := e.Merc.Receiver("tower.index.movies", e.Movies); err != nil {
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
		case m := <-e.Series:
			e.Log.Debugf("indexing series: %#v", m)
			resp, err := e.Client.IndexMedia(m)
			if err != nil {
				e.Log.Errorf("index media failed: %s", err)
				e.Log.Debugf("response: %#v", resp)
			}
		case m := <-e.Movies:
			e.Log.Debugf("indexing movie: %#v", m)
			resp, err := e.Client.IndexMedia(m)
			if err != nil {
				e.Log.Errorf("index media failed: %s", err)
				e.Log.Debugf("response: %#v", resp)
			}
		case m := <-e.Releases:
			e.Log.Debugf("indexing release: %#v", m)
			resp, err := e.Client.IndexRelease(m)
			if err != nil {
				e.Log.Errorf("index release failed: %s", err)
				e.Log.Debugf("response: %#v", resp)
			}
		}
	}
}
