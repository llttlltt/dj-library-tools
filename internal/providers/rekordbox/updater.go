package rekordbox

import (
	"fmt"
	"github.com/llttlltt/dj-library-tools/internal/core/models"
)

// UpdateBatch applies specific field updates to the provided XML library.
func UpdateBatch(rbXML *RekordboxLibraryXML, matches []models.MetadataMatch, fields []string) int {
	fieldMap := make(map[string]bool)
	for _, f := range fields {
		fieldMap[f] = true
	}

	updateCount := 0
	for _, match := range matches {
		for i := range rbXML.Collection.TRACK {
			target := &rbXML.Collection.TRACK[i]
			if fmt.Sprintf("%d", target.TrackID) == match.Target.ID {
				if fieldMap["beatgrids"] {
					if rt, ok := match.Source.ImplementationState.(Track); ok {
						target.Tempo = rt.Tempo
					}
				}
				if fieldMap["rating"] {
					target.Rating = int32(match.Source.Rating)
				}
				if fieldMap["comment"] {
					target.Comments = match.Source.Comment
				}
				if fieldMap["genre"] {
					target.Genre = match.Source.Genre
				}
				if fieldMap["label"] {
					target.Label = match.Source.Label
				}
				if fieldMap["key"] {
					target.Tonality = match.Source.Key
				}
				if fieldMap["bpm"] {
					target.AverageBpm = fmt.Sprintf("%.2f", match.Source.BPM)
				}
				
				updateCount++
				break
			}
		}
	}

	return updateCount
}
