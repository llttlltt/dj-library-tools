package provider

import (
	"context"
	"fmt"
	djerrors "github.com/llttlltt/dj-library-tools/internal/core/errors"
	"github.com/llttlltt/dj-library-tools/internal/core/models"
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

func (s *gatedTrackService) List(ctx context.Context, ectx ExecutionContext, query string) ([]models.Track, error) {
	return s.base.List(ctx, ectx, query)
}

func (s *gatedTrackService) Update(ctx context.Context, ectx ExecutionContext, query string, changes map[string]string) (int, error) {
	if !s.caps.CanUpdateMetadata {
		return 0, djerrors.ErrReadOnly
	}
	if !ectx.Apply {
		ectx.Feedback.OnPreview(fmt.Sprintf("update tracks matching %q with %v", query, changes))
		return 0, nil
	}
	return s.base.Update(ctx, ectx, query, changes)
}

func (s *gatedTrackService) UpdateBatch(ctx context.Context, ectx ExecutionContext, matches []models.MetadataMatch, fields []string) error {
	if !s.caps.CanUpdateMetadata {
		return djerrors.ErrReadOnly
	}
	if !ectx.Apply {
		ectx.Feedback.OnPreview(fmt.Sprintf("sync metadata fields %v for %d matched tracks", fields, len(matches)))
		return nil
	}
	return s.base.UpdateBatch(ctx, ectx, matches, fields)
}

func (s *gatedTrackService) Delete(ctx context.Context, ectx ExecutionContext, query string) (int, error) {
	if !s.caps.CanWrite {
		return 0, djerrors.ErrReadOnly
	}
	if !ectx.Apply {
		ectx.Feedback.OnPreview(fmt.Sprintf("delete tracks matching %q", query))
		return 0, nil
	}
	return s.base.Delete(ctx, ectx, query)
}

func (s *gatedTrackService) Groups() TrackGroupService {
	return &gatedTrackGroupService{base: s.base.Groups(), caps: s.caps}
}

func (s *gatedTrackService) Sort(ctx context.Context, ectx ExecutionContext, tracks []models.Track, field string) {
	s.base.Sort(ctx, ectx, tracks, field)
}

type gatedTrackGroupService struct {
	base TrackGroupService
	caps ProviderCapabilities
}

func (s *gatedTrackGroupService) Add(ctx context.Context, ectx ExecutionContext, tracks []models.Track, target models.ResourceGroup) (int, error) {
	if !s.caps.CanWrite {
		return 0, djerrors.ErrReadOnly
	}
	if !ectx.Apply {
		ectx.Feedback.OnPreview(fmt.Sprintf("add %d tracks to %q", len(tracks), target.Name))
		return len(tracks), nil
	}
	return s.base.Add(ctx, ectx, tracks, target)
}

func (s *gatedTrackGroupService) Remove(ctx context.Context, ectx ExecutionContext, tracks []models.Track, group models.ResourceGroup) (int, error) {
	if !s.caps.CanWrite {
		return 0, djerrors.ErrReadOnly
	}
	if !ectx.Apply {
		ectx.Feedback.OnPreview(fmt.Sprintf("remove %d tracks from %q", len(tracks), group.Name))
		return len(tracks), nil
	}
	return s.base.Remove(ctx, ectx, tracks, group)
}

func (s *gatedTrackGroupService) Move(ctx context.Context, ectx ExecutionContext, tracks []models.Track, from models.ResourceGroup, to models.ResourceGroup) (int, error) {
	if !s.caps.CanWrite {
		return 0, djerrors.ErrReadOnly
	}
	if !ectx.Apply {
		ectx.Feedback.OnPreview(fmt.Sprintf("move %d tracks from %q to %q", len(tracks), from.Name, to.Name))
		return len(tracks), nil
	}
	return s.base.Move(ctx, ectx, tracks, from, to)
}

// gatedGroupService enforces management rules for playlists and folders.
type gatedGroupService struct {
	base GroupService
	caps ProviderCapabilities
}

func (s *gatedGroupService) List(ctx context.Context, ectx ExecutionContext, query string) ([]models.ResourceGroup, error) {
	return s.base.List(ctx, ectx, query)
}

func (s *gatedGroupService) Create(ctx context.Context, ectx ExecutionContext, parent models.ResourceGroup, name string, gt models.GroupKind, pos int) (models.ResourceGroup, error) {
	if !s.caps.CanManageGroups {
		return models.ResourceGroup{}, djerrors.ErrUnsupportedResource
	}
	if !ectx.Apply {
		ectx.Feedback.OnPreview(fmt.Sprintf("create %s %q in folder %q", gt, name, parent.Name))
		return models.ResourceGroup{Name: name, Kind: gt}, nil
	}
	return s.base.Create(ctx, ectx, parent, name, gt, pos)
}

func (s *gatedGroupService) Update(ctx context.Context, ectx ExecutionContext, group models.ResourceGroup, newName string, newParent *models.ResourceGroup) error {
	if !s.caps.CanManageGroups {
		return djerrors.ErrReadOnly
	}
	if !ectx.Apply {
		if newName != "" {
			ectx.Feedback.OnPreview(fmt.Sprintf("rename %s %q to %q", group.GetKind(), group.Name, newName))
		}
		if newParent != nil {
			ectx.Feedback.OnPreview(fmt.Sprintf("move %s %q into folder %q", group.GetKind(), group.Name, newParent.Name))
		}
		return nil
	}
	return s.base.Update(ctx, ectx, group, newName, newParent)
}

func (s *gatedGroupService) Delete(ctx context.Context, ectx ExecutionContext, group models.ResourceGroup) error {
	if !s.caps.CanManageGroups {
		return djerrors.ErrReadOnly
	}
	if !ectx.Apply {
		ectx.Feedback.OnPreview(fmt.Sprintf("delete %s %q", group.GetKind(), group.Name))
		return nil
	}
	return s.base.Delete(ctx, ectx, group)
}

func (s *gatedGroupService) Sort(ctx context.Context, ectx ExecutionContext, groups []models.ResourceGroup, field string) {
	s.base.Sort(ctx, ectx, groups, field)
}
