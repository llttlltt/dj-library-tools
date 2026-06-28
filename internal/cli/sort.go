package cli

import (
	"github.com/llttlltt/dj-library-tools/internal/models"
)

func sortTracks(src *Selection, tracks []models.Track, field string) {
	if field == "" {
		return
	}
	src.Provider.SortTracks(getExecContext(), tracks, field)
}

func sortGroups(src *Selection, groups []models.ResourceGroup, field string) {
	if field == "" {
		return
	}
	src.Provider.SortGroups(getExecContext(), groups, field)
}
