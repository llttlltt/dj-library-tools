package library

import (
	"testing"

	"github.com/llttlltt/dj-library-tools/internal/models"
)

func TestEngine_Ls(t *testing.T) {
	mock := &MockLibrary{
		Tracks: []models.Track{
			{ID: "1", Title: "Oceans", Artist: "Four Tet"},
		},
	}
	eng := NewEngine(mock)

	matched, err := eng.Ls("title:Oceans", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(matched) != 1 {
		t.Errorf("got %d, want 1", len(matched))
	}
}
