package sync

import (
	"github.com/llttlltt/dj-library-tools/internal/engine"
	"github.com/llttlltt/dj-library-tools/internal/plex"
)

// PlexSyncFolder is the top-level folder name injected into Rekordbox by djlt.
const PlexSyncFolder = "Plex Sync"

// Engine manages sync operations against a music library.
type Engine struct {
	PlexClient *plex.Client
	Library    engine.WritableLibrary
	Matcher    *Matcher
}

// NewEngine creates a sync Engine backed by the given library.
// plexClient may be nil if the caller only needs playlist write operations or Save.
func NewEngine(plexClient *plex.Client, lib engine.WritableLibrary) *Engine {
	return &Engine{
		PlexClient: plexClient,
		Library:    lib,
		Matcher:    NewMatcher(lib.GetTracks()),
	}
}

// SyncResult holds the outcome of a single playlist injection.
type SyncResult struct {
	PlaylistName   string
	TracksInjected int
	// Updated is true when an existing playlist was replaced; false when newly created.
	Updated bool
}

// UpsertPlaylist creates or replaces a named playlist inside folder.
// When folder is empty the playlist is placed at the root level.
// position is the 0-indexed position in the folder. -1 appends to the end.
// trackIDs must be Rekordbox TrackID strings (KeyType=0).
func (e *Engine) UpsertPlaylist(folder, name string, trackIDs []string, position int) *SyncResult {
	updated := e.Library.UpdatePlaylist(name, trackIDs)
	if !updated {
		e.Library.AddPlaylist(folder, name, trackIDs, position)
	}

	return &SyncResult{
		PlaylistName:   name,
		TracksInjected: len(trackIDs),
		Updated:        updated,
	}
}

// InjectPlaylist upserts a named playlist under PlexSyncFolder.
// Preserved for backward compatibility; delegates to UpsertPlaylist.
func (e *Engine) InjectPlaylist(name string, trackIDs []string) *SyncResult {
	return e.UpsertPlaylist(PlexSyncFolder, name, trackIDs, -1)
}

// AddTracksToPlaylist appends trackIDs to a named playlist anywhere in the tree.
// Duplicate IDs are silently ignored.
// Returns (true, addedCount) if the playlist was found, (false, 0) otherwise.
func (e *Engine) AddTracksToPlaylist(name string, trackIDs []string) (bool, int) {
	return e.Library.AddTracksToPlaylist(name, trackIDs)
}

// RemoveTracksFromPlaylist removes all trackIDs present in the given slice from a named playlist.
// Returns (true, removedCount) if the playlist was found, (false, 0) otherwise.
func (e *Engine) RemoveTracksFromPlaylist(name string, trackIDs []string) (bool, int) {
	return e.Library.RemoveTracksFromPlaylist(name, trackIDs)
}

// CreateFolder creates a new folder node at the specified position.
func (e *Engine) CreateFolder(folder, name string, position int) bool {
	return e.Library.CreateFolder(folder, name, position)
}

// RenameNode renames the first node matching name and nodeType anywhere in the tree.
// nodeType: 0=folder, 1=playlist.
// Returns false if no matching node is found.
func (e *Engine) RenameNode(name, newName string, nodeType int32) bool {
	return e.Library.RenameNode(name, newName, nodeType)
}

// MoveNode detaches the first node matching name and nodeType from its current location
// and re-attaches it inside targetFolder (creating the folder if it does not exist).
// Returns false if the node is not found.
func (e *Engine) MoveNode(name string, nodeType int32, targetFolder string) bool {
	return e.Library.MoveNode(name, nodeType, targetFolder)
}

// RemoveNode removes the first node matching name and nodeType from anywhere in the tree.
// Returns false if no matching node is found.
func (e *Engine) RemoveNode(name string, nodeType int32) bool {
	return e.Library.RemoveNode(name, nodeType)
}

// RemovePlaylist removes a named playlist from anywhere in the tree.
// Preserved for backward compatibility; delegates to RemoveNode.
func (e *Engine) RemovePlaylist(name string) bool {
	return e.RemoveNode(name, 1)
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

// Save writes the modified library back to disk.
func (e *Engine) Save(path string) error {
	return e.Library.Save(path)
}

// SaveXML is an alias for Save for backward compatibility.
func (e *Engine) SaveXML(path string) error {
	return e.Save(path)
}
