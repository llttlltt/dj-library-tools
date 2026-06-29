package rekordbox

import (
	"encoding/xml"
)

// These structs represent the XML structure of a Rekordbox library.
// The definitions are based on the Rekordbox developer documentation:
// https://rekordbox.com/en/support/developer

type RekordboxLibraryXML struct {
	XMLName    xml.Name   `xml:"DJ_PLAYLISTS"`
	Version    string     `xml:"Version,attr"`
	Product    Product    `xml:"PRODUCT"`
	Collection Collection `xml:"COLLECTION"`
	Playlists  Playlists  `xml:"PLAYLISTS"`

	// Internal state for surgical saving
	OriginalRaw       []byte     `xml:"-"`
	CollectionChanged bool       `xml:"-"`
	Format            *XMLFormat `xml:"-"`
	PlaylistsChanged  bool       `xml:"-"`
}

type Product struct {
	XMLName xml.Name `xml:"PRODUCT"`
	Name    string   `xml:"Name,attr"`
	Version string   `xml:"Version,attr"`
	Company string   `xml:"Company,attr"`
}

type Collection struct {
	XMLName xml.Name `xml:"COLLECTION"`
	Entries int32    `xml:"Entries,attr"`
	TRACK   []Track  `xml:"TRACK"`
}

type Track struct {
	XMLName xml.Name `xml:"TRACK"`

	// --- Standard Rekordbox XML Specification (v1.0.0) ---
	TrackID      int     `xml:"TrackID,attr"`
	Name         string  `xml:"Name,attr"`
	Artist       string  `xml:"Artist,attr"`
	Composer     string  `xml:"Composer,attr"`
	Album        string  `xml:"Album,attr"`
	Grouping     string  `xml:"Grouping,attr"`
	Genre        string  `xml:"Genre,attr"`
	Kind         string  `xml:"Kind,attr"`
	Size         int64   `xml:"Size,attr"`
	TotalTime    int32   `xml:"TotalTime,attr"`
	DiscNumber   int32   `xml:"DiscNumber,attr"`
	TrackNumber  int32   `xml:"TrackNumber,attr"`
	Year         int32   `xml:"Year,attr"`
	AverageBpm   string  `xml:"AverageBpm,attr"`
	DateModified string  `xml:"DateModified,attr,omitempty"`
	DateAdded    string  `xml:"DateAdded,attr"`
	BitRate      int32   `xml:"BitRate,attr"`
	SampleRate   float64 `xml:"SampleRate,attr"`
	Comments     string  `xml:"Comments,attr"`
	PlayCount    int32   `xml:"PlayCount,attr"`
	LastPlayed   string  `xml:"LastPlayed,attr,omitempty"`
	Rating       int32   `xml:"Rating,attr"`
	Location     string  `xml:"Location,attr"`
	Remixer      string  `xml:"Remixer,attr"`
	Tonality     string  `xml:"Tonality,attr"`
	Label        string  `xml:"Label,attr"`
	Mix          string  `xml:"Mix,attr"`
	Colour       string  `xml:"Colour,attr,omitempty"`

	Tempo        []Tempo        `xml:"TEMPO"`
	PositionMark []PositionMark `xml:"POSITION_MARK"`
}

type Tempo struct {
	XMLName xml.Name `xml:"TEMPO"`

	// --- Standard Rekordbox XML Specification (v1.0.0) ---
	Inizio  string `xml:"Inizio,attr"`
	Bpm     string `xml:"Bpm,attr"`
	Metro   string `xml:"Metro,attr"`
	Battito int32  `xml:"Battito,attr"`
}

type PositionMark struct {
	XMLName xml.Name `xml:"POSITION_MARK"`

	// --- Standard Rekordbox XML Specification (v1.0.0) ---
	Name  string `xml:"Name,attr"`
	Type  int32  `xml:"Type,attr"`
	Start string `xml:"Start,attr"`
	End   string `xml:"End,attr,omitempty"`
	Num   int32  `xml:"Num,attr"`

	// --- Extensions (Undocumented Rekordbox Attributes) ---
	Red   int32 `xml:"Red,attr,omitempty"`
	Green int32 `xml:"Green,attr,omitempty"`
	Blue  int32 `xml:"Blue,attr,omitempty"`
}

type Playlists struct {
	XMLName xml.Name `xml:"PLAYLISTS"`
	Node    RootNode `xml:"NODE"`
}

type BaseNode struct {
	XMLName xml.Name `xml:"NODE"`
}

type RootNode struct {
	BaseNode
	XMLName xml.Name `xml:"NODE"`
	Type    int32    `xml:"Type,attr"`
	Name    string   `xml:"Name,attr"`
	Count   int32    `xml:"Count,attr"`
	Node    []Node   `xml:"NODE"`
}

type Node struct {
	BaseNode
	XMLName xml.Name `xml:"NODE"`
	Name    string   `xml:"Name,attr"`
	Type    int32    `xml:"Type,attr"`
	// Count, Entries, and KeyType are optional: folders use Count, playlists use
	// KeyType+Entries. Using *int32 with omitempty ensures zero-valued absent
	// attributes are never written (a nil pointer is omitted; &0 is written as
	// "0", which is correct for KeyType=0).
	Count   *int32   `xml:"Count,attr,omitempty"`
	Entries *int32   `xml:"Entries,attr,omitempty"`
	KeyType *int32   `xml:"KeyType,attr,omitempty"`
	Node    []Node   `xml:"NODE"`
	TRACK   []struct {
		Key string `xml:"Key,attr"`
	} `xml:"TRACK"`
}

// PtrInt32 returns a pointer to the given int32 value.
func PtrInt32(v int32) *int32 { return &v }

func (r *RekordboxLibraryXML) FindGroupInTree(nodes *[]Node, parent *Node, name string, nodeType int32) (*Node, *Node, *[]Node, int) {
	for i := range *nodes {
		n := &(*nodes)[i]
		if n.Name == name && n.Type == nodeType {
			return n, parent, nodes, i
		}
		if len(n.Node) > 0 {
			if found, foundParent, foundSlice, idx := r.FindGroupInTree(&n.Node, n, name, nodeType); found != nil {
				return found, foundParent, foundSlice, idx
			}
		}
	}
	return nil, nil, nil, -1
}

// DerefInt32 safely dereferences an *int32, returning 0 for a nil pointer.
func DerefInt32(p *int32) int32 {
	if p == nil {
		return 0
	}
	return *p
}
