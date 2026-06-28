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
	Engine *library.Engine
	path   string
	rbXML  *rekordbox.RekordboxLibraryXML
}

func (p *RekordboxProvider) Capabilities() provider.ProviderCapabilities {
	return provider.ProviderCapabilities{
		CanWrite:          true,
		CanManageGroups:   true,
		CanUpdateMetadata: true,
		SupportsCues:      true,
		SupportsBeatgrids: true,
		IsFileBased:       true,
	}
}

func (p *RekordboxProvider) GetContainmentPolicy() provider.ContainmentPolicy {
	return provider.ContainmentPolicy{
		AllowTracksInFolders:    false,
		AllowFoldersInPlaylists: false,
		AllowNestedFolders:      true,
	}
}

func (p *RekordboxProvider) Name() string {
	return "rb"
}

func (p *RekordboxProvider) GetResources(ctx provider.ExecutionContext, resource string, query string) ([]models.Resource, error) {
	var items []models.Resource
	switch resource {
	case "tracks":
		tracks, err := p.Engine.Ls(query, p)
		if err != nil { return nil, err }
		for _, t := range tracks { items = append(items, t) }
	case "playlists":
		fullQuery := "type:1"
		if query != "" { fullQuery = "(" + query + ") && type:1" }
		groups, err := p.Engine.LsGroups(fullQuery)
		if err != nil { return nil, err }
		for _, g := range groups { items = append(items, g) }
	case "folders":
		fullQuery := "type:0"
		if query != "" { fullQuery = "(" + query + ") && type:0" }
		groups, err := p.Engine.LsGroups(fullQuery)
		if err != nil { return nil, err }
		for _, g := range groups { items = append(items, g) }
	default:
		return nil, provider.ErrUnsupportedResource
	}
	return items, nil
}

func (p *RekordboxProvider) SortTracks(ctx provider.ExecutionContext, tracks []models.Track, field string) {
	utils.SortTracksAgnostic(tracks, field)
}

func (p *RekordboxProvider) SortGroups(ctx provider.ExecutionContext, groups []models.ResourceGroup, field string) {
	utils.SortGroupsAgnostic(groups, field)
}

func (p *RekordboxProvider) CanTranscode() bool {
	return true
}

func (p *RekordboxProvider) AddTracks(ctx provider.ExecutionContext, target models.ResourceGroup, tracks []models.Track) (int, error) {
	if err := p.ValidateAddTracks(target); err != nil {
		return 0, err
	}
	var ids []string
	for _, t := range tracks {
		ids = append(ids, t.ID)
		if ctx.Verbose {
			fmt.Printf("  %s %s - %s\n", color.GreenString("+"), t.Artist, t.Title)
		}
	}
	return p.Engine.Library.(library.WritableLibrary).AddTracks(target.ID, ids)
}

func (p *RekordboxProvider) RemoveTracks(ctx provider.ExecutionContext, target models.ResourceGroup, tracks []models.Track) (int, error) {
	var ids []string
	for _, t := range tracks {
		ids = append(ids, t.ID)
		if ctx.Verbose {
			fmt.Printf("  %s %s - %s\n", color.RedString("-"), t.Artist, t.Title)
		}
	}
	return p.Engine.Library.(library.WritableLibrary).RemoveTracks(target.ID, ids)
}

func (p *RekordboxProvider) CreateGroup(ctx provider.ExecutionContext, parent models.ResourceGroup, name string, groupType models.GroupType, position int) (models.ResourceGroup, error) {
	if err := p.ValidateCreateGroup(parent, groupType); err != nil {
		return models.ResourceGroup{}, err
	}
	if ctx.Verbose {
		fmt.Printf("Creating %s %q in %q\n", groupType, name, parent.Name)
	}
	return p.Engine.Library.(library.WritableLibrary).CreateGroup(parent.ID, name, groupType, position)
}

func (p *RekordboxProvider) DeleteGroup(ctx provider.ExecutionContext, node models.ResourceGroup) error {
	if ctx.Verbose {
		fmt.Printf("Deleting %s %q\n", node.GetKind(), node.Name)
	}
	return p.Engine.Library.(library.WritableLibrary).DeleteGroup(node.ID, node.Type)
}

func (p *RekordboxProvider) RenameGroup(ctx provider.ExecutionContext, node models.ResourceGroup, newName string, groupType models.GroupType) error {
	if ctx.Verbose {
		fmt.Printf("Renaming %q -> %q\n", node.Name, newName)
	}
	return p.Engine.Library.(library.WritableLibrary).RenameGroup(node.ID, newName, groupType)
}

func (p *RekordboxProvider) MoveGroup(ctx provider.ExecutionContext, node models.ResourceGroup, targetParent models.ResourceGroup) error {
	if err := p.ValidateMoveGroup(node, targetParent); err != nil {
		return err
	}
	if ctx.Verbose {
		fmt.Printf("Moving %q into %q\n", node.Name, targetParent.Name)
	}
	return p.Engine.Library.(library.WritableLibrary).MoveGroup(node.ID, node.Type, targetParent.ID)
}

func (p *RekordboxProvider) MoveTracks(ctx provider.ExecutionContext, source models.ResourceGroup, target models.ResourceGroup, tracks []models.Track) (int, error) {
	if err := p.ValidateAddTracks(target); err != nil {
		return 0, err
	}
	
	ids := make([]string, len(tracks))
	for i, t := range tracks {
		ids[i] = t.ID
	}

	if ctx.Verbose {
		fmt.Printf("Moving %d tracks from %q to %q\n", len(tracks), source.Name, target.Name)
	}

	if !ctx.DryRun {
		added, err := p.Engine.Library.(library.WritableLibrary).AddTracks(target.ID, ids)
		if err != nil {
			return 0, err
		}
		removed, err := p.Engine.Library.(library.WritableLibrary).RemoveTracks(source.ID, ids)
		if err != nil {
			return added, err
		}
		return removed, nil // Return number of tracks successfully moved
	}

	return len(tracks), nil
}

func (p *RekordboxProvider) Save(ctx provider.ExecutionContext, path string) error {
	if path == "" {
		path = p.path
	}
	if ctx.Verbose {
		fmt.Printf("Saving Rekordbox XML to %s\n", path)
	}
	return p.Engine.Library.(library.WritableLibrary).Save(path)
}

func NewRekordboxProvider(eng *library.Engine, path string) *RekordboxProvider {
	return &RekordboxProvider{Engine: eng, path: path}
}

func NewRekordboxProviderWithXML(eng *library.Engine, rbXML *rekordbox.RekordboxLibraryXML, path string) *RekordboxProvider {
	return &RekordboxProvider{Engine: eng, rbXML: rbXML, path: path}
}

func (p *RekordboxProvider) CustomMatch(track models.Track, field string, op query.Operator, value string) bool {
	target := strings.ToLower(value)
	operator := op

	if rt, ok := track.Raw.(rekordbox.Track); ok {
		if field == "hotcues" {
			for _, pm := range rt.PositionMark {
				if pm.Num == -1 {
					continue
				}
				if p.matchCueMetadata(pm, target, operator) {
					return true
				}
			}
		} else if field == "memorycues" {
			for _, pm := range rt.PositionMark {
				if pm.Num != -1 {
					continue
				}
				if p.matchCueMetadata(pm, target, operator) {
					return true
				}
			}
		}
	}
	return false
}

func (p *RekordboxProvider) matchCueMetadata(pm rekordbox.PositionMark, target string, op query.Operator) bool {
	if op == query.OpExact {
		if strings.EqualFold(pm.Name, target) {
			return true
		}
	} else if strings.Contains(strings.ToLower(pm.Name), target) {
		return true
	}

	colorName := strings.ToLower(p.getHotCueColorName(pm))
	if op == query.OpExact {
		if colorName == target {
			return true
		}
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

func (p *RekordboxProvider) GetTrackColorName(hex string) string {
	switch strings.ToUpper(hex) {
	case "0XFF007F": return "pink"
	case "0XFF0000": return "red"
	case "0XFFA500": return "orange"
	case "0XFFFF00": return "yellow"
	case "0X00FF00": return "green"
	case "0X25FDE9": return "aqua"
	case "0X0000FF": return "blue"
	case "0X660099": return "purple"
	}
	return hex
}

func (p *RekordboxProvider) Sync(ctx provider.ExecutionContext, tracks []models.Track, sourceQuery string, targetQuery string, options provider.SyncOptions) error {
	var rbLib *rekordbox.Library
	if p.rbXML != nil {
		rbLib = rekordbox.NewLibrary(p.rbXML)
	} else {
		rbXML, err := rekordbox.ReadRekordboxLibrary(p.path)
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

	if p.rbXML == nil && !ctx.DryRun {
		return rbLib.Save(p.path)
	}
	return nil
}

func (p *RekordboxProvider) UpdateTracks(ctx provider.ExecutionContext, query string, changes map[string]string) (int, error) {
	return 0, fmt.Errorf("metadata modification not yet fully refactored for rekordbox")
}

func (p *RekordboxProvider) ValidateAddTracks(target models.ResourceGroup) error {
	if target.Type == models.GroupTypeFolder {
		return fmt.Errorf("cannot add tracks to folder %q (rekordbox tracks must live in playlists)", target.Name)
	}
	return nil
}

func (p *RekordboxProvider) ValidateMoveGroup(src models.ResourceGroup, target models.ResourceGroup) error {
	if target.Type == models.GroupTypePlaylist {
		return fmt.Errorf("cannot move group into playlist %q (containers must live in folders)", target.Name)
	}
	return nil
}

func (p *RekordboxProvider) ValidateCreateGroup(parent models.ResourceGroup, groupType models.GroupType) error {
	if parent.Name != "" && parent.Type == models.GroupTypePlaylist {
		return fmt.Errorf("cannot create group inside playlist %q (containers must live in folders)", parent.Name)
	}
	return nil
}

func (p *RekordboxProvider) IdentifyGroup(name string, groupType models.GroupType) string {
	return name
}

func (p *RekordboxProvider) SupportedResources() []string {
	return []string{"tracks", "playlists", "folders"}
}

func (p *RekordboxProvider) MetadataCapabilities() []string {
	return []string{"bpm", "key", "rating", "comment", "genre", "label", "color", "beatgrids"}
}

func (p *RekordboxProvider) UpdateMetadata(ctx provider.ExecutionContext, matches []models.MetadataMatch, fields []string) error {
	rbXML, err := rekordbox.ReadRekordboxLibrary(p.path)
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
		return rekordbox.WriteRekordboxLibrary(p.path, rbXML)
	}
	return nil
}

func (p *RekordboxProvider) Fix(ctx provider.ExecutionContext, resource string, query string) error {
	return nil
}
