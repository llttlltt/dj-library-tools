package cli

import (
	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/resolver"
)

func sortTracks(src *resolver.Selection, tracks []models.Track, field string) {
	if field == "" {
		return
	}
	src.Provider.Tracks().Sort(getExecContext(), tracks, field)
}

func sortGroups(src *resolver.Selection, groups []models.ResourceGroup, field string) {
	if field == "" {
		return
	}
	src.Provider.Groups().Sort(getExecContext(), groups, field)
}
