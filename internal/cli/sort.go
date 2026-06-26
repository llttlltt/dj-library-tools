package cli

import (
	"sort"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/models"
)

func sortTracks(tracks []models.Track, field string) {
	if field == "" {
		return
	}

	desc := false
	if strings.HasPrefix(field, "-") {
		desc = true
		field = field[1:]
	}

	sort.Slice(tracks, func(i, j int) bool {
		f := strings.ToLower(field)
		res := false
		switch f {
		case "bpm":
			res = tracks[i].BPM < tracks[j].BPM
		case "artist":
			res = strings.ToLower(tracks[i].Artist) < strings.ToLower(tracks[j].Artist)
		case "title":
			res = strings.ToLower(tracks[i].Title) < strings.ToLower(tracks[j].Title)
		case "album":
			res = strings.ToLower(tracks[i].Album) < strings.ToLower(tracks[j].Album)
		case "key":
			res = tracks[i].Key < tracks[j].Key
		case "rating":
			res = tracks[i].Rating < tracks[j].Rating
		case "added":
			// Date sorting requires more care, but for now we'll do string
			res = false
		default:
			return false
		}

		if desc {
			return !res
		}
		return res
	})
}

func sortNodes(results []models.Node, field string) {
	if field == "" {
		return
	}

	desc := false
	if strings.HasPrefix(field, "-") {
		desc = true
		field = field[1:]
	}

	sort.Slice(results, func(i, j int) bool {
		f := strings.ToLower(field)
		res := false
		switch f {
		case "name":
			res = strings.ToLower(results[i].Name) < strings.ToLower(results[j].Name)
		case "entries", "count":
			res = results[i].Entries < results[j].Entries
		default:
			return false
		}

		if desc {
			return !res
		}
		return res
	})
}
