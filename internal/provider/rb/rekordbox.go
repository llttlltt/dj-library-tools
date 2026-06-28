package rb

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/library"
	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/provider/factory"
	"github.com/llttlltt/dj-library-tools/internal/query"
	"github.com/llttlltt/dj-library-tools/internal/rekordbox"
	"github.com/llttlltt/dj-library-tools/internal/sync"
	"github.com/llttlltt/dj-library-tools/internal/utils"
	"github.com/fatih/color"
)

func init() {
	factory.Register("rb", func(opts factory.ProviderOptions) (provider.Provider, error) {
		if opts.FilePath == "" {
			return nil, fmt.Errorf("rekordbox XML library required via --file flag")
		}
		rbXML, err := rekordbox.ReadRekordboxLibrary(opts.FilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read rekordbox library: %w", err)
		}
		eng := library.NewEngine(rekordbox.NewLibrary(rbXML))
		return NewRekordboxProviderWithXML(eng, rbXML, opts.FilePath), nil
	})
	factory.Register("rekordbox", func(opts factory.ProviderOptions) (provider.Provider, error) {
		return factory.NewProvider("rb", opts)
	})
}

type RekordboxProvider struct {
	engine *library.Engine
	path   string
	rbXML  *rekordbox.RekordboxLibraryXML
}

func NewRekordboxProvider(eng *library.Engine, path string) *RekordboxProvider {
	return &RekordboxProvider{engine: eng, path: path}
}

func NewRekordboxProviderWithXML(eng *library.Engine, rbXML *rekordbox.RekordboxLibraryXML, path string) *RekordboxProvider {
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

// --- Track Service ---

type rekordboxTrackService struct{ *RekordboxProvider }

func (s *rekordboxTrackService) List(ctx provider.ExecutionContext, query string) ([]models.Track, error) {
	return s.engine.Ls(query, s)
}

func (s *rekordboxTrackService) Update(ctx provider.ExecutionContext, query string, changes map[string]string) (int, error) {
	return 0, fmt.Errorf("update tracks by query not yet implemented for rekordbox")
}

func (s *rekordboxTrackService) UpdateBatch(ctx provider.ExecutionContext, matches []models.MetadataMatch, fields []string) error {
	rbXML, err := rekordbox.ReadRekordboxLibrary(s.path)
	if err != nil {
		return err
	}

	fieldMap := make(map[string]bool)
	for _, f := range fields {
		fieldMap[f] = true
	}

	blue := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	updateCount := 0
	for _, match := range matches {
		for i := range rbXML.Collection.TRACK {
			target := &rbXML.Collection.TRACK[i]
			if fmt.Sprintf("%d", target.TrackID) == match.Target.ID {
				if ctx.Verbose {
					fmt.Printf("[%s]\n", blue(match.Target.Artist+" - "+match.Target.Title))
				}

				if fieldMap["beatgrids"] {
					if rt, ok := match.Source.Raw.(rekordbox.Track); ok {
						if ctx.Verbose {
							fmt.Printf("  %s Beatgrids: %d -> %d\n", yellow("~"), len(target.Tempo), len(rt.Tempo))
						}
						target.Tempo = rt.Tempo
					}
				}
				if fieldMap["rating"] {
					if ctx.Verbose {
						fmt.Printf("  %s Rating: %d -> %d\n", yellow("~"), target.Rating, match.Source.Rating)
					}
					target.Rating = int32(match.Source.Rating)
				}
				if fieldMap["comment"] {
					if ctx.Verbose {
						fmt.Printf("  %s Comment: %q -> %q\n", yellow("~"), target.Comments, match.Source.Comment)
					}
					target.Comments = match.Source.Comment
				}
				if fieldMap["genre"] {
					if ctx.Verbose {
						fmt.Printf("  %s Genre: %q -> %q\n", yellow("~"), target.Genre, match.Source.Genre)
					}
					target.Genre = match.Source.Genre
				}
				if fieldMap["label"] {
					if ctx.Verbose {
						fmt.Printf("  %s Label: %q -> %q\n", yellow("~"), target.Label, match.Source.Label)
					}
					target.Label = match.Source.Label
				}
				if fieldMap["key"] {
					if ctx.Verbose {
						fmt.Printf("  %s Key: %q -> %q\n", yellow("~"), target.Tonality, match.Source.Key)
					}
					target.Tonality = match.Source.Key
				}
				if fieldMap["bpm"] {
					newBpm := fmt.Sprintf("%.2f", match.Source.BPM)
					if ctx.Verbose {
						fmt.Printf("  %s BPM: %q -> %q\n", yellow("~"), target.AverageBpm, newBpm)
					}
					target.AverageBpm = newBpm
				}
				
				updateCount++
				break
			}
		}
	}

	if ctx.Verbose {
		fmt.Printf("\nSuccessfully updated %d tracks.\n", updateCount)
	}

	if !ctx.DryRun {
		return rekordbox.WriteRekordboxLibrary(s.path, rbXML)
	}
	return nil
}

func (s *rekordboxTrackService) Delete(ctx provider.ExecutionContext, query string) (int, error) {
	return 0, fmt.Errorf("track deletion not supported by rekordbox provider")
}

func (s *rekordboxTrackService) Groups() provider.TrackGroupService {
	return s
}

func (s *rekordboxTrackService) Add(ctx provider.ExecutionContext, tracks []models.Track, target models.ResourceGroup) (int, error) {
	if target.Type == models.GroupTypeFolder {
		return 0, fmt.Errorf("cannot add tracks to folder %q (rekordbox tracks must live in playlists)", target.Name)
	}
	var ids []string
	for _, t := range tracks {
		ids = append(ids, t.ID)
		if ctx.Verbose {
			fmt.Printf("  %s %s - %s\n", color.GreenString("+"), t.Artist, t.Title)
		}
	}
	return s.engine.Library.(library.WritableLibrary).AddTracks(target.ID, ids)
}

func (s *rekordboxTrackService) Remove(ctx provider.ExecutionContext, tracks []models.Track, group models.ResourceGroup) (int, error) {
	var ids []string
	for _, t := range tracks {
		ids = append(ids, t.ID)
		if ctx.Verbose {
			fmt.Printf("  %s %s - %s\n", color.RedString("-"), t.Artist, t.Title)
		}
	}
	return s.engine.Library.(library.WritableLibrary).RemoveTracks(group.ID, ids)
}

func (s *rekordboxTrackService) Move(ctx provider.ExecutionContext, tracks []models.Track, from models.ResourceGroup, to models.ResourceGroup) (int, error) {
	if to.Type == models.GroupTypeFolder {
		return 0, fmt.Errorf("cannot move tracks into folder %q (rekordbox tracks must live in playlists)", to.Name)
	}
	
	ids := make([]string, len(tracks))
	for i, t := range tracks {
		ids[i] = t.ID
	}

	if ctx.Verbose {
		fmt.Printf("Moving %d tracks from %q to %q\n", len(tracks), from.Name, to.Name)
	}

	if !ctx.DryRun {
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

	return len(tracks), nil
}

func (s *rekordboxTrackService) Sort(ctx provider.ExecutionContext, tracks []models.Track, field string) {
	utils.SortTracksAgnostic(tracks, field)
}

// --- Group Service ---

type rekordboxGroupService struct{ *RekordboxProvider }

func (s *rekordboxGroupService) List(ctx provider.ExecutionContext, query string) ([]models.ResourceGroup, error) {
	return s.engine.LsGroups(query)
}

func (s *rekordboxGroupService) Create(ctx provider.ExecutionContext, parent models.ResourceGroup, name string, groupType models.GroupType, position int) (models.ResourceGroup, error) {
	if parent.Name != "" && parent.Type == models.GroupTypePlaylist {
		return models.ResourceGroup{}, fmt.Errorf("cannot create group inside playlist %q (containers must live in folders)", parent.Name)
	}
	if ctx.Verbose {
		fmt.Printf("Creating %s %q in %q\n", groupType, name, parent.Name)
	}
	return s.engine.Library.(library.WritableLibrary).CreateGroup(parent.ID, name, groupType, position)
}

func (s *rekordboxGroupService) Update(ctx provider.ExecutionContext, group models.ResourceGroup, newName string, newParent *models.ResourceGroup) error {
	if newName != "" {
		if err := s.engine.Library.(library.WritableLibrary).RenameGroup(group.ID, newName, group.Type); err != nil {
			return err
		}
	}
	if newParent != nil {
		if newParent.Type == models.GroupTypePlaylist {
			return fmt.Errorf("cannot move group into playlist %q (containers must live in folders)", newParent.Name)
		}
		return s.engine.Library.(library.WritableLibrary).MoveGroup(group.ID, group.Type, newParent.ID)
	}
	return nil
}

func (s *rekordboxGroupService) Delete(ctx provider.ExecutionContext, group models.ResourceGroup) error {
	if ctx.Verbose {
		fmt.Printf("Deleting %s %q\n", group.GetKind(), group.Name)
	}
	return s.engine.Library.(library.WritableLibrary).DeleteGroup(group.ID, group.Type)
}

func (s *rekordboxGroupService) Sort(ctx provider.ExecutionContext, groups []models.ResourceGroup, field string) {
	utils.SortGroupsAgnostic(groups, field)
}

// --- System Service ---

type rekordboxSystemService struct{ *RekordboxProvider }

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
	return []string{"bpm", "key", "rating", "comment", "genre", "label", "color", "beatgrids"}
}

func (s *rekordboxSystemService) SupportedResources() []string {
	return []string{"tracks", "playlists", "folders"}
}

func (s *rekordboxSystemService) Save(ctx provider.ExecutionContext, path string) error {
	if path == "" {
		path = s.path
	}
	if ctx.Verbose {
		fmt.Printf("Saving Rekordbox XML to %s\n", path)
	}
	return s.engine.Library.(library.WritableLibrary).Save(path)
}

func (s *rekordboxSystemService) Fix(ctx provider.ExecutionContext, resource string, query string) error {
	return nil
}

func (s *rekordboxSystemService) Sync(ctx provider.ExecutionContext, tracks []models.Track, sourceQuery string, targetQuery string, options provider.SyncOptions) error {
	var rbLib *rekordbox.Library
	if s.rbXML != nil {
		rbLib = rekordbox.NewLibrary(s.rbXML)
	} else {
		rbXML, err := rekordbox.ReadRekordboxLibrary(s.path)
		if err != nil {
			return err
		}
		rbLib = rekordbox.NewLibrary(rbXML)
	}

	orch := sync.NewOrchestrator(rbLib, ctx.DryRun, ctx.Verbose)

	err := orch.SyncToLibrary(tracks, sourceQuery, targetQuery, sync.SyncOptions{
		ExportDest:   options.ExportDest,
		ExportFormat: options.ExportFormat,
		PathMaps:     options.PathMaps,
	}, options.AppendOnly)
	
	if err != nil {
		return err
	}

	if s.rbXML == nil && !ctx.DryRun {
		return rbLib.Save(s.path)
	}
	return nil
}

func (s *rekordboxSystemService) Identify(name string, groupType models.GroupType) string {
	return name
}

// --- Custom Matching (Internal for Engine) ---

func (p *RekordboxProvider) CustomMatch(track models.Track, field string, op query.Operator, value string) bool {
	target := strings.ToLower(value)

	if rt, ok := track.Raw.(rekordbox.Track); ok {
		if field == "hotcues" {
			for _, pm := range rt.PositionMark {
				if pm.Num == -1 { continue }
				if p.matchCueMetadata(pm, target, op) { return true }
			}
		} else if field == "memorycues" {
			for _, pm := range rt.PositionMark {
				if pm.Num != -1 { continue }
				if p.matchCueMetadata(pm, target, op) { return true }
			}
		}
	}
	return false
}

func (p *RekordboxProvider) matchCueMetadata(pm rekordbox.PositionMark, target string, op query.Operator) bool {
	if op == query.OpExact {
		if strings.EqualFold(pm.Name, target) { return true }
	} else if strings.Contains(strings.ToLower(pm.Name), target) {
		return true
	}

	colorName := strings.ToLower(p.getHotCueColorName(pm))
	if op == query.OpExact {
		if colorName == target { return true }
	} else if strings.Contains(colorName, target) {
		return true
	}

	return false
}

func (p *RekordboxProvider) getHotCueColorName(pm rekordbox.PositionMark) string {
	rgb := fmt.Sprintf("%02X%02X%02X", pm.Red, pm.Green, pm.Blue)
	switch rgb {
	case "E62828": return "red"
	case "DE44CF": return "hotpink"
	case "FFFF00", "B4BE04", "C3AF04": return "yellow"
	case "28E214", "10B176": return "green"
	case "00E0FF", "50B4FF": return "aqua"
	case "305AFF", "6473FF": return "blue"
	case "B432FF", "AA72FF": return "purple"
	case "E0641B", "FFA500": return "orange"
	}
	return ""
}
