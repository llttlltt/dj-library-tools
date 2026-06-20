package plex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	PlexTvBaseURL = "https://plex.tv/api/v2"
	ClientName    = "djlt"
	AppID         = "dj-library-tools"
)

type Client struct {
	HTTPClient *http.Client
	Token      string
}

func NewClient(token string) *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		Token: token,
	}
}

func (c *Client) newRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Plex-Client-Identifier", AppID)
	req.Header.Set("X-Plex-Product", ClientName)
	req.Header.Set("X-Plex-Device", "Terminal")
	req.Header.Set("Accept", "application/json")

	if c.Token != "" {
		req.Header.Set("X-Plex-Token", c.Token)
	}

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("plex api error: %s", resp.Status)
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}

	return nil
}
