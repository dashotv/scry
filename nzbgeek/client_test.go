package nzbgeek

import (
	"fmt"
	"os"
	"testing"
)

var geek = os.Getenv("NZBGEEK_URL")
var key = os.Getenv("NZBGEEK_KEY")

func printNzbs(response *Response) {
	for _, i := range response.Channel.Item {
		fmt.Printf(" - %s\n", i.Title)
	}
}

func TestClient_TvSearch(t *testing.T) {
	client := NewClient(geek, key)
	options := &TvSearchOptions{
		//Season:  "3",
		//Episode: "4",
		TvdbID: "376729", // Tower of God
	}
	response, err := client.TvSearch(options)
	if err != nil {
		t.Error(err)
	}
	printNzbs(response)
}

func TestClient_MovieSearch(t *testing.T) {
	client := NewClient(geek, key)
	ids := []string{
		"4263482",   // The Witch 2015
		"tt3794354", // Sonic the Hedgehog
	}

	for _, id := range ids {
		options := &MovieSearchOptions{ImdbID: id}
		response, err := client.MovieSearch(options)
		if err != nil {
			t.Error(err)
		}
		printNzbs(response)
	}
}
