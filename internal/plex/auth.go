package plex

import (
	"context"
	"fmt"
	"net/http"
)

type PinResponse struct {
	ID        int    `json:"id"`
	Code      string `json:"code"`
	AuthToken string `json:"authToken"`
}

func (c *Client) RequestPin(ctx context.Context) (*PinResponse, error) {
	url := fmt.Sprintf("%s/pins?strong=true", PlexTvBaseURL)
	req, err := c.newRequest(ctx, http.MethodPost, url)
	if err != nil {
		return nil, err
	}

	var pin PinResponse
	if err := c.do(req, &pin); err != nil {
		return nil, err
	}

	return &pin, nil
}

func (c *Client) CheckPin(ctx context.Context, id int) (*PinResponse, error) {
	url := fmt.Sprintf("%s/pins/%d", PlexTvBaseURL, id)
	req, err := c.newRequest(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	var pin PinResponse
	if err := c.do(req, &pin); err != nil {
		return nil, err
	}

	return &pin, nil
}

func (c *Client) GetResources(ctx context.Context) ([]Resource, error) {
	url := fmt.Sprintf("%s/resources", PlexTvBaseURL)
	req, err := c.newRequest(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	var resources []Resource
	if err := c.do(req, &resources); err != nil {
		return nil, err
	}

	return resources, nil
}

func (c *Client) GetPlaylists(ctx context.Context, baseURL string) ([]Playlist, error) {
	url := fmt.Sprintf("%s/playlists", baseURL)
	req, err := c.newRequest(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	var container struct {
		MediaContainer MediaContainer `json:"MediaContainer"`
	}
	if err := c.do(req, &container); err != nil {
		return nil, err
	}

	return container.MediaContainer.Metadata, nil
}

func (c *Client) GetPlaylistTracks(ctx context.Context, baseURL, playlistKey string) ([]Track, error) {
	url := fmt.Sprintf("%s%s", baseURL, playlistKey)
	req, err := c.newRequest(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	var container struct {
		MediaContainer TrackContainer `json:"MediaContainer"`
	}
	if err := c.do(req, &container); err != nil {
		return nil, err
	}

	return container.MediaContainer.Metadata, nil
}

