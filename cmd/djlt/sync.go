package main

import (
	"fmt"
	"os"

	"github.com/llttlltt/dj-library-tools/internal/plex"
	"github.com/llttlltt/dj-library-tools/internal/sync"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync playlists between providers",
}

var syncPlexCmd = &cobra.Command{
	Use:   "plex [playlist-id]",
	Short: "Sync a Plex playlist to Rekordbox XML",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		playlistID := args[0]
		token := os.Getenv("PLEX_TOKEN")
		if token == "" {
			return fmt.Errorf("PLEX_TOKEN environment variable not set")
		}

		if xmlPath == "" {
			return fmt.Errorf("rekordbox XML path not specified (use --xml or -x)")
		}

		rbXML, err := rekordbox.ReadRekordboxLibrary(xmlPath)
		if err != nil {
			return fmt.Errorf("failed to read rekordbox library: %w", err)
		}

		client := plex.NewClient(token)
		resources, err := client.GetResources()
		if err != nil {
			return fmt.Errorf("failed to get resources: %w", err)
		}

		var targetServer plex.Resource
		found := false
		for _, res := range resources {
			if res.Provides == "server" && len(res.Connections) > 0 {
				targetServer = res
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("no plex servers found")
		}

		baseURL := targetServer.Connections[0].URI
		serverClient := plex.NewClient(targetServer.AccessToken)

		engine := sync.NewEngine(serverClient, rbXML)

		// Fetch the playlist title first to name it in Rekordbox
		playlists, err := serverClient.GetPlaylists(baseURL)
		if err != nil {
			return fmt.Errorf("failed to get playlists: %w", err)
		}

		var targetPlaylist plex.Playlist
		for _, pl := range playlists {
			if pl.RatingKey == playlistID {
				targetPlaylist = pl
				break
			}
		}

		if targetPlaylist.RatingKey == "" {
			return fmt.Errorf("playlist with ID %s not found", playlistID)
		}

		fmt.Printf("Syncing playlist: %s (%d tracks)...\n", targetPlaylist.Title, targetPlaylist.LeafCount)

		tracks, err := serverClient.GetPlaylistTracks(baseURL, targetPlaylist.Key)
		if err != nil {
			return fmt.Errorf("failed to get tracks: %w", err)
		}

		var rbTracks []string
		for _, track := range tracks {
			match := engine.Matcher.Match(track)
			if match.RBTrack != nil && match.Confidence >= 0.8 {
				fmt.Printf("  Matched: %s - %s (Confidence: %.2f)\n", track.Artist, track.Title, match.Confidence)
				rbTracks = append(rbTracks, fmt.Sprintf("%d", match.RBTrack.TrackID))
			} else {
				fmt.Printf("  No match for: %s - %s\n", track.Artist, track.Title)
			}
		}

		// Inject into XML
		err = injectPlaylist(rbXML, targetPlaylist.Title, rbTracks)
		if err != nil {
			return fmt.Errorf("failed to inject playlist: %w", err)
		}

		err = rekordbox.WriteRekordboxLibrary(xmlPath, rbXML)
		if err != nil {
			return fmt.Errorf("failed to save XML: %w", err)
		}

		fmt.Printf("Successfully synced %d/%d tracks to Rekordbox XML\n", len(rbTracks), len(tracks))

		return nil
	},
}

func injectPlaylist(lib *rekordbox.RekordboxLibraryXML, name string, trackIDs []string) error {
	// Find or create "Plex Sync" folder
	var syncFolder *rekordbox.Node
	for i := range lib.Playlists.Node.Node {
		if lib.Playlists.Node.Node[i].Name == "Plex Sync" && lib.Playlists.Node.Node[i].Type == 0 {
			syncFolder = &lib.Playlists.Node.Node[i]
			break
		}
	}

	if syncFolder == nil {
		lib.Playlists.Node.Node = append(lib.Playlists.Node.Node, rekordbox.Node{
			BaseNode: rekordbox.BaseNode{
				Type: 0,
				Name: "Plex Sync",
			},
		})
		syncFolder = &lib.Playlists.Node.Node[len(lib.Playlists.Node.Node)-1]
	}

	// Create or overwrite playlist inside sync folder
	newPlaylist := rekordbox.Node{
		BaseNode: rekordbox.BaseNode{
			Type: 1,
			Name: name,
		},
		KeyType: 0, // Track ID
		Entries: int32(len(trackIDs)),
	}

	for _, id := range trackIDs {
		newPlaylist.TRACK = append(newPlaylist.TRACK, struct {
			Key string `xml:"Key,attr"`
		}{Key: id})
	}

	// Remove existing if it exists
	found := false
	for i := range syncFolder.Node {
		if syncFolder.Node[i].Name == name && syncFolder.Node[i].Type == 1 {
			syncFolder.Node[i] = newPlaylist
			found = true
			break
		}
	}

	if !found {
		syncFolder.Node = append(syncFolder.Node, newPlaylist)
		syncFolder.Count++
	}

	return nil
}

func init() {
	syncCmd.AddCommand(syncPlexCmd)
	rootCmd.AddCommand(syncCmd)
}
