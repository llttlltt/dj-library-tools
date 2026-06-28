package provider

import (
	"fmt"
	"github.com/llttlltt/dj-library-tools/internal/djerr"
	"github.com/llttlltt/dj-library-tools/internal/models"
)

// GatedProvider wraps a provider and enforces capabilities at the interface level.
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

// gatedTrackService enforces read-only status for track updates.
type gatedTrackService struct {
	base TrackService
	caps ProviderCapabilities
}

func (s *gatedTrackService) List(ctx ExecutionContext, query string) ([]models.Track, error) {
	return s.base.List(ctx, query)
}

func (s *gatedTrackService) Update(ctx ExecutionContext, query string, changes map[string]string) (int, error) {
	if !s.caps.CanWrite {
		return 0, djerr.ErrReadOnly
	}
	return s.base.Update(ctx, query, changes)
}

func (s *gatedTrackService) UpdateBatch(ctx ExecutionContext, matches []models.MetadataMatch, fields []string) error {
	if !s.caps.CanUpdateMetadata {
		return djerr.ErrReadOnly
	}
	return s.base.UpdateBatch(ctx, matches, fields)
}

func (s *gatedTrackService) Delete(ctx ExecutionContext, query string) (int, error) {
	if !s.caps.CanWrite {
		return 0, djerr.ErrReadOnly
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
	return s.base.Add(ctx, tracks, target)
}

func (s *gatedTrackGroupService) Remove(ctx ExecutionContext, tracks []models.Track, group models.ResourceGroup) (int, error) {
	if !s.caps.CanWrite {
		return 0, djerr.ErrReadOnly
	}
	return s.base.Remove(ctx, tracks, group)
}

func (s *gatedTrackGroupService) Move(ctx ExecutionContext, tracks []models.Track, from models.ResourceGroup, to models.ResourceGroup) (int, error) {
	if !s.caps.CanWrite {
		return 0, djerr.ErrReadOnly
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
		return models.ResourceGroup{}, fmt.Errorf("group management is not supported by this provider")
	}
	return s.base.Create(ctx, parent, name, gt, pos)
}

func (s *gatedGroupService) Update(ctx ExecutionContext, group models.ResourceGroup, newName string, newParent *models.ResourceGroup) error {
	if !s.caps.CanManageGroups {
		return djerr.ErrReadOnly
	}
	return s.base.Update(ctx, group, newName, newParent)
}

func (s *gatedGroupService) Delete(ctx ExecutionContext, group models.ResourceGroup) error {
	if !s.caps.CanManageGroups {
		return djerr.ErrReadOnly
	}
	return s.base.Delete(ctx, group)
}

func (s *gatedGroupService) Sort(ctx ExecutionContext, groups []models.ResourceGroup, field string) {
	s.base.Sort(ctx, groups, field)
}
