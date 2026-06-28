package models

import "strconv"

// ResourceGroup represents a container like a playlist or folder.
type ResourceGroup struct {
	ID           string    `query:"id"`
	Name         string    `query:"name"`
	Items        int       `query:"items,numeric"`
	ParentFolder string    `query:"parent"`
	Kind         GroupKind `query:"kind"`

	ImplementationState interface{} `query:"-"`
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

// GetQueryValue provides a fallback for derived fields in groups.
func (g ResourceGroup) GetQueryValue(field string) (string, bool) {
	switch field {
	case "folder":
		return g.ParentFolder, true
	case "items":
		return strconv.Itoa(g.Items), true
	}
	return "", false
}
