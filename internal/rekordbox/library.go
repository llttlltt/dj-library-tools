package rekordbox

import (
	"fmt"
	"github.com/llttlltt/dj-library-tools/internal/models"
)

// Library is an adapter that makes RekordboxLibraryXML satisfy the Library interface.
type Library struct {
	XML *RekordboxLibraryXML
}

func (r *Library) GetResources(kind string) []models.Resource {
	var results []models.Resource
	switch kind {
	case "track":
		for _, rt := range r.XML.Collection.TRACK {
			results = append(results, ToNeutralTrack(rt))
		}
	case "group":
		groups := r.GetPlaylists()
		for _, g := range groups {
			results = append(results, g)
		}
	}
	return results
}

func (r *Library) GetTracks() []models.Track {
	rbTracks := r.XML.Collection.TRACK
	tracks := make([]models.Track, len(rbTracks))
	for i, rt := range rbTracks {
		tracks[i] = ToNeutralTrack(rt)
	}
	return tracks
}

func (r *Library) GetPlaylists() []models.ResourceGroup {
	var results []models.ResourceGroup
	r.collectAllGroups(r.XML.Playlists.Node.Node, "", &results)
	return results
}

func (r *Library) collectAllGroups(nodes []Node, parent string, out *[]models.ResourceGroup) {
	for _, n := range nodes {
		*out = append(*out, ToNeutralGroup(n, parent))
		if len(n.Node) > 0 {
			r.collectAllGroups(n.Node, n.Name, out)
		}
	}
}

func (r *Library) GetMembershipMap() map[string][]string {
	m := make(map[string][]string)
	r.walkRekordboxPlaylists(r.XML.Playlists.Node.Node, m)
	return m
}

func (r *Library) walkRekordboxPlaylists(nodes []Node, m map[string][]string) {
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

func (r *Library) CreateGroup(parentID, name string, groupKind models.GroupKind, position int) (models.ResourceGroup, error) {
	r.XML.PlaylistsChanged = true
	var container *[]Node
	var folderNode *Node
	
	if parentID == "" {
		container = &r.XML.Playlists.Node.Node
	} else {
		folderNode = r.findOrCreateContainer(parentID)
		container = &folderNode.Node
	}

	nodeType := 1 // Playlist
	if groupKind == models.GroupKindFolder {
		nodeType = 0
	}

	node := Node{
		Name:    name,
		Type:    int32(nodeType),
		KeyType: PtrInt32(0),
	}

	if nodeType == 1 {
		node.Entries = PtrInt32(0)
	} else {
		node.Count = PtrInt32(0)
	}

	if position < 0 || position >= len(*container) {
		*container = append(*container, node)
	} else {
		*container = append((*container)[:position+1], (*container)[position:]...)
		(*container)[position] = node
	}

	if folderNode != nil {
		if folderNode.Count == nil {
			folderNode.Count = PtrInt32(1)
		} else {
			*folderNode.Count++
		}
	}
	return ToNeutralGroup(node, parentID), nil
}

func (r *Library) DeleteGroup(groupID string, groupKind models.GroupKind) error {
	r.XML.PlaylistsChanged = true
	nodeType := int32(1)
	if groupKind == models.GroupKindFolder {
		nodeType = 0
	}

	_, parentNode, parentSlice, idx := r.findGroupInTree(&r.XML.Playlists.Node.Node, nil, groupID, nodeType)
	if idx == -1 {
		return fmt.Errorf("group not found")
	}
	*parentSlice = append((*parentSlice)[:idx], (*parentSlice)[idx+1:]...)
	if parentNode != nil && parentNode.Count != nil && *parentNode.Count > 0 {
		*parentNode.Count--
	}
	return nil
}

func (r *Library) AddTracks(groupID string, trackIDs []string) (int, error) {
	r.XML.PlaylistsChanged = true
	node, _, _, _ := r.findGroupInTree(&r.XML.Playlists.Node.Node, nil, groupID, 1)
	if node == nil {
		return 0, fmt.Errorf("playlist not found")
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
	node.Entries = PtrInt32(int32(len(node.TRACK)))
	return added, nil
}

func (r *Library) RemoveTracks(groupID string, trackIDs []string) (int, error) {
	r.XML.PlaylistsChanged = true
	node, _, _, _ := r.findGroupInTree(&r.XML.Playlists.Node.Node, nil, groupID, 1)
	if node == nil {
		return 0, fmt.Errorf("playlist not found")
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
	node.Entries = PtrInt32(int32(len(node.TRACK)))
	return before - len(node.TRACK), nil
}

func (r *Library) UpdateGroup(groupID string, trackIDs []string) error {
	r.XML.PlaylistsChanged = true
	node, _, _, _ := r.findGroupInTree(&r.XML.Playlists.Node.Node, nil, groupID, 1)
	if node == nil {
		return fmt.Errorf("playlist not found")
	}
	node.TRACK = nil
	for _, id := range trackIDs {
		node.TRACK = append(node.TRACK, struct {
			Key string `xml:"Key,attr"`
		}{Key: id})
	}
	node.Entries = PtrInt32(int32(len(trackIDs)))
	return nil
}

func (r *Library) RenameGroup(groupID, newName string, groupKind models.GroupKind) error {
	r.XML.PlaylistsChanged = true
	nodeType := int32(1)
	if groupKind == models.GroupKindFolder {
		nodeType = 0
	}
	node, _, _, _ := r.findGroupInTree(&r.XML.Playlists.Node.Node, nil, groupID, nodeType)
	if node == nil {
		return fmt.Errorf("group not found")
	}
	node.Name = newName
	return nil
}

func (r *Library) MoveGroup(groupID string, groupKind models.GroupKind, targetParentID string) error {
	r.XML.PlaylistsChanged = true
	nodeType := int32(1)
	if groupKind == models.GroupKindFolder {
		nodeType = 0
	}
	node, parentNode, parentSlice, idx := r.findGroupInTree(&r.XML.Playlists.Node.Node, nil, groupID, nodeType)
	if node == nil {
		return fmt.Errorf("group not found")
	}

	moved := *node
	*parentSlice = append((*parentSlice)[:idx], (*parentSlice)[idx+1:]...)
	if parentNode != nil && parentNode.Count != nil && *parentNode.Count > 0 {
		*parentNode.Count--
	}

	target := r.findOrCreateContainer(targetParentID)
	target.Node = append(target.Node, moved)
	if target.Count == nil {
		target.Count = PtrInt32(1)
	} else {
		*target.Count++
	}
	return nil
}

func (r *Library) UpdateMetadata(matches []models.MetadataMatch, fields []string) error {
	UpdateBatch(r.XML, matches, fields)
	return nil
}

func (r *Library) Save(path string) error {
	return WriteRekordboxLibrary(path, r.XML)
}

func (r *Library) findOrCreateContainer(name string) *Node {
	nodes := &r.XML.Playlists.Node.Node
	for i := range *nodes {
		if (*nodes)[i].Name == name && (*nodes)[i].Type == 0 {
			return &(*nodes)[i]
		}
	}
	*nodes = append(*nodes, Node{
		Name:  name,
		Type:  0,
		Count: PtrInt32(0),
	})
	return &(*nodes)[len(*nodes)-1]
}

func (r *Library) findGroupInTree(nodes *[]Node, parent *Node, name string, nodeType int32) (*Node, *Node, *[]Node, int) {
	for i := range *nodes {
		n := &(*nodes)[i]
		if n.Name == name && n.Type == nodeType {
			return n, parent, nodes, i
		}
		if len(n.Node) > 0 {
			if found, foundParent, foundSlice, idx := r.findGroupInTree(&n.Node, n, name, nodeType); found != nil {
				return found, foundParent, foundSlice, idx
			}
		}
	}
	return nil, nil, nil, -1
}

// NewLibrary creates a new Library wrapper for a Rekordbox XML.
func NewLibrary(xml *RekordboxLibraryXML) *Library {
	return &Library{XML: xml}
}
