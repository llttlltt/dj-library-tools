package cli

import (
	"sort"
	"strconv"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

func sortTracks(tracks []rekordbox.Track, field string) {
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
			bi, _ := strconv.ParseFloat(tracks[i].AverageBpm, 64)
			bj, _ := strconv.ParseFloat(tracks[j].AverageBpm, 64)
			res = bi < bj
		case "artist":
			res = strings.ToLower(tracks[i].Artist) < strings.ToLower(tracks[j].Artist)
		case "title":
			res = strings.ToLower(tracks[i].Name) < strings.ToLower(tracks[j].Name)
		case "album":
			res = strings.ToLower(tracks[i].Album) < strings.ToLower(tracks[j].Album)
		case "key":
			res = tracks[i].Tonality < tracks[j].Tonality
		case "rating":
			res = tracks[i].Rating < tracks[j].Rating
		case "playcount":
			res = tracks[i].PlayCount < tracks[j].PlayCount
		case "added":
			res = tracks[i].DateAdded < tracks[j].DateAdded
		default:
			return false
		}

		if desc {
			return !res
		}
		return res
	})
}

func sortNodes(results []provider.NodeResult, field string) {
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
