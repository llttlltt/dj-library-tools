package rekordbox

import (
	"context"
	"fmt"

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
		var res FixData
		var err error
		switch fixType {
		case provider.FixDuplicates:
			for _, target := range targets {
				if target == "members" {
					res, err = FixDuplicateMembers(ctx, s.rbXML, selection, ectx.Apply)
				}
				if target == "tracks" {
					res, err = FixDuplicateTracks(ctx, s.rbXML, selection, ectx.Apply)
				}
				if err != nil { return totalAffected, err }
				s.reportFix(ectx, res, string(fixType), target)
				totalAffected += res.TotalApplied
			}
		case provider.FixMetadata:
			res, err = FixMetadataNormalization(ctx, s.rbXML, selection, targets, ectx.Apply)
			if err != nil { return totalAffected, err }
			s.reportFix(ectx, res, string(fixType), "metadata")
			totalAffected += res.TotalApplied
		case provider.FixPaths:
			res, err = FixPaths(ctx, s.rbXML, selection, targets, ectx.Apply)
			if err != nil { return totalAffected, err }
			s.reportFix(ectx, res, string(fixType), "paths")
			totalAffected += res.TotalApplied
		}
	}

	return totalAffected, nil
}

func (s *rekordboxSystemService) reportFix(ectx provider.ExecutionContext, res FixData, fixType, target string) {
	if res.TotalIdentified == 0 {
		return
	}
	for _, d := range res.Details {
		if d.Table != nil && ectx.Verbose {
			ectx.Feedback.OnTable(d.Table.Headers, d.Table.Rows)
		}
		if d.From != "" && ectx.Verbose {
			fmt.Printf("Relocating %s\n  From: %s\n  To:   %s\n", d.TrackName, d.From, d.To)
		}
	}

	fmt.Printf("Identified %d %s/%s issues.\n", res.TotalIdentified, fixType, target)
	if ectx.Apply {
		fmt.Printf("Applied %d repairs.\n", res.TotalApplied)
	} else {
		fmt.Println("Run with --apply to commit these repairs.")
	}
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
