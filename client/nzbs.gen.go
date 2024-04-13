// Code generated by github.com/dashotv/golem. DO NOT EDIT.
package client

import (
	"context"
	"fmt"

	"github.com/dashotv/fae"
	"github.com/dashotv/scry/nzbgeek"
)

type NzbsService struct {
	client *Client
}

// NewNzbs makes a new client for accessing Nzbs services.
func NewNzbsService(client *Client) *NzbsService {
	return &NzbsService{
		client: client,
	}
}

type NzbsMovieRequest struct {
	Imdbid string `json:"imdbid"`
	Tmdbid string `json:"tmdbid"`
}

type NzbsMovieResponse struct {
	*Response
	Result []nzbgeek.SearchResult `json:"result"`
}

func (s *NzbsService) Movie(ctx context.Context, req *NzbsMovieRequest) (*NzbsMovieResponse, error) {
	result := &NzbsMovieResponse{Response: &Response{}}
	resp, err := s.client.Resty.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(result).
		SetQueryParam("imdbid", fmt.Sprintf("%v", req.Imdbid)).
		SetQueryParam("tmdbid", fmt.Sprintf("%v", req.Tmdbid)).
		Get("/nzbs/movie")
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

type NzbsTvRequest struct {
	Tvdbid  string `json:"tvdbid"`
	Season  int    `json:"season"`
	Episode int    `json:"episode"`
}

type NzbsTvResponse struct {
	*Response
	Result []nzbgeek.SearchResult `json:"result"`
}

func (s *NzbsService) Tv(ctx context.Context, req *NzbsTvRequest) (*NzbsTvResponse, error) {
	result := &NzbsTvResponse{Response: &Response{}}
	resp, err := s.client.Resty.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(result).
		SetQueryParam("tvdbid", fmt.Sprintf("%v", req.Tvdbid)).
		SetQueryParam("season", fmt.Sprintf("%v", req.Season)).
		SetQueryParam("episode", fmt.Sprintf("%v", req.Episode)).
		Get("/nzbs/tv")
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
