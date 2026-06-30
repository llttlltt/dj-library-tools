package util

import (
	"sort"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/core/models"
)

func SortTracksAgnostic(tracks []models.Track, field string) {
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
		default:
			return false
		}

		if desc {
			return !res
		}
		return res
	})
}

func SortGroupsAgnostic(groups []models.ResourceGroup, field string) {
	if field == "" {
		return
	}

	desc := false
	if strings.HasPrefix(field, "-") {
		desc = true
		field = field[1:]
	}

	sort.Slice(groups, func(i, j int) bool {
		f := strings.ToLower(field)
		res := false
		switch f {
		case "name":
			res = strings.ToLower(groups[i].Name) < strings.ToLower(groups[j].Name)
		case "items":
			res = groups[i].Items < groups[j].Items
		default:
			return false
		}

		if desc {
			return !res
		}
		return res
	})
}
