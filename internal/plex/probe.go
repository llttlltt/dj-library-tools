package plex

import (
	"context"
	"sync"
	"time"
)

type ConnectionResult struct {
	BaseURL   string
	Playlists []Playlist
	Tracks    []Track
	Err       error
}

// ProbeBestConnection tries all connections in parallel and returns the first successful result.
func (c *Client) ProbeBestConnection(resource Resource) (*ConnectionResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	results := make(chan *ConnectionResult, len(resource.Connections))
	var wg sync.WaitGroup

	for _, conn := range resource.Connections {
		wg.Add(1)
		go func(uri string) {
			defer wg.Done()
			
			// Try to get playlists as a health check
			playlists, err := c.GetPlaylists(ctx, uri)
			if err == nil {
				select {
				case results <- &ConnectionResult{BaseURL: uri, Playlists: playlists}:
					cancel() // Cancel other probes once we have a winner
				case <-ctx.Done():
				}
			}
		}(conn.URI)
	}

	// Closer goroutine
	go func() {
		wg.Wait()
		close(results)
	}()

	select {
	case res := <-results:
		if res != nil {
			return res, nil
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return nil, nil
}
