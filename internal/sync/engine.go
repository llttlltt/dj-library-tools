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
// plexClient may be nil if the caller only needs playlist write operations or SaveXML.
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

// UpsertPlaylist creates or replaces a named playlist inside folder.
// When folder is empty the playlist is placed at the root level.
// position is the 0-indexed position in the folder. -1 appends to the end.
// trackIDs must be Rekordbox TrackID strings (KeyType=0).
func (e *Engine) UpsertPlaylist(folder, name string, trackIDs []string, position int) *SyncResult {
	e.RBXML.PlaylistsChanged = true
	var container *[]rekordbox.Node
	var folderNode *rekordbox.Node
	if folder == "" {
		container = &e.RBXML.Playlists.Node.Node
	} else {
		folderNode = e.findOrCreateFolder(folder)
		container = &folderNode.Node
	}

	node := rekordbox.Node{
		BaseNode: rekordbox.BaseNode{Type: 1, Name: name},
		KeyType:  0,
		Entries:  int32(len(trackIDs)),
	}
	for _, id := range trackIDs {
		node.TRACK = append(node.TRACK, struct {
			Key string `xml:"Key,attr"`
		}{Key: id})
	}

	// Check if updating
	for i := range *container {
		if (*container)[i].Name == name && (*container)[i].Type == 1 {
			(*container)[i] = node
			return &SyncResult{PlaylistName: name, TracksInjected: len(trackIDs), Updated: true}
		}
	}

	// New playlist - handle position
	if position < 0 || position >= len(*container) {
		*container = append(*container, node)
	} else {
		// Insert at position
		*container = append((*container)[:position+1], (*container)[position:]...)
		(*container)[position] = node
	}

	if folderNode != nil {
		folderNode.Count++
	}
	return &SyncResult{PlaylistName: name, TracksInjected: len(trackIDs), Updated: false}
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
	e.RBXML.PlaylistsChanged = true
	node, _, _, _ := findNodeInTree(&e.RBXML.Playlists.Node.Node, nil, name, 1)
	if node == nil {
		return false, 0
	}

	existing := make(map[string]struct{}, len(node.TRACK))
	for _, t := range node.TRACK {
		existing[t.Key] = struct{}{}
	}

	added := 0
	for _, id := range trackIDs {
		if _, dup := existing[id]; !dup {
			node.TRACK = append(node.TRACK, struct {
				Key string `xml:"Key,attr"`
			}{Key: id})
			existing[id] = struct{}{}
			added++
		}
	}
	node.Entries = int32(len(node.TRACK))
	return true, added
}

// RemoveTracksFromPlaylist removes all trackIDs present in the given slice from a named playlist.
// Returns (true, removedCount) if the playlist was found, (false, 0) otherwise.
func (e *Engine) RemoveTracksFromPlaylist(name string, trackIDs []string) (bool, int) {
	e.RBXML.PlaylistsChanged = true
	node, _, _, _ := findNodeInTree(&e.RBXML.Playlists.Node.Node, nil, name, 1)
	if node == nil {
		return false, 0
	}

	toRemove := make(map[string]struct{}, len(trackIDs))
	for _, id := range trackIDs {
		toRemove[id] = struct{}{}
	}

	before := len(node.TRACK)
	kept := node.TRACK[:0]
	for _, t := range node.TRACK {
		if _, remove := toRemove[t.Key]; !remove {
			kept = append(kept, t)
		}
	}
	node.TRACK = kept
	node.Entries = int32(len(node.TRACK))
	return true, before - len(node.TRACK)
}

// RenameNode renames the first node matching name and nodeType anywhere in the tree.
// nodeType: 0=folder, 1=playlist.
// Returns false if no matching node is found.
func (e *Engine) RenameNode(name, newName string, nodeType int32) bool {
	e.RBXML.PlaylistsChanged = true
	node, _, _, _ := findNodeInTree(&e.RBXML.Playlists.Node.Node, nil, name, nodeType)
	if node == nil {
		return false
	}
	node.Name = newName
	return true
}

// MoveNode detaches the first node matching name and nodeType from its current location
// and re-attaches it inside targetFolder (creating the folder if it does not exist).
// Returns false if the node is not found.
func (e *Engine) MoveNode(name string, nodeType int32, targetFolder string) bool {
	e.RBXML.PlaylistsChanged = true
	node, parentNode, parentSlice, idx := findNodeInTree(&e.RBXML.Playlists.Node.Node, nil, name, nodeType)
	if node == nil {
		return false
	}

	moved := *node
	*parentSlice = append((*parentSlice)[:idx], (*parentSlice)[idx+1:]...)
	if parentNode != nil && parentNode.Count > 0 {
		parentNode.Count--
	}

	target := e.findOrCreateFolder(targetFolder)
	target.Node = append(target.Node, moved)
	target.Count++
	return true
}

// RemoveNode removes the first node matching name and nodeType from anywhere in the tree.
// Returns false if no matching node is found.
func (e *Engine) RemoveNode(name string, nodeType int32) bool {
	e.RBXML.PlaylistsChanged = true
	_, parentNode, parentSlice, idx := findNodeInTree(&e.RBXML.Playlists.Node.Node, nil, name, nodeType)
	if idx == -1 {
		return false
	}
	*parentSlice = append((*parentSlice)[:idx], (*parentSlice)[idx+1:]...)
	if parentNode != nil && parentNode.Count > 0 {
		parentNode.Count--
	}
	return true
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

// findNodeInTree recursively searches for the first node matching name and nodeType.
// parent is the node whose Node slice is being searched (nil at root level).
// Returns (found node, parent node, parent slice, index) or (nil, nil, nil, -1) when not found.
func findNodeInTree(nodes *[]rekordbox.Node, parent *rekordbox.Node, name string, nodeType int32) (*rekordbox.Node, *rekordbox.Node, *[]rekordbox.Node, int) {
	for i := range *nodes {
		n := &(*nodes)[i]
		if n.Name == name && n.Type == nodeType {
			return n, parent, nodes, i
		}
		if len(n.Node) > 0 {
			if found, foundParent, foundSlice, idx := findNodeInTree(&n.Node, n, name, nodeType); found != nil {
				return found, foundParent, foundSlice, idx
			}
		}
	}
	return nil, nil, nil, -1
}
