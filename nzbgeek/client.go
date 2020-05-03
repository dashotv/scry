package nzbgeek

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

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
	Season  string
	Episode string
	RageID  string
	TvdbID  string
}

type SearchResponse struct {
	Offset int
	Total  int
	Result []SearchResult
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
	Title       string `json:"title"`
	Guid        string
	Link        string
	Comments    string
	Published   CustomTime `json:"pubDate"`
	Category    string
	Description string
	Enclosure   struct {
		Attributes struct {
			URL    string
			Length string
			Type   string
		} `json:"@attributes"`
	}
	Attributes []struct {
		Attribute struct {
			Name  string
			Value string
		} `json:"@attributes"`
	}
}

func (c *Client) TvSearch(options *TvSearchOptions) (*Response, error) {
	response := &Response{}
	params := url.Values{}
	params.Add("t", "tvsearch")
	if options.TvdbID != "" {
		params.Add("tvdbid", options.TvdbID)
	}
	if options.Season != "" {
		params.Add("season", options.Season)
	}
	if options.Episode != "" {
		params.Add("ep", options.Episode)
	}
	if options.RageID != "" {
		params.Add("rid", options.RageID)
	}
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
