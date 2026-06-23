package engine

import (
	"testing"

	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

// makeLibraryWithPlaylists builds a synthetic library with:
//
//	ROOT
//	  └── My Sets (folder)
//	        ├── Summer Vibes (playlist, 2 tracks: IDs 1,2)
//	        └── Winter Warmers (playlist, 1 track: ID 3)
//	  └── All Techno (playlist, 2 tracks: IDs 1,3) — root level
func makeLibraryWithPlaylists() *rekordbox.RekordboxLibraryXML {
	return &rekordbox.RekordboxLibraryXML{
		Collection: rekordbox.Collection{
			Entries: 3,
			TRACK: []rekordbox.Track{
				{TrackID: 1, Name: "Track One", Genre: "Techno", AverageBpm: "128.0"},
				{TrackID: 2, Name: "Track Two", Genre: "House", AverageBpm: "122.0"},
				{TrackID: 3, Name: "Track Three", Genre: "Techno", AverageBpm: "132.0"},
			},
		},
		Playlists: rekordbox.Playlists{
			Node: rekordbox.RootNode{
				BaseNode: rekordbox.BaseNode{Type: 0, Name: "ROOT"},
				Count:    2,
				Node: []rekordbox.Node{
					{
						BaseNode: rekordbox.BaseNode{Type: 0, Name: "My Sets"},
						Count:    2,
						Node: []rekordbox.Node{
							{
								BaseNode: rekordbox.BaseNode{Type: 1, Name: "Summer Vibes"},
								KeyType:  0,
								Entries:  2,
								TRACK: []struct {
									Key string `xml:"Key,attr"`
								}{{Key: "1"}, {Key: "2"}},
							},
							{
								BaseNode: rekordbox.BaseNode{Type: 1, Name: "Winter Warmers"},
								KeyType:  0,
								Entries:  1,
								TRACK: []struct {
									Key string `xml:"Key,attr"`
								}{{Key: "3"}},
							},
						},
					},
					{
						BaseNode: rekordbox.BaseNode{Type: 1, Name: "All Techno"},
						KeyType:  0,
						Entries:  2,
						TRACK: []struct {
							Key string `xml:"Key,attr"`
						}{{Key: "1"}, {Key: "3"}},
					},
				},
			},
		},
	}
}

func TestLsPlaylists(t *testing.T) {
	lib := makeLibraryWithPlaylists()
	eng := NewEngine(lib)

	tests := []struct {
		name      string
		query     string
		wantCount int
		wantNames []string
	}{
		{
			name:      "name substring",
			query:     "name:Summer",
			wantCount: 1,
			wantNames: []string{"Summer Vibes"},
		},
		{
			name:      "folder filter",
			query:     "folder:My Sets",
			wantCount: 2,
			wantNames: []string{"Summer Vibes", "Winter Warmers"},
		},
		{
			name:      "root level playlist by name",
			query:     "name:All Techno",
			wantCount: 1,
			wantNames: []string{"All Techno"},
		},
		{
			name:      "empty query returns all playlists",
			query:     "",
			wantCount: 3,
		},
		{
			name:      "no match returns empty",
			query:     "name:Nonexistent",
			wantCount: 0,
		},
		{
			name:      "entries range",
			query:     "entries:2..99",
			wantCount: 2,
			wantNames: []string{"Summer Vibes", "All Techno"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := eng.LsPlaylists(tt.query)
			if err != nil {
				t.Fatalf("LsPlaylists(%q) error: %v", tt.query, err)
			}
			if len(results) != tt.wantCount {
				t.Errorf("got %d results, want %d", len(results), tt.wantCount)
			}
			for _, wantName := range tt.wantNames {
				found := false
				for _, r := range results {
					if r.Node.Name == wantName {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected result with name %q not found", wantName)
				}
			}
		})
	}
}

func TestLsFolders(t *testing.T) {
	lib := makeLibraryWithPlaylists()
	eng := NewEngine(lib)

	tests := []struct {
		name      string
		query     string
		wantCount int
		wantNames []string
	}{
		{
			name:      "name match",
			query:     "name:My Sets",
			wantCount: 1,
			wantNames: []string{"My Sets"},
		},
		{
			name:      "no match",
			query:     "name:Nonexistent",
			wantCount: 0,
		},
		{
			name:      "empty query returns all folders",
			query:     "",
			wantCount: 1,
			wantNames: []string{"My Sets"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := eng.LsFolders(tt.query)
			if err != nil {
				t.Fatalf("LsFolders(%q) error: %v", tt.query, err)
			}
			if len(results) != tt.wantCount {
				t.Errorf("got %d results, want %d", len(results), tt.wantCount)
			}
			for _, wantName := range tt.wantNames {
				found := false
				for _, r := range results {
					if r.Node.Name == wantName {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected result with name %q not found", wantName)
				}
			}
		})
	}
}

func TestLsPlaylists_ParentFolderPropagated(t *testing.T) {
	lib := makeLibraryWithPlaylists()
	eng := NewEngine(lib)

	results, err := eng.LsPlaylists("folder:My Sets")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.ParentFolder != "My Sets" {
			t.Errorf("playlist %q: ParentFolder = %q, want %q", r.Node.Name, r.ParentFolder, "My Sets")
		}
	}

	rootResults, err := eng.LsPlaylists("name:All Techno")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rootResults) != 1 {
		t.Fatalf("expected 1 result, got %d", len(rootResults))
	}
	if rootResults[0].ParentFolder != "" {
		t.Errorf("root-level playlist ParentFolder = %q, want empty string", rootResults[0].ParentFolder)
	}
}
