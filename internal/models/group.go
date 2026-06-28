package models

import "strconv"

// ResourceGroup represents a container like a playlist or folder.
type ResourceGroup struct {
	ID           string
	Name         string
	Items        int
	ParentFolder string
	Kind         GroupKind

	ImplementationState interface{}
}

type GroupKind string

const (
	GroupKindFolder   GroupKind = "folder"
	GroupKindPlaylist GroupKind = "playlist"
)

func (g GroupKind) String() string {
	return string(g)
}

func (g ResourceGroup) GetID() string   { return g.ID }
func (g ResourceGroup) GetName() string { return g.Name }
func (g ResourceGroup) GetKind() string { return g.Kind.String() }

// Value returns a string representation of a group property for querying.
func (g ResourceGroup) Value(key string) string {
	switch key {
	case "id":
		return g.ID
	case "name":
		return g.Name
	case "parent", "folder":
		return g.ParentFolder
	case "items":
		return strconv.Itoa(g.Items)
	case "kind":
		return string(g.Kind)
	}
	return ""
}
