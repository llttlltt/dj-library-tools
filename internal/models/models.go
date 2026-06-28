package models

// Resource is the interface for any item in a music library (Track, Playlist, Folder).
type Resource interface {
	GetID() string
	GetName() string
	GetKind() string // "track", "playlist", "folder"
}

// Track is the provider-neutral representation of a music track.
type Track struct {
	ID            string
	Title         string
	Artist        string
	Album         string
	BPM           float64
	Key           string
	Genre         string
	Comment       string
	Label         string
	Year          int
	Location      string
	Duration      int
	Rating        int // 0-5
	Plays         int
	DateAdded     string
	DateModified  string
	Color         string
	Bitrate       int
	SampleRate    int
	Size          int64
	Remixer       string
	Mix           string

	// Advanced DJ Metadata
	HotCues       int
	MemoryCues    int
	BeatgridCount int

	Raw interface{}
}

func (t Track) GetID() string   { return t.ID }
func (t Track) GetName() string { return t.Title }
func (t Track) GetKind() string { return "track" }

// ResourceGroup represents a container like a playlist or folder.
type ResourceGroup struct {
	ID           string
	Name         string
	Items        int
	ParentFolder string
	Type         GroupType
	Raw          interface{}
}

type GroupType int

const (
	GroupTypeFolder   GroupType = 0
	GroupTypePlaylist GroupType = 1
)

func (n ResourceGroup) GetID() string   { return n.ID }
func (n ResourceGroup) GetName() string { return n.Name }
func (n ResourceGroup) GetKind() string {
	if n.Type == GroupTypeFolder {
		return "folder"
	}
	return "playlist"
}
