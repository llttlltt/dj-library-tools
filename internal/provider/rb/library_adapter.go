package rb

import (
	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/rekordbox"
)

// RekordboxLibrary is an adapter that makes RekordboxLibraryXML satisfy the Library interface.
type RekordboxLibrary struct {
	XML *rekordbox.RekordboxLibraryXML
}

func (r *RekordboxLibrary) GetTracks() []models.Track {
	rbTracks := r.XML.Collection.TRACK
	tracks := make([]models.Track, len(rbTracks))
	for i, rt := range rbTracks {
		tracks[i] = ToNeutralTrack(rt)
	}
	return tracks
}

func (r *RekordboxLibrary) GetPlaylists() []models.ResourceGroup {
	var results []models.ResourceGroup
	r.collectAllNodes(r.XML.Playlists.Node.Node, "", &results)
	return results
}

func (r *RekordboxLibrary) collectAllNodes(nodes []rekordbox.Node, parent string, out *[]models.ResourceGroup) {
	for _, n := range nodes {
		*out = append(*out, ToNeutralGroup(n, parent))
		if len(n.Node) > 0 {
			r.collectAllNodes(n.Node, n.Name, out)
		}
	}
}

func (r *RekordboxLibrary) GetMembershipMap() map[string][]string {
	m := make(map[string][]string)
	r.walkRekordboxPlaylists(r.XML.Playlists.Node.Node, m)
	return m
}

func (r *RekordboxLibrary) walkRekordboxPlaylists(nodes []rekordbox.Node, m map[string][]string) {
	for _, node := range nodes {
		if node.Type == 1 {
			for _, t := range node.TRACK {
				m[t.Key] = append(m[t.Key], node.Name)
			}
		}
		if len(node.Node) > 0 {
			r.walkRekordboxPlaylists(node.Node, m)
		}
	}
}

func (r *RekordboxLibrary) AddPlaylist(folder, name string, trackIDs []string, position int) {
	r.XML.PlaylistsChanged = true
	var container *[]rekordbox.Node
	var folderNode *rekordbox.Node
	if folder == "" {
		container = &r.XML.Playlists.Node.Node
	} else {
		folderNode = r.findOrCreateFolder(folder)
		container = &folderNode.Node
	}

	node := rekordbox.Node{
		Name:    name,
		Type:    1,
		KeyType: rekordbox.PtrInt32(0),
		Entries: rekordbox.PtrInt32(int32(len(trackIDs))),
	}
	for _, id := range trackIDs {
		node.TRACK = append(node.TRACK, struct {
			Key string `xml:"Key,attr"`
		}{Key: id})
	}

	if position < 0 || position >= len(*container) {
		*container = append(*container, node)
	} else {
		*container = append((*container)[:position+1], (*container)[position:]...)
		(*container)[position] = node
	}

	if folderNode != nil {
		if folderNode.Count == nil {
			folderNode.Count = rekordbox.PtrInt32(1)
		} else {
			*folderNode.Count++
		}
	}
}

func (r *RekordboxLibrary) UpdatePlaylist(name string, trackIDs []string) bool {
	r.XML.PlaylistsChanged = true
	node, _, _, _ := r.findNodeInTree(&r.XML.Playlists.Node.Node, nil, name, 1)
	if node == nil {
		return false
	}
	node.TRACK = nil
	for _, id := range trackIDs {
		node.TRACK = append(node.TRACK, struct {
			Key string `xml:"Key,attr"`
		}{Key: id})
	}
	node.Entries = rekordbox.PtrInt32(int32(len(trackIDs)))
	return true
}

func (r *RekordboxLibrary) AddTracksToPlaylist(name string, trackIDs []string) (bool, int) {
	r.XML.PlaylistsChanged = true
	node, _, _, _ := r.findNodeInTree(&r.XML.Playlists.Node.Node, nil, name, 1)
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
	node.Entries = rekordbox.PtrInt32(int32(len(node.TRACK)))
	return true, added
}

func (r *RekordboxLibrary) RemoveTracksFromPlaylist(name string, trackIDs []string) (bool, int) {
	r.XML.PlaylistsChanged = true
	node, _, _, _ := r.findNodeInTree(&r.XML.Playlists.Node.Node, nil, name, 1)
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
	node.Entries = rekordbox.PtrInt32(int32(len(node.TRACK)))
	return true, before - len(node.TRACK)
}

func (r *RekordboxLibrary) CreateFolder(folder, name string, position int) bool {
	r.XML.PlaylistsChanged = true
	var container *[]rekordbox.Node
	var folderNode *rekordbox.Node
	if folder == "" {
		container = &r.XML.Playlists.Node.Node
	} else {
		folderNode = r.findOrCreateFolder(folder)
		container = &folderNode.Node
	}

	node := rekordbox.Node{
		Name:  name,
		Type:  0,
		Count: rekordbox.PtrInt32(0),
	}

	if position < 0 || position >= len(*container) {
		*container = append(*container, node)
	} else {
		*container = append((*container)[:position+1], (*container)[position:]...)
		(*container)[position] = node
	}

	if folderNode != nil {
		if folderNode.Count == nil {
			folderNode.Count = rekordbox.PtrInt32(1)
		} else {
			*folderNode.Count++
		}
	}
	return true
}

func (r *RekordboxLibrary) RenameGroup(name, newName string, nodeType int32) bool {
	r.XML.PlaylistsChanged = true
	node, _, _, _ := r.findNodeInTree(&r.XML.Playlists.Node.Node, nil, name, nodeType)
	if node == nil {
		return false
	}
	node.Name = newName
	return true
}

func (r *RekordboxLibrary) MoveGroup(name string, nodeType int32, targetFolder string) bool {
	r.XML.PlaylistsChanged = true
	node, parentNode, parentSlice, idx := r.findNodeInTree(&r.XML.Playlists.Node.Node, nil, name, nodeType)
	if node == nil {
		return false
	}

	moved := *node
	*parentSlice = append((*parentSlice)[:idx], (*parentSlice)[idx+1:]...)
	if parentNode != nil && parentNode.Count != nil && *parentNode.Count > 0 {
		*parentNode.Count--
	}

	target := r.findOrCreateFolder(targetFolder)
	target.Node = append(target.Node, moved)
	if target.Count == nil {
		target.Count = rekordbox.PtrInt32(1)
	} else {
		*target.Count++
	}
	return true
}

func (r *RekordboxLibrary) RemoveGroup(name string, nodeType int32) bool {
	r.XML.PlaylistsChanged = true
	_, parentNode, parentSlice, idx := r.findNodeInTree(&r.XML.Playlists.Node.Node, nil, name, nodeType)
	if idx == -1 {
		return false
	}
	*parentSlice = append((*parentSlice)[:idx], (*parentSlice)[idx+1:]...)
	if parentNode != nil && parentNode.Count != nil && *parentNode.Count > 0 {
		*parentNode.Count--
	}
	return true
}

func (r *RekordboxLibrary) Save(path string) error {
	return rekordbox.WriteRekordboxLibrary(path, r.XML)
}

func (r *RekordboxLibrary) findOrCreateFolder(name string) *rekordbox.Node {
	nodes := &r.XML.Playlists.Node.Node
	for i := range *nodes {
		if (*nodes)[i].Name == name && (*nodes)[i].Type == 0 {
			return &(*nodes)[i]
		}
	}
	*nodes = append(*nodes, rekordbox.Node{
		Name:  name,
		Type:  0,
		Count: rekordbox.PtrInt32(0),
	})
	return &(*nodes)[len(*nodes)-1]
}

func (r *RekordboxLibrary) findNodeInTree(nodes *[]rekordbox.Node, parent *rekordbox.Node, name string, nodeType int32) (*rekordbox.Node, *rekordbox.Node, *[]rekordbox.Node, int) {
	for i := range *nodes {
		n := &(*nodes)[i]
		if n.Name == name && n.Type == nodeType {
			return n, parent, nodes, i
		}
		if len(n.Node) > 0 {
			if found, foundParent, foundSlice, idx := r.findNodeInTree(&n.Node, n, name, nodeType); found != nil {
				return found, foundParent, foundSlice, idx
			}
		}
	}
	return nil, nil, nil, -1
}

// NewRekordboxLibrary creates a new Library wrapper for a Rekordbox XML.
func NewRekordboxLibrary(xml *rekordbox.RekordboxLibraryXML) *RekordboxLibrary {
	return &RekordboxLibrary{XML: xml}
}
