package rekordbox

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/services/library"
	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/llttlltt/dj-library-tools/internal/providers"
	"github.com/llttlltt/dj-library-tools/internal/providers/factory"
	"github.com/llttlltt/dj-library-tools/internal/services/sync"
	"github.com/llttlltt/dj-library-tools/internal/core/util"
)

func init() {
	factory.Register("rb", func(opts factory.ProviderOptions) (provider.Provider, error) {
		if opts.FilePath == "" {
			return nil, fmt.Errorf("rekordbox XML library required via --file flag")
		}
		rbXML, err := ReadRekordboxLibrary(opts.FilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read rekordbox library: %w", err)
		}
		eng := library.NewEngine(NewLibrary(rbXML))
		return NewRekordboxProviderWithXML(eng, rbXML, opts.FilePath), nil
	})
	factory.Register("rekordbox", func(opts factory.ProviderOptions) (provider.Provider, error) {
		return factory.NewProvider("rb", opts)
	})
}

type RekordboxProvider struct {
	engine *library.Engine
	path   string
	rbXML  *RekordboxLibraryXML
}

func NewRekordboxProvider(eng *library.Engine, path string) *RekordboxProvider {
	return &RekordboxProvider{engine: eng, path: path}
}

func NewRekordboxProviderWithXML(eng *library.Engine, rbXML *RekordboxLibraryXML, path string) *RekordboxProvider {
	return &RekordboxProvider{engine: eng, rbXML: rbXML, path: path}
}

func (p *RekordboxProvider) Name() string { return "rb" }

func (p *RekordboxProvider) Tracks() provider.TrackService {
	return &rekordboxTrackService{p}
}

func (p *RekordboxProvider) Groups() provider.GroupService {
	return &rekordboxGroupService{p}
}

func (p *RekordboxProvider) System() provider.SystemService {
	return &rekordboxSystemService{p}
}

type rekordboxSystemService struct{ *RekordboxProvider }

// --- Track Service ---

type rekordboxTrackService struct{ *RekordboxProvider }

func (s *rekordboxTrackService) List(ctx context.Context, ectx provider.ExecutionContext, query string) ([]models.Track, error) {
	return s.engine.Ls(query, nil)
}

func (s *rekordboxTrackService) Update(ctx context.Context, ectx provider.ExecutionContext, query string, changes map[string]string) (int, error) {
	return 0, fmt.Errorf("update tracks by query not yet implemented for rekordbox")
}

func (s *rekordboxTrackService) UpdateBatch(ctx context.Context, ectx provider.ExecutionContext, matches []models.MetadataMatch, fields []string) error {
	if s.rbXML == nil {
		return fmt.Errorf("rekordbox library not loaded")
	}

	count := UpdateBatch(s.rbXML, matches, fields)

	if ectx.Verbose {
		fmt.Printf("\nSuccessfully updated %d tracks.\n", count)
	}

	s.rbXML.CollectionChanged = true
	return nil
}

func (s *rekordboxTrackService) Delete(ctx context.Context, ectx provider.ExecutionContext, query string) (int, error) {
	return 0, fmt.Errorf("track deletion not supported by rekordbox provider")
}

func (s *rekordboxTrackService) Groups() provider.TrackGroupService {
	return s
}

func (s *rekordboxTrackService) Add(ctx context.Context, ectx provider.ExecutionContext, tracks []models.Track, target models.ResourceGroup) (int, error) {
	if target.Kind == models.GroupKindFolder {
		return 0, fmt.Errorf("cannot add tracks to folder %q (rekordbox tracks must live in playlists)", target.Name)
	}
	var ids []string
	for _, t := range tracks {
		ids = append(ids, t.ID)
	}
	return s.engine.Library.(library.WritableLibrary).AddTracks(target.ID, ids)
}

func (s *rekordboxTrackService) Remove(ctx context.Context, ectx provider.ExecutionContext, tracks []models.Track, group models.ResourceGroup) (int, error) {
	var ids []string
	for _, t := range tracks {
		ids = append(ids, t.ID)
	}
	return s.engine.Library.(library.WritableLibrary).RemoveTracks(group.ID, ids)
}

func (s *rekordboxTrackService) Move(ctx context.Context, ectx provider.ExecutionContext, tracks []models.Track, from models.ResourceGroup, to models.ResourceGroup) (int, error) {
	if to.Kind == models.GroupKindFolder {
		return 0, fmt.Errorf("cannot move tracks into folder %q (rekordbox tracks must live in playlists)", to.Name)
	}
	
	ids := make([]string, len(tracks))
	for i, t := range tracks {
		ids[i] = t.ID
	}

	if ectx.Verbose {
		fmt.Printf("Moving %d tracks from %q to %q\n", len(tracks), from.Name, to.Name)
	}

	added, err := s.engine.Library.(library.WritableLibrary).AddTracks(to.ID, ids)
	if err != nil {
		return 0, err
	}
	removed, err := s.engine.Library.(library.WritableLibrary).RemoveTracks(from.ID, ids)
	if err != nil {
		return added, err
	}
	return removed, nil 
}

func (s *rekordboxTrackService) Sort(ctx context.Context, ectx provider.ExecutionContext, tracks []models.Track, field string) {
	util.SortTracksAgnostic(tracks, field)
}

// --- Group Service ---

type rekordboxGroupService struct{ *RekordboxProvider }

func (s *rekordboxGroupService) List(ctx context.Context, ectx provider.ExecutionContext, query string) ([]models.ResourceGroup, error) {
	return s.engine.LsGroups(query)
}

func (s *rekordboxGroupService) Create(ctx context.Context, ectx provider.ExecutionContext, parent models.ResourceGroup, name string, groupKind models.GroupKind, position int) (models.ResourceGroup, error) {
	if parent.Name != "" && parent.Kind == models.GroupKindPlaylist {
		return models.ResourceGroup{}, fmt.Errorf("cannot create group inside playlist %q (containers must live in folders)", parent.Name)
	}
	if ectx.Verbose {
		fmt.Printf("Creating %s %q in %q\n", groupKind, name, parent.Name)
	}
	return s.engine.Library.(library.WritableLibrary).CreateGroup(parent.ID, name, groupKind, position)
}

func (s *rekordboxGroupService) Update(ctx context.Context, ectx provider.ExecutionContext, group models.ResourceGroup, newName string, newParent *models.ResourceGroup) error {
	if newName != "" {
		if err := s.engine.Library.(library.WritableLibrary).RenameGroup(group.ID, newName, group.Kind); err != nil {
			return err
		}
	}
	if newParent != nil {
		if newParent.Kind == models.GroupKindPlaylist {
			return fmt.Errorf("cannot move group into playlist %q (containers must live in folders)", newParent.Name)
		}
		return s.engine.Library.(library.WritableLibrary).MoveGroup(group.ID, group.Kind, newParent.ID)
	}
	return nil
}

func (s *rekordboxGroupService) Delete(ctx context.Context, ectx provider.ExecutionContext, group models.ResourceGroup) error {
	if ectx.Verbose {
		fmt.Printf("Deleting %s %q\n", group.GetKind(), group.Name)
	}
	return s.engine.Library.(library.WritableLibrary).DeleteGroup(group.ID, group.Kind)
}

func (s *rekordboxGroupService) Sort(ctx context.Context, ectx provider.ExecutionContext, groups []models.ResourceGroup, field string) {
	util.SortGroupsAgnostic(groups, field)
}

// --- System Service ---

func (s *rekordboxSystemService) Capabilities() provider.ProviderCapabilities {
	return provider.ProviderCapabilities{
		CanWrite:          true,
		CanManageGroups:   true,
		CanUpdateMetadata: true,
		SupportsCues:      true,
		SupportsBeatgrids: true,
		IsFileBased:       true,
	}
}

func (s *rekordboxSystemService) Containment() provider.ContainmentPolicy {
	return provider.ContainmentPolicy{
		AllowTracksInFolders:    false,
		AllowFoldersInPlaylists: false,
		AllowNestedFolders:      true,
	}
}

func (s *rekordboxSystemService) MetadataCapabilities() []string {
	return provider.ResolveAvailableFields(s.Capabilities())
}

func (s *rekordboxSystemService) SupportedResources() []string {
	return []string{"tracks", "playlists", "folders"}
}

func (s *rekordboxSystemService) TableHeaders() []string {
	return []string{"bpm", "key", "artist", "title"}
}

func (s *rekordboxSystemService) Save(ctx context.Context, ectx provider.ExecutionContext, path string) error {
	if path == "" {
		path = s.path
	}
	if ectx.Verbose {
		fmt.Printf("Saving Rekordbox XML to %s\n", path)
	}
	return s.engine.Library.(library.WritableLibrary).Save(path)
}

func (s *rekordboxSystemService) Fix(ctx context.Context, ectx provider.ExecutionContext, selection provider.Selection, options provider.FixOptions) (int, error) {
	totalAffected := 0

	for fixType, targets := range options.Actions {
		switch fixType {
		case provider.FixDuplicates:
			for _, target := range targets {
				if target == "members" {
					affected, err := s.fixDuplicateMembers(ctx, ectx, selection)
					if err != nil {
						return totalAffected, err
					}
					totalAffected += affected
				}
				if target == "tracks" {
					affected, err := s.fixDuplicateTracks(ctx, ectx, selection)
					if err != nil {
						return totalAffected, err
					}
					totalAffected += affected
				}
			}
		case provider.FixMetadata:
			affected, err := s.fixMetadataNormalization(ctx, ectx, selection, targets)
			if err != nil {
				return totalAffected, err
			}
			totalAffected += affected
		case provider.FixPaths:
			affected, err := s.fixPaths(ctx, ectx, selection, targets)
			if err != nil {
				return totalAffected, err
			}
			totalAffected += affected
		}
	}

	return totalAffected, nil
}

func (s *rekordboxSystemService) fixPaths(ctx context.Context, ectx provider.ExecutionContext, selection provider.Selection, targets []string) (int, error) {
	totalRelocated := 0
	
	// Determine strategy
	strategies := make(map[string]bool)
	for _, t := range targets {
		strategies[strings.ToLower(t)] = true
	}

	// 1. Identify tracks to check
	var toProcess []*Track
	if len(selection.Tracks) > 0 {
		idMap := make(map[string]*Track)
		for i := range s.rbXML.Collection.TRACK {
			id := fmt.Sprintf("%d", s.rbXML.Collection.TRACK[i].TrackID)
			idMap[id] = &s.rbXML.Collection.TRACK[i]
		}
		for _, t := range selection.Tracks {
			if rt, ok := idMap[t.ID]; ok {
				toProcess = append(toProcess, rt)
			}
		}
	} else {
		for i := range s.rbXML.Collection.TRACK {
			toProcess = append(toProcess, &s.rbXML.Collection.TRACK[i])
		}
	}

	for _, rt := range toProcess {
		select {
		case <-ctx.Done():
			return totalRelocated, ctx.Err()
		default:
		}
		originalLocation := rt.Location
		if originalLocation == "" {
			continue
		}

		// Decode URI to physical path
		u, err := url.Parse(originalLocation)
		if err != nil {
			continue
		}
		currentPath := u.Path
		
		// Check if file exists
		if _, err := os.Stat(currentPath); err == nil {
			continue // File is already healthy
		}

		// File is missing, try to find it
		newPath := ""

		// Strategy: Normalize (Case sensitivity fix, extension fix)
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
			totalRelocated++
			if ectx.Verbose {
				fmt.Printf("Relocating %s\n  From: %s\n  To:   %s\n", rt.Name, currentPath, newPath)
			}
			if ectx.Apply {
				u.Path = newPath
				rt.Location = u.String()
				s.rbXML.CollectionChanged = true
			}
		} else {
			if ectx.Verbose {
				fmt.Printf("Missing file for %s: %s\n", rt.Name, currentPath)
			}
		}
	}

	if totalRelocated > 0 {
		fmt.Printf("Identified %d missing tracks that can be relocated.\n", totalRelocated)
		if ectx.Apply {
			fmt.Printf("Relocated %d tracks in the library.\n", totalRelocated)
		} else {
			fmt.Println("Run with --apply to commit these relocations.")
		}
	} else {
		fmt.Println("No paths required repair.")
	}

	return totalRelocated, nil
}

func (s *rekordboxSystemService) fixMetadataNormalization(ctx context.Context, ectx provider.ExecutionContext, selection provider.Selection, targets []string) (int, error) {
	totalFixed := 0

	// Map of TrackID to Track pointer for selection lookup
	idMap := make(map[string]*Track)
	for i := range s.rbXML.Collection.TRACK {
		id := fmt.Sprintf("%d", s.rbXML.Collection.TRACK[i].TrackID)
		idMap[id] = &s.rbXML.Collection.TRACK[i]
	}

	// Determine which tracks to process
	var toProcess []*Track
	if len(selection.Tracks) > 0 {
		for _, t := range selection.Tracks {
			if rt, ok := idMap[t.ID]; ok {
				toProcess = append(toProcess, rt)
			}
		}
	} else {
		for i := range s.rbXML.Collection.TRACK {
			toProcess = append(toProcess, &s.rbXML.Collection.TRACK[i])
		}
	}

	// Field targets
	normalizeAll := len(targets) == 0 || (len(targets) == 1 && targets[0] == "all")
	targetMap := make(map[string]bool)
	for _, t := range targets {
		targetMap[strings.ToLower(t)] = true
	}

	for _, rt := range toProcess {
		select {
		case <-ctx.Done():
			return totalFixed, ctx.Err()
		default:
		}
		fixed := false
		
		// 1. Trim Whitespace
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

		// 2. Clear "None" or placeholder values
		if normalizeAll || targetMap["comments"] {
			lowerComment := strings.ToLower(rt.Comments)
			if lowerComment == "none" || lowerComment == "nil" {
				rt.Comments = ""
				fixed = true
			}
		}

		if fixed {
			totalFixed++
			if ectx.Verbose {
				fmt.Printf("Normalized metadata for: %s - %s\n", rt.Artist, rt.Name)
			}
		}
	}

	if totalFixed > 0 {
		fmt.Printf("Identified %d tracks with metadata requiring normalization.\n", totalFixed)
		if ectx.Apply {
			s.rbXML.CollectionChanged = true
			fmt.Printf("Applied normalization to %d tracks.\n", totalFixed)
		} else {
			fmt.Println("Run with --apply to persist changes.")
		}
	} else {
		fmt.Println("No metadata normalization required.")
	}

	return totalFixed, nil
}

func (s *rekordboxSystemService) fixDuplicateTracks(ctx context.Context, ectx provider.ExecutionContext, selection provider.Selection) (int, error) {
	totalRemoved := 0

	// Use the entire collection if no tracks are selected
	tracks := selection.Tracks
	if len(tracks) == 0 {
		for _, rt := range s.rbXML.Collection.TRACK {
			tracks = append(tracks, ToNeutralTrack(rt))
		}
	}

	// 1. Identify duplicates based on Artist, Title, and Size
	// We use Artist+Title+Size as a strong proxy for identical files.
	type identity struct {
		artist, title string
		size          int64
	}

	seen := make(map[identity]string) // identity -> first ID seen
	toRemove := make(map[string]bool)

	for _, t := range tracks {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}
		id := identity{
			artist: strings.ToLower(t.Artist),
			title:  strings.ToLower(t.Title),
			size:   t.Size,
		}

		if firstID, dup := seen[id]; dup {
			toRemove[t.ID] = true
			if ectx.Verbose {
				fmt.Printf("Duplicate track found: %s - %s (ID: %s, Dupe of: %s)\n", t.Artist, t.Title, t.ID, firstID)
			}
		} else {
			seen[id] = t.ID
		}
	}

	if len(toRemove) == 0 {
		fmt.Println("No duplicate tracks found in the master collection.")
		return 0, nil
	}

	fmt.Printf("Found %d duplicate tracks in the master collection.\n", len(toRemove))

	if ectx.Apply {
		// 2. Remove from COLLECTION
		before := len(s.rbXML.Collection.TRACK)
		newTracks := s.rbXML.Collection.TRACK[:0]
		for _, rt := range s.rbXML.Collection.TRACK {
			id := fmt.Sprintf("%d", rt.TrackID)
			if !toRemove[id] {
				newTracks = append(newTracks, rt)
			}
		}
		s.rbXML.Collection.TRACK = newTracks
		s.rbXML.Collection.Entries = int32(len(newTracks))
		s.rbXML.CollectionChanged = true

		// 3. Remove from PLAYLISTS (all memberships)
		s.removeTrackIDsFromPlaylists(&s.rbXML.Playlists.Node.Node, toRemove)

		totalRemoved = before - len(newTracks)
		fmt.Printf("Removed %d duplicate tracks and cleaned up playlist memberships.\n", totalRemoved)
	} else {
		fmt.Println("Run with --apply to remove these duplicates.")
	}

	return len(toRemove), nil
}

func (s *rekordboxSystemService) removeTrackIDsFromPlaylists(nodes *[]Node, ids map[string]bool) {
	for i := range *nodes {
		node := &(*nodes)[i]
		if node.Type == 1 { // Playlist
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
				s.rbXML.PlaylistsChanged = true
			}
		}
		if len(node.Node) > 0 {
			s.removeTrackIDsFromPlaylists(&node.Node, ids)
		}
	}
}
func (s *rekordboxSystemService) fixDuplicateMembers(ctx context.Context, ectx provider.ExecutionContext, selection provider.Selection) (int, error) {
	totalRemoved := 0

	// Build track lookup map if verbose to show track names
	var trackLookup map[string]models.Track
	if ectx.Verbose {
		trackLookup = make(map[string]models.Track)
		for _, rt := range s.rbXML.Collection.TRACK {
			t := ToNeutralTrack(rt)
			trackLookup[t.ID] = t
		}
	}

	// We only care about playlists for duplicate membership fixing
	for _, res := range selection.Items {
		select {
		case <-ctx.Done():
			return totalRemoved, ctx.Err()
		default:
		}
		group, ok := res.(models.ResourceGroup)
		if !ok || group.Kind != models.GroupKindPlaylist {
			continue
		}

		node, _, _, _ := s.rbXML.FindGroupInTree(&s.rbXML.Playlists.Node.Node, nil, group.Name, 1)
		if node == nil {
			continue
		}

		totalTracks := len(node.TRACK)
		seenAt := make(map[string]int) // ID -> 1-based index
		var kept []struct {
			Key string `xml:"Key,attr"`
		}

		removed := 0
		type removedInfo struct {
			pos, firstPos int
			artist, title, id string
		}
		var removedRows []removedInfo

		for i, t := range node.TRACK {
			pos := i + 1
			if firstPos, dup := seenAt[t.Key]; dup {
				removed++
				if ectx.Verbose {
					info := removedInfo{pos: pos, firstPos: firstPos, id: t.Key}
					if track, ok := trackLookup[t.Key]; ok {
						info.artist = track.Artist
						info.title = track.Title
					} else {
						info.artist = "[Unknown]"
						info.title = "[Unknown]"
					}
					removedRows = append(removedRows, info)
				}
			} else {
				seenAt[t.Key] = pos
				kept = append(kept, t)
			}
		}

		if ectx.Verbose && len(removedRows) > 0 {
			headers := []string{"pos", "id", "title", "artist"}
			var rows [][]string
			for _, r := range removedRows {
				rows = append(rows, []string{
					fmt.Sprintf("%d, %d", r.firstPos, r.pos),
					r.id,
					r.title,
					r.artist,
				})
			}
			ectx.Feedback.OnTable(headers, rows)
			fmt.Println()
		}

		fmt.Printf("%s:\n", group.Name)
		fmt.Printf("- Total tracks: %d\n", totalTracks)
		if removed > 0 {
			fmt.Printf("- Duplicate tracks: %d\n", removed)
			fmt.Printf("- Remaining tracks: %d\n", totalTracks-removed)
		}
		fmt.Println()

		if removed > 0 && ectx.Apply {
			node.TRACK = kept
			node.Entries = PtrInt32(int32(len(kept)))
			s.rbXML.PlaylistsChanged = true
		}
		totalRemoved += removed
	}

	return totalRemoved, nil
}

func (s *rekordboxSystemService) Sync(ctx context.Context, ectx provider.ExecutionContext, tracks []models.Track, targetQuery string, options provider.SyncOptions) error {
	if s.rbXML == nil {
		return fmt.Errorf("rekordbox library not loaded")
	}
	rbLib := NewLibrary(s.rbXML)

	err := sync.SyncToLibrary(ctx, rbLib, tracks, targetQuery, sync.SyncOptions{
		ExportDest:     options.ExportDest,
		ExportFormat:   options.ExportFormat,
		PathMaps:       options.PathMaps,
		MetadataFields: options.MetadataFields,
		MatchFields:    options.MatchFields,
	}, ectx.Apply, ectx.Verbose, options.AppendOnly)

	if err != nil {
		return err
	}

	return nil
}

func (s *rekordboxSystemService) Identify(name string, groupType models.GroupKind) string {
	return Identify(name, groupType)
}
