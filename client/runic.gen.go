// Code generated by github.com/dashotv/golem. DO NOT EDIT.
package client

import (
	"context"
	"fmt"

	"github.com/dashotv/fae"
	"github.com/dashotv/scry/search"
)

type RunicService struct {
	client *Client
}

// NewRunic makes a new client for accessing Runic services.
func NewRunicService(client *Client) *RunicService {
	return &RunicService{
		client: client,
	}
}

type RunicIndexRequest struct {
	Start      int    `json:"start"`
	Limit      int    `json:"limit"`
	Type       string `json:"type"`
	Text       string `json:"text"`
	Year       int    `json:"year"`
	Season     int    `json:"season"`
	Episode    int    `json:"episode"`
	Group      string `json:"group"`
	Website    string `json:"website"`
	Resolution int    `json:"resolution"`
	Source     string `json:"source"`
	Uncensored bool   `json:"uncensored"`
	Bluray     bool   `json:"bluray"`
	Verified   bool   `json:"verified"`
	Exact      bool   `json:"exact"`
}

type RunicIndexResponse struct {
	*Response
	Result *search.RunicSearchResponse `json:"result"`
	Total  int64                       `json:"total"`
}

func (s *RunicService) Index(ctx context.Context, req *RunicIndexRequest) (*RunicIndexResponse, error) {
	result := &RunicIndexResponse{Response: &Response{}}
	resp, err := s.client.Resty.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(result).
		SetQueryParam("start", fmt.Sprintf("%v", req.Start)).
		SetQueryParam("limit", fmt.Sprintf("%v", req.Limit)).
		SetQueryParam("type", fmt.Sprintf("%v", req.Type)).
		SetQueryParam("text", fmt.Sprintf("%v", req.Text)).
		SetQueryParam("year", fmt.Sprintf("%v", req.Year)).
		SetQueryParam("season", fmt.Sprintf("%v", req.Season)).
		SetQueryParam("episode", fmt.Sprintf("%v", req.Episode)).
		SetQueryParam("group", fmt.Sprintf("%v", req.Group)).
		SetQueryParam("website", fmt.Sprintf("%v", req.Website)).
		SetQueryParam("resolution", fmt.Sprintf("%v", req.Resolution)).
		SetQueryParam("source", fmt.Sprintf("%v", req.Source)).
		SetQueryParam("uncensored", fmt.Sprintf("%v", req.Uncensored)).
		SetQueryParam("bluray", fmt.Sprintf("%v", req.Bluray)).
		SetQueryParam("verified", fmt.Sprintf("%v", req.Verified)).
		SetQueryParam("exact", fmt.Sprintf("%v", req.Exact)).
		Get("/runic/")
	if err != nil {
		return nil, fae.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return nil, fae.Errorf("%d: %v", resp.StatusCode(), resp.String())
	}
	if result.Error {
		return nil, fae.New(result.Message)
	}

	return result, nil
}