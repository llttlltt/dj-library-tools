package sync

import (
	"github.com/llttlltt/dj-library-tools/internal/plex"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

// PlexSyncFolder is the top-level folder name injected into Rekordbox by djlt.
const PlexSyncFolder = "Plex Sync"

// Engine manages sync operations against a Rekordbox XML library.
// PlexClient may be nil when only XML injection is needed.
type Engine struct {
	PlexClient *plex.Client
	RBXML      *rekordbox.RekordboxLibraryXML
	Matcher    *Matcher
}

// NewEngine creates a sync Engine backed by the given Rekordbox XML.
// plexClient may be nil if the caller only needs InjectPlaylist / SaveXML.
func NewEngine(plexClient *plex.Client, rbXML *rekordbox.RekordboxLibraryXML) *Engine {
	return &Engine{
		PlexClient: plexClient,
		RBXML:      rbXML,
		Matcher:    NewMatcher(rbXML.Collection.TRACK),
	}
}

// SyncResult holds the outcome of a single playlist injection.
type SyncResult struct {
	PlaylistName   string
	TracksInjected int
	// Updated is true when an existing playlist was replaced; false when newly created.
	Updated bool
}

// InjectPlaylist upserts a named playlist under PlexSyncFolder in the XML.
// trackIDs must be Rekordbox TrackID strings (KeyType=0).
// The folder is created if it does not already exist.
// If a playlist with the same name already exists in the folder it is replaced,
// otherwise a new one is appended.
func (e *Engine) InjectPlaylist(name string, trackIDs []string) *SyncResult {
	folder := e.findOrCreateFolder(PlexSyncFolder)

	node := rekordbox.Node{
		BaseNode: rekordbox.BaseNode{
			Type: 1,
			Name: name,
		},
		KeyType: 0,
		Entries: int32(len(trackIDs)),
	}
	for _, id := range trackIDs {
		node.TRACK = append(node.TRACK, struct {
			Key string `xml:"Key,attr"`
		}{Key: id})
	}

	for i := range folder.Node {
		if folder.Node[i].Name == name && folder.Node[i].Type == 1 {
			folder.Node[i] = node
			return &SyncResult{PlaylistName: name, TracksInjected: len(trackIDs), Updated: true}
		}
	}

	folder.Node = append(folder.Node, node)
	folder.Count++
	return &SyncResult{PlaylistName: name, TracksInjected: len(trackIDs), Updated: false}
}

// RemovePlaylist removes a named playlist from PlexSyncFolder.
// Returns true if the playlist was found and removed.
func (e *Engine) RemovePlaylist(name string) bool {
	for i := range e.RBXML.Playlists.Node.Node {
		folder := &e.RBXML.Playlists.Node.Node[i]
		if folder.Name != PlexSyncFolder || folder.Type != 0 {
			continue
		}
		for j := range folder.Node {
			if folder.Node[j].Name == name && folder.Node[j].Type == 1 {
				folder.Node = append(folder.Node[:j], folder.Node[j+1:]...)
				if folder.Count > 0 {
					folder.Count--
				}
				return true
			}
		}
	}
	return false
}

// MatchTracks matches a slice of Plex tracks against the Rekordbox collection,
// returning only results at or above minConfidence.
func (e *Engine) MatchTracks(plexTracks []plex.Track, minConfidence float64) []MatchResult {
	out := make([]MatchResult, 0, len(plexTracks))
	for _, t := range plexTracks {
		m := e.Matcher.Match(t)
		if m.RBTrack != nil && m.Confidence >= minConfidence {
			out = append(out, m)
		}
	}
	return out
}

// SaveXML writes the modified library back to disk.
func (e *Engine) SaveXML(path string) error {
	return rekordbox.WriteRekordboxLibrary(path, e.RBXML)
}

// findOrCreateFolder returns a pointer to the named top-level folder node,
// creating it (Type=0) if it does not exist.
func (e *Engine) findOrCreateFolder(name string) *rekordbox.Node {
	nodes := &e.RBXML.Playlists.Node.Node
	for i := range *nodes {
		if (*nodes)[i].Name == name && (*nodes)[i].Type == 0 {
			return &(*nodes)[i]
		}
	}
	*nodes = append(*nodes, rekordbox.Node{
		BaseNode: rekordbox.BaseNode{Type: 0, Name: name},
	})
	return &(*nodes)[len(*nodes)-1]
}
