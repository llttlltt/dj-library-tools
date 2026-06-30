package orchestrator

import (
	"context"
	"fmt"

	"strings"

	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/llttlltt/dj-library-tools/internal/core/query"
	"github.com/llttlltt/dj-library-tools/internal/providers"
	"github.com/llttlltt/dj-library-tools/internal/services/resolver"
)

type Orchestrator struct {
	Feedback             provider.Feedback
	RekordboxPrimaryPath string
}

type Options struct {
	RekordboxPrimaryPath string
}

func New(fb provider.Feedback, opts Options) *Orchestrator {
	if fb == nil {
		fb = provider.NoopFeedback{}
	}
	return &Orchestrator{
		Feedback:             fb,
		RekordboxPrimaryPath: opts.RekordboxPrimaryPath,
	}
}

type RunOptions struct {
	FilePath      string
	Apply         bool
	Verbose       bool
	FilterMissing bool
	FilterExists  bool
}

func (o *Orchestrator) buildResolveOptions(opts RunOptions) resolver.ResolveOptions {
	return resolver.ResolveOptions{
		FilePath:             opts.FilePath,
		RekordboxPrimaryPath: o.RekordboxPrimaryPath,
		FilterMissing:        opts.FilterMissing,
		FilterExists:         opts.FilterExists,
		Apply:                opts.Apply,
		Verbose:              opts.Verbose,
		Feedback:             o.Feedback,
	}
}

func (o *Orchestrator) buildExecContext(opts RunOptions) provider.ExecutionContext {
	return provider.ExecutionContext{
		Apply:    opts.Apply,
		Verbose:  opts.Verbose,
		Feedback: o.Feedback,
	}
}

type StatResult struct {
	Count      int            `json:"Count"`
	AvgBPM     float64        `json:"AvgBPM"`
	Genres     map[string]int `json:"Genres"`
	Labels     map[string]int `json:"Labels"`
	Keys       map[string]int `json:"Keys"`
	Artists    map[string]int `json:"Artists"`
	TotalTempo float64        `json:"TotalTempo"`
}

type ListResult struct {
	Resource       string
	Tracks         []models.Track
	Groups         []models.ResourceGroup
	DefaultColumns []string
}

func (o *Orchestrator) List(ctx context.Context, locStr string, queryOverride string, opts RunOptions, sortBy string) (*ListResult, error) {
	sel, prov, err := resolver.ResolveSelection(locStr, queryOverride, o.buildResolveOptions(opts))
	if err != nil {
		return nil, err
	}

	if sortBy != "" {
		if sel.Location.Resource == "tracks" {
			if _, ok := models.TrackFields[sortBy]; !ok {
				return nil, fmt.Errorf("invalid sort field %q; valid fields are: %v", sortBy, strings.Join(query.AllowedTrackFields, ", "))
			}
			prov.Tracks().Sort(ctx, o.buildExecContext(opts), sel.Tracks, sortBy)
		} else {
			if _, ok := models.GroupFields[sortBy]; !ok {
				return nil, fmt.Errorf("invalid sort field %q; valid fields are: %v", sortBy, strings.Join(query.AllowedGroupFields, ", "))
			}
			prov.Groups().Sort(ctx, o.buildExecContext(opts), sel.Groups, sortBy)
		}
	}

	return &ListResult{
		Resource:       sel.Location.Resource,
		Tracks:         sel.Tracks,
		Groups:         sel.Groups,
		DefaultColumns: prov.System().TableHeaders(),
	}, nil
}

func (o *Orchestrator) Stats(ctx context.Context, locStr, queryOverride string, opts RunOptions) (*StatResult, error) {
	sel, _, err := resolver.ResolveSelection(locStr, queryOverride, o.buildResolveOptions(opts))
	if err != nil {
		return nil, err
	}

	if sel.Location.Resource != "tracks" {
		return nil, fmt.Errorf("stats only available for track resources")
	}

	res := &StatResult{
		Count:   len(sel.Tracks),
		Genres:  make(map[string]int),
		Labels:  make(map[string]int),
		Keys:    make(map[string]int),
		Artists: make(map[string]int),
	}

	if len(sel.Tracks) == 0 {
		return res, nil
	}

	totalBPM := 0.0
	for _, t := range sel.Tracks {
		if t.Genre != "" {
			res.Genres[t.Genre]++
		}
		if t.Label != "" {
			res.Labels[t.Label]++
		}
		if t.Key != "" {
			res.Keys[t.Key]++
		}
		if t.Artist != "" {
			res.Artists[t.Artist]++
		}
		totalBPM += t.BPM
	}
	res.AvgBPM = totalBPM / float64(len(sel.Tracks))

	return res, nil
}

type SyncOptions struct {
	ExportDest     string
	ExportFormat   string
	PathMaps       map[string]string
	AppendOnly     bool
	MetadataFields []string
	MatchFields    []string
}

type FixKind string

const (
	FixDuplicates FixKind = "duplicates"
	FixMetadata   FixKind = "metadata"
	FixPaths      FixKind = "paths"
	FixOrphans    FixKind = "orphans"
)

type FixOptions struct {
	Actions map[FixKind][]string
}

func (o *Orchestrator) Sync(ctx context.Context, sourceLoc, targetLoc string, queryOverride string, opts RunOptions, syncOpts SyncOptions) error {
	src, _, err := resolver.ResolveSelection(sourceLoc, queryOverride, o.buildResolveOptions(opts))
	if err != nil {
		return err
	}

	tgt, prov, err := resolver.ResolveSelection(targetLoc, "", o.buildResolveOptions(opts))
	if err != nil {
		return err
	}

	resolvedTargetID := tgt.Location.Query
	if len(tgt.Groups) > 0 {
		resolvedTargetID = tgt.Groups[0].ID
	}

	pSyncOpts := provider.SyncOptions{
		ExportDest:     syncOpts.ExportDest,
		ExportFormat:   syncOpts.ExportFormat,
		PathMaps:       syncOpts.PathMaps,
		AppendOnly:     syncOpts.AppendOnly,
		MetadataFields: syncOpts.MetadataFields,
		MatchFields:    syncOpts.MatchFields,
	}

	err = prov.System().Sync(ctx, o.buildExecContext(opts), src.Tracks, resolvedTargetID, pSyncOpts)
	if err != nil {
		return err
	}

	if opts.Apply {
		return prov.System().Save(ctx, o.buildExecContext(opts), opts.FilePath)
	}
	return nil
}

type SyncDiff struct {
	TargetName  string
	CurrentIDs  []string
	AddedIDs    []string
	RemovedIDs  []string
	SourceIDs   []string
	TrackLookup map[string]models.Track
}

func (o *Orchestrator) GetSyncDiff(ctx context.Context, sourceLoc, targetLoc string, queryOverride string, opts RunOptions, appendOnly bool) (*SyncDiff, error) {
	src, _, err := resolver.ResolveSelection(sourceLoc, queryOverride, o.buildResolveOptions(opts))
	if err != nil {
		return nil, err
	}

	tgt, prov, err := resolver.ResolveSelection(targetLoc, "", o.buildResolveOptions(opts))
	if err != nil {
		return nil, err
	}

	diff := &SyncDiff{
		TargetName:  tgt.Location.Query,
		TrackLookup: make(map[string]models.Track),
	}
	if diff.TargetName == "" {
		diff.TargetName = tgt.Location.Resource
	}

	if len(tgt.Groups) > 0 {
		group := tgt.Groups[0]
		diff.TargetName = group.Name

		var currentIDs []string
		allTracks, _ := prov.Tracks().List(ctx, o.buildExecContext(opts), "")
		for _, t := range allTracks {
			diff.TrackLookup[t.ID] = t
			for _, m := range t.Playlists {
				if m.Name == group.Name && m.Folder == group.ParentFolder {
					currentIDs = append(currentIDs, t.ID)
					break
				}
			}
		}

		var sourceIDs []string
		for _, t := range src.Tracks {
			sourceIDs = append(sourceIDs, t.ID)
			diff.TrackLookup[t.ID] = t
		}

		added, removed := calculateSyncDiff(currentIDs, sourceIDs)
		diff.CurrentIDs = currentIDs
		diff.AddedIDs = added
		diff.RemovedIDs = removed
		diff.SourceIDs = sourceIDs
	}

	return diff, nil
}

func (o *Orchestrator) Fix(ctx context.Context, locStr string, queryOverride string, opts RunOptions, fixOpts FixOptions) (int, error) {
	sel, prov, err := resolver.ResolveSelection(locStr, queryOverride, o.buildResolveOptions(opts))
	if err != nil {
		return 0, err
	}

	if len(sel.Items) == 0 {
		return 0, nil
	}

	pFixOpts := provider.FixOptions{
		Actions: make(map[provider.FixType][]string),
	}
	for k, v := range fixOpts.Actions {
		pFixOpts.Actions[provider.FixType(k)] = v
	}

	count, err := prov.System().Fix(ctx, o.buildExecContext(opts), *sel, pFixOpts)
	if err != nil {
		return count, err
	}

	if opts.Apply {
		return count, prov.System().Save(ctx, o.buildExecContext(opts), opts.FilePath)
	}
	return count, nil
}

func (o *Orchestrator) Edit(ctx context.Context, locStr string, queryOverride string, opts RunOptions, changes map[string]string) (int, error) {
	sel, prov, err := resolver.ResolveSelection(locStr, queryOverride, o.buildResolveOptions(opts))
	if err != nil {
		return 0, err
	}

	count, err := prov.Tracks().Update(ctx, o.buildExecContext(opts), sel.Location.Query, changes)
	if err != nil {
		return count, err
	}

	if opts.Apply {
		return count, prov.System().Save(ctx, o.buildExecContext(opts), opts.FilePath)
	}
	return count, nil
}

func (o *Orchestrator) Make(ctx context.Context, locStr string, name string, opts RunOptions, groupKind models.GroupKind, position int, fromLoc string) (models.ResourceGroup, error) {
	_, prov, err := resolver.ResolveSelection(locStr, "", o.buildResolveOptions(opts))
	if err != nil {
		return models.ResourceGroup{}, err
	}

	newNode, err := prov.Groups().Create(ctx, o.buildExecContext(opts), models.ResourceGroup{}, name, groupKind, position)
	if err != nil {
		return models.ResourceGroup{}, err
	}

	if fromLoc != "" {
		src, _, err := resolver.ResolveSelection(fromLoc, "", o.buildResolveOptions(opts))
		if err != nil {
			return newNode, err
		}
		if len(src.Tracks) > 0 {
			_, err = prov.Tracks().Groups().Add(ctx, o.buildExecContext(opts), src.Tracks, newNode)
			if err != nil {
				return newNode, err
			}
		}
	}

	if opts.Apply {
		return newNode, prov.System().Save(ctx, o.buildExecContext(opts), opts.FilePath)
	}
	return newNode, nil
}

func (o *Orchestrator) Move(ctx context.Context, locStr string, queryOverride string, opts RunOptions, moveTo string, moveFrom string, moveName string) error {
	src, prov, err := resolver.ResolveSelection(locStr, queryOverride, o.buildResolveOptions(opts))
	if err != nil {
		return err
	}

	if moveName != "" {
		if err := prov.Groups().Update(ctx, o.buildExecContext(opts), src.Groups[0], moveName, nil); err != nil {
			return err
		}
	} else if src.Location.Resource == "tracks" {
		org, _, err := resolver.ResolveSelection(moveFrom, "", o.buildResolveOptions(opts))
		if err != nil {
			return err
		}
		tgt, _, err := resolver.ResolveSelection(moveTo, "", o.buildResolveOptions(opts))
		if err != nil {
			return err
		}
		for _, origin := range org.Groups {
			for _, target := range tgt.Groups {
				if _, err := prov.Tracks().Groups().Move(ctx, o.buildExecContext(opts), src.Tracks, origin, target); err != nil {
					return err
				}
			}
		}
	} else {
		tgt, _, err := resolver.ResolveSelection(moveTo, "", o.buildResolveOptions(opts))
		if err != nil {
			return err
		}
		targetParent := tgt.Groups[0]
		for _, t := range src.Groups {
			if err := prov.Groups().Update(ctx, o.buildExecContext(opts), t, "", &targetParent); err != nil {
				continue
			}
		}
	}

	if opts.Apply {
		return prov.System().Save(ctx, o.buildExecContext(opts), opts.FilePath)
	}
	return nil
}

func (o *Orchestrator) Delete(ctx context.Context, locStr string, queryOverride string, opts RunOptions, fromLocs []string, recursive bool) error {
	sel, prov, err := resolver.ResolveSelection(locStr, queryOverride, o.buildResolveOptions(opts))
	if err != nil {
		return err
	}

	if len(fromLocs) == 0 {
		for _, item := range sel.Items {
			if node, ok := item.(models.ResourceGroup); ok {
				if recursive && node.Kind == models.GroupKindFolder {
					children, _ := prov.Groups().List(ctx, o.buildExecContext(opts), "parent:"+node.Name)
					for _, c := range children {
						prov.Groups().Delete(ctx, o.buildExecContext(opts), c)
					}
				}
				if err := prov.Groups().Delete(ctx, o.buildExecContext(opts), node); err != nil {
					return err
				}
			}
		}
	} else {
		for _, fromStr := range fromLocs {
			org, _, err := resolver.ResolveSelection(fromStr, "", o.buildResolveOptions(opts))
			if err != nil {
				return err
			}
			for _, target := range org.Groups {
				if _, err := prov.Tracks().Groups().Remove(ctx, o.buildExecContext(opts), sel.Tracks, target); err != nil {
					return err
				}
			}
		}
	}

	if opts.Apply {
		return prov.System().Save(ctx, o.buildExecContext(opts), opts.FilePath)
	}
	return nil
}

func calculateSyncDiff(current, target []string) (added, removed []string) {
	currentMap := make(map[string]bool)
	for _, id := range current {
		currentMap[id] = true
	}
	targetMap := make(map[string]bool)
	for _, id := range target {
		targetMap[id] = true
	}

	for _, id := range target {
		if !currentMap[id] {
			added = append(added, id)
		}
	}
	for _, id := range current {
		if !targetMap[id] {
			removed = append(removed, id)
		}
	}
	return
}
