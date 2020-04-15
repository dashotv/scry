package nzbgeek

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

var geek = os.Getenv("NZBGEEK_URL")
var key = os.Getenv("NZBGEEK_KEY")

func TestClient_TvSearch(t *testing.T) {
	client := NewClient(geek, key)
	options := &TvSearchOptions{
		Season:  "1",
		Episode: "4",
		TvdbID:  "321231",
	}
	response, err := client.TvSearch(options)
	if err != nil {
		t.Error(err)
	}
	spew.Dump(response)
}
