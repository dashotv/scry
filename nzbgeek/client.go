package nzbgeek

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Client struct {
	URL string
	Key string
}

func NewClient(URL, key string) *Client {
	return &Client{
		URL: URL,
		Key: key,
	}
}

type SearchOptions struct {
	T string
}

type TvSearchOptions struct {
	SearchOptions
	Season  string
	Episode string
	RageID  string
	TvdbID  string
}

func (o *TvSearchOptions) Params() url.Values {
	params := url.Values{}
	params.Add("t", "tvsearch")
	if o.TvdbID != "" {
		params.Add("tvdbid", o.TvdbID)
	}
	if o.Season != "" {
		params.Add("season", o.Season)
	}
	if o.Episode != "" {
		params.Add("ep", o.Episode)
	}
	if o.RageID != "" {
		params.Add("rid", o.RageID)
	}
	return params
}

type MovieSearchOptions struct {
	SearchOptions
	ImdbID string
}

func (o *MovieSearchOptions) Params() url.Values {
	params := url.Values{}
	params.Add("t", "movie")
	if o.ImdbID != "" {
		id := strings.Replace(o.ImdbID, "tt", "", 1)
		params.Add("imdbid", id)
	}
	return params
}

type Response struct {
	Attributes struct {
		Version string
	}
	Channel struct {
		Title string
		Item  []SearchResult
	}
}

type SearchResult struct {
	Title       string     `json:"title"`
	Guid        string     `json:"guid"`
	Link        string     `json:"link"`
	Comments    string     `json:"comments"`
	Published   CustomTime `json:"pubDate"`
	Category    string     `json:"category"`
	Description string     `json:"description"`
	Enclosure   struct {
		Attributes struct {
			URL    string `json:"url"`
			Length string `json:"length"`
			Type   string `json:"type"`
		} `json:"@attributes"`
	} `json:"enclosure"`
	Attributes []struct {
		Attribute struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"@attributes"`
	} `json:"attributes"`
}

func (c *Client) TvSearch(options *TvSearchOptions) (*Response, error) {
	response := &Response{}
	params := options.Params()
	err := c.request("", params, response)
	if err != nil {
		return nil, err
	}
	if len(response.Channel.Item) == 0 {
		response.Channel.Item = []SearchResult{}
	}
	return response, err
}

func (c *Client) MovieSearch(options *MovieSearchOptions) (*Response, error) {
	response := &Response{}
	params := options.Params()
	err := c.request("", params, response)
	if err != nil {
		return nil, err
	}
	if len(response.Channel.Item) == 0 {
		response.Channel.Item = []SearchResult{}
	}
	return response, err
}

func (c *Client) request(path string, params url.Values, target interface{}) error {
	var err error
	var request *http.Request
	var response *http.Response
	var body []byte

	params.Add("limit", "50")
	params.Add("apikey", c.Key)
	params.Add("o", "json")

	u := fmt.Sprintf("%s/%s", c.URL, path)

	if request, err = http.NewRequest("GET", u, nil); err != nil {
		return errors.Wrap(err, "creating "+u+" request failed")
	}
	request.URL.RawQuery = params.Encode()

	logrus.Debugf("request url: %s", request.URL)

	client := &http.Client{}
	if response, err = client.Do(request); err != nil {
		//log.Fatal(err)
		return errors.Wrap(err, "error making http request")
	}
	defer response.Body.Close()

	if body, err = ioutil.ReadAll(response.Body); err != nil {
		return errors.Wrap(err, "reading request body")
	}

	logrus.Debugf("body: %s", string(body))

	if target == nil {
		return nil
	}

	if err = json.Unmarshal(body, &target); err != nil {
		return errors.Wrap(err, "json unmarshal")
	}

	return nil
}
