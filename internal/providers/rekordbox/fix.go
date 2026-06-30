package rekordbox

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/llttlltt/dj-library-tools/internal/providers"
)

type FixResult struct {
	TotalAffected int
	Details       []FixDetail
}

type FixDetail struct {
	TrackName string
	Message   string
	From      string
	To        string
	Table     *TableData
}

type FixData struct {
	TotalIdentified int
	TotalApplied    int
	Details         []FixDetail
}

type TableData struct {
	Headers []string
	Rows    [][]string
}

func FixPaths(ctx context.Context, rbXML *RekordboxLibraryXML, selection provider.Selection, targets []string, apply bool) (FixData, error) {
	res := FixData{}
	strategies := make(map[string]bool)
	for _, t := range targets {
		strategies[strings.ToLower(t)] = true
	}

	var toProcess []*Track
	if len(selection.Tracks) > 0 {
		idMap := make(map[string]*Track)
		for i := range rbXML.Collection.TRACK {
			id := fmt.Sprintf("%d", rbXML.Collection.TRACK[i].TrackID)
			idMap[id] = &rbXML.Collection.TRACK[i]
		}
		for _, t := range selection.Tracks {
			if rt, ok := idMap[t.ID]; ok {
				toProcess = append(toProcess, rt)
			}
		}
	} else {
		for i := range rbXML.Collection.TRACK {
			toProcess = append(toProcess, &rbXML.Collection.TRACK[i])
		}
	}

	for _, rt := range toProcess {
		select {
		case <-ctx.Done():
			return res, ctx.Err()
		default:
		}
		originalLocation := rt.Location
		if originalLocation == "" {
			continue
		}

		u, err := url.Parse(originalLocation)
		if err != nil {
			continue
		}
		currentPath := u.Path

		if _, err := os.Stat(currentPath); err == nil {
			continue
		}

		newPath := ""
		if strategies["normalize"] {
			exts := []string{".mp3", ".m4a", ".wav", ".flac", ".aif", ".aiff"}
			base := strings.TrimSuffix(currentPath, filepath.Ext(currentPath))
			for _, ext := range exts {
				testPath := base + ext
				if _, err := os.Stat(testPath); err == nil {
					newPath = testPath
					break
				}
			}
		}

		if newPath != "" {
			res.TotalIdentified++
			res.Details = append(res.Details, FixDetail{
				TrackName: rt.Name,
				From:      currentPath,
				To:        newPath,
			})
			if apply {
				u.Path = newPath
				rt.Location = u.String()
				rbXML.CollectionChanged = true
				res.TotalApplied++
			}
		}
	}
	return res, nil
}

func FixMetadataNormalization(ctx context.Context, rbXML *RekordboxLibraryXML, selection provider.Selection, targets []string, apply bool) (FixData, error) {
	res := FixData{}
	idMap := make(map[string]*Track)
	for i := range rbXML.Collection.TRACK {
		id := fmt.Sprintf("%d", rbXML.Collection.TRACK[i].TrackID)
		idMap[id] = &rbXML.Collection.TRACK[i]
	}

	var toProcess []*Track
	if len(selection.Tracks) > 0 {
		for _, t := range selection.Tracks {
			if rt, ok := idMap[t.ID]; ok {
				toProcess = append(toProcess, rt)
			}
		}
	} else {
		for i := range rbXML.Collection.TRACK {
			toProcess = append(toProcess, &rbXML.Collection.TRACK[i])
		}
	}

	normalizeAll := len(targets) == 0 || (len(targets) == 1 && targets[0] == "all")
	targetMap := make(map[string]bool)
	for _, t := range targets {
		targetMap[strings.ToLower(t)] = true
	}

	for _, rt := range toProcess {
		select {
		case <-ctx.Done():
			return res, ctx.Err()
		default:
		}
		fixed := false
		if normalizeAll || targetMap["artist"] {
			if rt.Artist != strings.TrimSpace(rt.Artist) {
				rt.Artist = strings.TrimSpace(rt.Artist)
				fixed = true
			}
		}
		if normalizeAll || targetMap["title"] {
			if rt.Name != strings.TrimSpace(rt.Name) {
				rt.Name = strings.TrimSpace(rt.Name)
				fixed = true
			}
		}
		if normalizeAll || targetMap["album"] {
			if rt.Album != strings.TrimSpace(rt.Album) {
				rt.Album = strings.TrimSpace(rt.Album)
				fixed = true
			}
		}
		if normalizeAll || targetMap["genre"] {
			if rt.Genre != strings.TrimSpace(rt.Genre) {
				rt.Genre = strings.TrimSpace(rt.Genre)
				fixed = true
			}
		}
		if normalizeAll || targetMap["comments"] {
			lowerComment := strings.ToLower(rt.Comments)
			if lowerComment == "none" || lowerComment == "nil" {
				rt.Comments = ""
				fixed = true
			}
		}

		if fixed {
			res.TotalIdentified++
			res.Details = append(res.Details, FixDetail{TrackName: fmt.Sprintf("%s - %s", rt.Artist, rt.Name)})
			if apply {
				rbXML.CollectionChanged = true
				res.TotalApplied++
			}
		}
	}
	return res, nil
}

func FixDuplicateTracks(ctx context.Context, rbXML *RekordboxLibraryXML, selection provider.Selection, apply bool) (FixData, error) {
	res := FixData{}
	tracks := selection.Tracks
	if len(tracks) == 0 {
		for _, rt := range rbXML.Collection.TRACK {
			tracks = append(tracks, ToNeutralTrack(rt))
		}
	}

	type identity struct {
		artist, title string
		size          int64
	}
	seen := make(map[identity]string)
	toRemove := make(map[string]bool)

	for _, t := range tracks {
		select {
		case <-ctx.Done():
			return res, ctx.Err()
		default:
		}
		id := identity{artist: strings.ToLower(t.Artist), title: strings.ToLower(t.Title), size: t.Size}
		if firstID, dup := seen[id]; dup {
			toRemove[t.ID] = true
			res.Details = append(res.Details, FixDetail{TrackName: fmt.Sprintf("%s - %s (ID: %s, Dupe of: %s)", t.Artist, t.Title, t.ID, firstID)})
		} else {
			seen[id] = t.ID
		}
	}

	res.TotalIdentified = len(toRemove)
	if apply && len(toRemove) > 0 {
		before := len(rbXML.Collection.TRACK)
		newTracks := rbXML.Collection.TRACK[:0]
		for _, rt := range rbXML.Collection.TRACK {
			id := fmt.Sprintf("%d", rt.TrackID)
			if !toRemove[id] {
				newTracks = append(newTracks, rt)
			}
		}
		rbXML.Collection.TRACK = newTracks
		rbXML.Collection.Entries = int32(len(newTracks))
		rbXML.CollectionChanged = true
		RemoveTrackIDsFromPlaylists(&rbXML.Playlists.Node.Node, toRemove, &rbXML.PlaylistsChanged)
		res.TotalApplied = before - len(newTracks)
	}
	return res, nil
}

func RemoveTrackIDsFromPlaylists(nodes *[]Node, ids map[string]bool, changed *bool) {
	for i := range *nodes {
		node := &(*nodes)[i]
		if node.Type == 1 {
			before := len(node.TRACK)
			kept := node.TRACK[:0]
			for _, t := range node.TRACK {
				if !ids[t.Key] {
					kept = append(kept, t)
				}
			}
			if len(kept) != before {
				node.TRACK = kept
				node.Entries = PtrInt32(int32(len(kept)))
				*changed = true
			}
		}
		if len(node.Node) > 0 {
			RemoveTrackIDsFromPlaylists(&node.Node, ids, changed)
		}
	}
}

func FixDuplicateMembers(ctx context.Context, rbXML *RekordboxLibraryXML, selection provider.Selection, apply bool) (FixData, error) {
	res := FixData{}
	trackLookup := make(map[string]models.Track)
	for _, rt := range rbXML.Collection.TRACK {
		t := ToNeutralTrack(rt)
		trackLookup[t.ID] = t
	}

	for _, item := range selection.Items {
		select {
		case <-ctx.Done():
			return res, ctx.Err()
		default:
		}
		group, ok := item.(models.ResourceGroup)
		if !ok || group.Kind != models.GroupKindPlaylist {
			continue
		}

		node, _, _, _ := rbXML.FindGroupInTree(&rbXML.Playlists.Node.Node, nil, group.Name, 1)
		if node == nil {
			continue
		}

		seenAt := make(map[string]int)
		var kept []struct {
			Key string `xml:"Key,attr"`
		}
		removed := 0
		var tableRows [][]string

		for i, t := range node.TRACK {
			pos := i + 1
			if firstPos, dup := seenAt[t.Key]; dup {
				removed++
				artist, title := "[Unknown]", "[Unknown]"
				if tr, ok := trackLookup[t.Key]; ok {
					artist, title = tr.Artist, tr.Title
				}
				tableRows = append(tableRows, []string{fmt.Sprintf("%d, %d", firstPos, pos), t.Key, title, artist})
			} else {
				seenAt[t.Key] = pos
				kept = append(kept, t)
			}
		}

		if removed > 0 {
			res.TotalIdentified += removed
			detail := FixDetail{
				TrackName: group.Name,
				Message:   fmt.Sprintf("%s: %d duplicates found", group.Name, removed),
				Table: &TableData{
					Headers: []string{"pos", "id", "title", "artist"},
					Rows:    tableRows,
				},
			}
			res.Details = append(res.Details, detail)
			if apply {
				node.TRACK = kept
				node.Entries = PtrInt32(int32(len(kept)))
				rbXML.PlaylistsChanged = true
				res.TotalApplied += removed
			}
		}
	}
	return res, nil
}
