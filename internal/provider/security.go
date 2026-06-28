package provider

import (
	"fmt"
	"github.com/llttlltt/dj-library-tools/internal/djerr"
	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/fatih/color"
)

// GatedProvider wraps a provider and enforces capabilities and safety at the interface level.
type GatedProvider struct {
	Base Provider
}

func (p *GatedProvider) Name() string { return p.Base.Name() }

func (p *GatedProvider) Tracks() TrackService {
	return &gatedTrackService{
		base: p.Base.Tracks(),
		caps: p.Base.System().Capabilities(),
	}
}

func (p *GatedProvider) Groups() GroupService {
	return &gatedGroupService{
		base: p.Base.Groups(),
		caps: p.Base.System().Capabilities(),
	}
}

func (p *GatedProvider) System() SystemService {
	return p.Base.System()
}

// gatedTrackService enforces read-only status and apply-mode for track updates.
type gatedTrackService struct {
	base TrackService
	caps ProviderCapabilities
}

func (s *gatedTrackService) List(ctx ExecutionContext, query string) ([]models.Track, error) {
	return s.base.List(ctx, query)
}

func (s *gatedTrackService) Update(ctx ExecutionContext, query string, changes map[string]string) (int, error) {
	if !s.caps.CanUpdateMetadata {
		return 0, djerr.ErrReadOnly
	}
	if !ctx.Apply {
		fmt.Printf("[%s] Would update tracks matching %q with %v\n", color.YellowString("Preview"), query, changes)
		return 0, nil
	}
	return s.base.Update(ctx, query, changes)
}

func (s *gatedTrackService) UpdateBatch(ctx ExecutionContext, matches []models.MetadataMatch, fields []string) error {
	if !s.caps.CanUpdateMetadata {
		return djerr.ErrReadOnly
	}
	if !ctx.Apply {
		fmt.Printf("[%s] Would sync metadata fields %v for %d matched tracks\n", color.YellowString("Preview"), fields, len(matches))
		return nil
	}
	return s.base.UpdateBatch(ctx, matches, fields)
}

func (s *gatedTrackService) Delete(ctx ExecutionContext, query string) (int, error) {
	if !s.caps.CanWrite {
		return 0, djerr.ErrReadOnly
	}
	if !ctx.Apply {
		fmt.Printf("[%s] Would delete tracks matching %q\n", color.RedString("Preview"), query)
		return 0, nil
	}
	return s.base.Delete(ctx, query)
}

func (s *gatedTrackService) Groups() TrackGroupService {
	return &gatedTrackGroupService{base: s.base.Groups(), caps: s.caps}
}

func (s *gatedTrackService) Sort(ctx ExecutionContext, tracks []models.Track, field string) {
	s.base.Sort(ctx, tracks, field)
}

type gatedTrackGroupService struct {
	base TrackGroupService
	caps ProviderCapabilities
}

func (s *gatedTrackGroupService) Add(ctx ExecutionContext, tracks []models.Track, target models.ResourceGroup) (int, error) {
	if !s.caps.CanWrite {
		return 0, djerr.ErrReadOnly
	}
	if !ctx.Apply {
		fmt.Printf("[%s] Would add %d tracks to %q\n", color.GreenString("Preview"), len(tracks), target.Name)
		return len(tracks), nil
	}
	return s.base.Add(ctx, tracks, target)
}

func (s *gatedTrackGroupService) Remove(ctx ExecutionContext, tracks []models.Track, group models.ResourceGroup) (int, error) {
	if !s.caps.CanWrite {
		return 0, djerr.ErrReadOnly
	}
	if !ctx.Apply {
		fmt.Printf("[%s] Would remove %d tracks from %q\n", color.RedString("Preview"), len(tracks), group.Name)
		return len(tracks), nil
	}
	return s.base.Remove(ctx, tracks, group)
}

func (s *gatedTrackGroupService) Move(ctx ExecutionContext, tracks []models.Track, from models.ResourceGroup, to models.ResourceGroup) (int, error) {
	if !s.caps.CanWrite {
		return 0, djerr.ErrReadOnly
	}
	if !ctx.Apply {
		fmt.Printf("[%s] Would move %d tracks from %q to %q\n", color.YellowString("Preview"), len(tracks), from.Name, to.Name)
		return len(tracks), nil
	}
	return s.base.Move(ctx, tracks, from, to)
}

// gatedGroupService enforces management rules for playlists and folders.
type gatedGroupService struct {
	base GroupService
	caps ProviderCapabilities
}

func (s *gatedGroupService) List(ctx ExecutionContext, query string) ([]models.ResourceGroup, error) {
	return s.base.List(ctx, query)
}

func (s *gatedGroupService) Create(ctx ExecutionContext, parent models.ResourceGroup, name string, gt models.GroupKind, pos int) (models.ResourceGroup, error) {
	if !s.caps.CanManageGroups {
		return models.ResourceGroup{}, djerr.ErrUnsupportedResource
	}
	if !ctx.Apply {
		fmt.Printf("[%s] Would create %s %q in folder %q\n", color.GreenString("Preview"), gt, name, parent.Name)
		return models.ResourceGroup{Name: name, Kind: gt}, nil
	}
	return s.base.Create(ctx, parent, name, gt, pos)
}

func (s *gatedGroupService) Update(ctx ExecutionContext, group models.ResourceGroup, newName string, newParent *models.ResourceGroup) error {
	if !s.caps.CanManageGroups {
		return djerr.ErrReadOnly
	}
	if !ctx.Apply {
		if newName != "" {
			fmt.Printf("[%s] Would rename %s %q to %q\n", color.YellowString("Preview"), group.GetKind(), group.Name, newName)
		}
		if newParent != nil {
			fmt.Printf("[%s] Would move %s %q into folder %q\n", color.YellowString("Preview"), group.GetKind(), group.Name, newParent.Name)
		}
		return nil
	}
	return s.base.Update(ctx, group, newName, newParent)
}

func (s *gatedGroupService) Delete(ctx ExecutionContext, group models.ResourceGroup) error {
	if !s.caps.CanManageGroups {
		return djerr.ErrReadOnly
	}
	if !ctx.Apply {
		fmt.Printf("[%s] Would delete %s %q\n", color.RedString("Preview"), group.GetKind(), group.Name)
		return nil
	}
	return s.base.Delete(ctx, group)
}

func (s *gatedGroupService) Sort(ctx ExecutionContext, groups []models.ResourceGroup, field string) {
	s.base.Sort(ctx, groups, field)
}
