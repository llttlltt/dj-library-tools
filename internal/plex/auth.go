package plex

import (
	"fmt"
	"net/http"
)

type PinResponse struct {
	ID        int    `json:"id"`
	Code      string `json:"code"`
	AuthToken string `json:"authToken"`
}

func (c *Client) RequestPin() (*PinResponse, error) {
	url := fmt.Sprintf("%s/pins?strong=true", PlexTvBaseURL)
	req, err := c.newRequest(http.MethodPost, url)
	if err != nil {
		return nil, err
	}

	var pin PinResponse
	if err := c.do(req, &pin); err != nil {
		return nil, err
	}

	return &pin, nil
}

func (c *Client) GetResources() ([]Resource, error) {
	url := fmt.Sprintf("%s/resources", PlexTvBaseURL)
	req, err := c.newRequest(http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	var resources []Resource
	if err := c.do(req, &resources); err != nil {
		return nil, err
	}

	return resources, nil
}

func (c *Client) GetPlaylists(baseURL string) ([]Playlist, error) {
	url := fmt.Sprintf("%s/playlists", baseURL)
	req, err := c.newRequest(http.MethodGet, url)
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

func (c *Client) GetPlaylistTracks(baseURL, playlistKey string) ([]Track, error) {
	url := fmt.Sprintf("%s%s", baseURL, playlistKey)
	req, err := c.newRequest(http.MethodGet, url)
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

