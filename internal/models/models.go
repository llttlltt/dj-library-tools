package models

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
	Rating        int     // 0-5
	PlayCount     int
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
	
	Raw           interface{}
}

// Node represents a container like a playlist or folder.
type Node struct {
	Name         string
	Entries      int
	ParentFolder string
	Type         int // 0: Folder, 1: Playlist
	Raw          interface{}
}
