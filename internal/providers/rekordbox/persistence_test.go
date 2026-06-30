package rekordbox

import (
	"os"
	"testing"
	"time"

	"github.com/llttlltt/dj-library-tools/internal/services/library"
	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/llttlltt/dj-library-tools/internal/providers"
)

func TestStatefulPersistence(t *testing.T) {
	// 1. Create a temp Rekordbox XML
	content := `<?xml version="1.0" encoding="UTF-8"?>
<DJ_PLAYLISTS Version="1.0.0">
  <PRODUCT Name="rekordbox" Version="6.6.4" Company="AlphaTheta"/>
  <COLLECTION Entries="1">
    <TRACK TrackID="1" Name="Original Title" Artist="Original Artist" Location="file://localhost/path/to/track.mp3">
    </TRACK>
  </COLLECTION>
  <PLAYLISTS>
    <NODE Type="0" Name="ROOT" Count="0">
    </NODE>
  </PLAYLISTS>
</DJ_PLAYLISTS>`
	
	tmpFile, err := os.CreateTemp("", "rekordbox_*.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// 2. Load provider
	rbXML, _ := ReadRekordboxLibrary(tmpFile.Name())
	eng := library.NewEngine(NewLibrary(rbXML))
	p := NewRekordboxProviderWithXML(eng, rbXML, tmpFile.Name())

	ctx := provider.ExecutionContext{
		Apply:    false,
		Verbose:  true,
		Feedback: &mockFeedback{},
	}

	initialStat, _ := os.Stat(tmpFile.Name())
	initialModTime := initialStat.ModTime()

	// 3. Perform a mutation WITHOUT apply
	// We'll use UpdateBatch (which we just refactored)
	matches := []models.MetadataMatch{
		{
			Source: models.Track{
				ID:      "1",
				Comment: "New Comment",
			},
			Target: models.Track{
				ID: "1",
			},
		},
	}
	
	err = p.Tracks().UpdateBatch(ctx, matches, []string{"comment"})
	if err != nil {
		t.Fatalf("UpdateBatch failed: %v", err)
	}


	// 4. Assert in-memory state changed
	if rbXML.Collection.TRACK[0].Comments != "New Comment" {
		t.Errorf("expected in-memory comment to be updated, got %q", rbXML.Collection.TRACK[0].Comments)
	}

	// 5. Assert disk NOT changed (check modtime or content)
	// Give it a tiny bit of time to ensure if a write happened, modtime would definitely change
	time.Sleep(10 * time.Millisecond)
	newStat, _ := os.Stat(tmpFile.Name())
	if !newStat.ModTime().Equal(initialModTime) {
		t.Errorf("expected file modtime to be unchanged, but it changed from %v to %v", initialModTime, newStat.ModTime())
	}

	// 6. Perform a mutation WITH apply (explicit Save)
	ctx.Apply = true
	err = p.System().Save(ctx, "")
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// 7. Assert disk CHANGED
	newContent, _ := os.ReadFile(tmpFile.Name())
	if !contains(string(newContent), "New Comment") {
		t.Errorf("expected file content to contain 'New Comment', but it didn't")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (len(substr) == 0 || (func() bool {
		for i := 0; i <= len(s)-len(substr); i++ {
			if s[i:i+len(substr)] == substr {
				return true
			}
		}
		return false
	})())
}

type mockFeedback struct{}

func (f *mockFeedback) OnPreview(message string)           {}
func (f *mockFeedback) OnSuccess(message string)           {}
func (f *mockFeedback) OnWarning(message string)           {}
func (f *mockFeedback) OnTable(headers []string, rows [][]string) {}
