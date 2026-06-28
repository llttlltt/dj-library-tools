package rb

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/library"
	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/query"
	"github.com/llttlltt/dj-library-tools/internal/rekordbox"
	"github.com/llttlltt/dj-library-tools/internal/sync"
)

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

func (p *RekordboxProvider) GetTracks(query string) ([]models.Track, error) {
	return p.Engine.Ls(query, p)
}

func (p *RekordboxProvider) GetPlaylists(query string) ([]models.ResourceGroup, error) {
	return p.Engine.LsPlaylists(query)
}

func (p *RekordboxProvider) GetFolders(query string) ([]models.ResourceGroup, error) {
	return p.Engine.LsFolders(query)
}

func (p *RekordboxProvider) CanTranscode() bool {
	return true
}

func (p *RekordboxProvider) AddTracks(target models.ResourceGroup, tracks []models.Track) (int, error) {
	var ids []string
	for _, t := range tracks {
		ids = append(ids, t.ID)
	}
	_, added := p.Engine.Library.(library.WritableLibrary).AddTracksToPlaylist(target.Name, ids)
	return added, nil
}

func (p *RekordboxProvider) RemoveTracks(target models.ResourceGroup, tracks []models.Track) (int, error) {
	var ids []string
	for _, t := range tracks {
		ids = append(ids, t.ID)
	}
	_, removed := p.Engine.Library.(library.WritableLibrary).RemoveTracksFromPlaylist(target.Name, ids)
	return removed, nil
}

func (p *RekordboxProvider) CreateNode(parent models.ResourceGroup, name string, nodeType int) (models.ResourceGroup, error) {
	if nodeType == 0 {
		p.Engine.Library.(library.WritableLibrary).CreateFolder(parent.Name, name, -1)
	} else {
		p.Engine.Library.(library.WritableLibrary).AddPlaylist(parent.Name, name, nil, -1)
	}
	return models.ResourceGroup{Name: name, Type: models.GroupType(nodeType)}, nil
}

func (p *RekordboxProvider) DeleteNode(node models.ResourceGroup) error {
	p.Engine.Library.(library.WritableLibrary).RemoveNode(node.Name, int32(node.Type))
	return nil
}

func (p *RekordboxProvider) RenameNode(node models.ResourceGroup, newName string) error {
	p.Engine.Library.(library.WritableLibrary).RenameNode(node.Name, newName, int32(node.Type))
	return nil
}

func (p *RekordboxProvider) MoveNode(node models.ResourceGroup, targetParent models.ResourceGroup) error {
	p.Engine.Library.(library.WritableLibrary).MoveNode(node.Name, int32(node.Type), targetParent.Name)
	return nil
}

func (p *RekordboxProvider) Save(path string) error {
	if path == "" {
		path = p.path
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

func (p *RekordboxProvider) Sync(tracks []models.Track, sourceQuery string, targetQuery string, options provider.SyncOptions) error {
	var rbLib *RekordboxLibrary
	if p.rbXML != nil {
		rbLib = NewRekordboxLibrary(p.rbXML)
	} else {
		rbXML, err := rekordbox.ReadRekordboxLibrary(p.path)
		if err != nil {
			return err
		}
		rbLib = NewRekordboxLibrary(rbXML)
	}

	orch := sync.NewOrchestrator(rbLib, false, false)

	err := orch.SyncToLibrary(tracks, sourceQuery, targetQuery, sync.SyncOptions{
		ExportDest:   options.ExportDest,
		ExportFormat: options.ExportFormat,
		PathMaps:     options.PathMaps,
	}, options.AppendOnly)
	
	if err != nil {
		return err
	}

	if p.rbXML == nil {
		return rbLib.Save(p.path)
	}
	return nil
}
