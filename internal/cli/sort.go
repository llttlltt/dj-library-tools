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

	sort.Slice(tracks, func(i, j int) bool {
		f := strings.ToLower(field)
		switch f {
		case "bpm":
			bi, _ := strconv.ParseFloat(tracks[i].AverageBpm, 64)
			bj, _ := strconv.ParseFloat(tracks[j].AverageBpm, 64)
			return bi < bj
		case "artist":
			return strings.ToLower(tracks[i].Artist) < strings.ToLower(tracks[j].Artist)
		case "title":
			return strings.ToLower(tracks[i].Name) < strings.ToLower(tracks[j].Name)
		case "album":
			return strings.ToLower(tracks[i].Album) < strings.ToLower(tracks[j].Album)
		case "key":
			return tracks[i].Tonality < tracks[j].Tonality
		case "rating":
			return tracks[i].Rating < tracks[j].Rating
		case "playcount":
			return tracks[i].PlayCount < tracks[j].PlayCount
		case "added":
			return tracks[i].DateAdded < tracks[j].DateAdded
		default:
			return false
		}
	})
}

func sortNodes(results []provider.NodeResult, field string) {
	if field == "" {
		return
	}

	sort.Slice(results, func(i, j int) bool {
		f := strings.ToLower(field)
		switch f {
		case "name":
			return strings.ToLower(results[i].Name) < strings.ToLower(results[j].Name)
		case "entries", "count":
			return results[i].Entries < results[j].Entries
		default:
			return false
		}
	})
}
