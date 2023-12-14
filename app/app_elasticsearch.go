package app

import (
	"bytes"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
)

func init() {
	initializers = append(initializers, setupElasticsearch)
}

func setupElasticsearch(app *Application) error {
	c := elasticsearch.Config{
		Addresses: []string{
			app.Config.ElasticsearchURL,
		},
	}

	var err error
	app.ES, err = elasticsearch.NewClient(c)
	if err != nil {
		return err
	}

	return nil
}

type ElasticSearchIndexPayload struct {
	Index string
	Model interface{}
}

func (a *Application) ElasticSearchIndex(payload any) error {
	p, ok := payload.(*ElasticSearchIndexPayload)
	if !ok {
		return errors.New("invalid payload")
	}

	data, err := json.Marshal(p.Model)
	if err != nil {
		return errors.Wrap(err, "marshaling model")
	}

	resp, err := a.ES.Index(p.Index, bytes.NewReader(data))
	if err != nil {
		return errors.Wrap(err, "indexing model")
	}

	a.Log.Infof("elasticsearch index response: %+v", resp)
	return nil
}
