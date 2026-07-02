package orchestrator

import (
	"context"
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/llttlltt/dj-library-tools/internal/core/query"
	"github.com/llttlltt/dj-library-tools/internal/providers"
	"github.com/llttlltt/dj-library-tools/internal/providers/factory"
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

// NewWithFeedback returns a shallow copy of base with Feedback replaced by fb.
// Used by the workflow engine to capture per-Step output without recreating
// the full orchestrator configuration.
func NewWithFeedback(fb provider.Feedback, base *Orchestrator) *Orchestrator {
	if fb == nil {
		fb = provider.NoopFeedback{}
	}
	return &Orchestrator{
		Feedback:             fb,
		RekordboxPrimaryPath: base.RekordboxPrimaryPath,
	}
}

type RunOptions struct {
	FilePath string
	Apply    bool
	Verbose  bool
	Host     string
	Port     int
	Token    string
}

func (o *Orchestrator) buildResolveOptions(opts RunOptions) resolver.ResolveOptions {
	return resolver.ResolveOptions{
		FilePath:             opts.FilePath,
		RekordboxPrimaryPath: o.RekordboxPrimaryPath,
		Apply:                opts.Apply,
		Verbose:              opts.Verbose,
		Host:                 opts.Host,
		Port:                 opts.Port,
		Token:                opts.Token,
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
	if err := o.validateLocation(locStr); err != nil {
		return nil, err
	}
	sel, prov, err := resolver.ResolveSelection(ctx, locStr, queryOverride, o.buildResolveOptions(opts))
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
	sel, _, err := resolver.ResolveSelection(ctx, locStr, queryOverride, o.buildResolveOptions(opts))
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

type FixOptions struct {
	Actions map[provider.FixType][]string
}

func (o *Orchestrator) Sync(ctx context.Context, sourceLoc, targetLoc string, queryOverride string, opts RunOptions, syncOpts SyncOptions) error {
	if err := o.validateLocation(sourceLoc); err != nil {
		return err
	}
	if err := o.validateTargetLocation(targetLoc); err != nil {
		return err
	}
	src, _, err := resolver.ResolveSelection(ctx, sourceLoc, queryOverride, o.buildResolveOptions(opts))
	if err != nil {
		return err
	}

	tgt, prov, err := resolver.ResolveSelection(ctx, targetLoc, "", o.buildResolveOptions(opts))
	if err != nil {
		return err
	}

	resolvedTargetIDs := []string{tgt.Location.Query}
	if len(tgt.Groups) > 0 {
		resolvedTargetIDs = []string{}
		for _, g := range tgt.Groups {
			resolvedTargetIDs = append(resolvedTargetIDs, g.ID)
		}
	}

	pSyncOpts := provider.SyncOptions{
		ExportDest:     syncOpts.ExportDest,
		ExportFormat:   syncOpts.ExportFormat,
		PathMaps:       syncOpts.PathMaps,
		AppendOnly:     syncOpts.AppendOnly,
		MetadataFields: syncOpts.MetadataFields,
		MatchFields:    syncOpts.MatchFields,
	}

	for _, targetID := range resolvedTargetIDs {
		err = prov.System().Sync(ctx, o.buildExecContext(opts), src.Tracks, targetID, pSyncOpts)
		if err != nil {
			return err
		}
	}

	if opts.Apply {
		if !prov.System().Capabilities().CanWrite {
			return fmt.Errorf("provider %q does not support writing", prov.Name())
		}
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

func (o *Orchestrator) GetSyncDiff(ctx context.Context, sourceLoc, targetLoc string, queryOverride string, opts RunOptions, appendOnly bool) ([]*SyncDiff, error) {
	if err := o.validateLocation(sourceLoc); err != nil {
		return nil, err
	}
	if err := o.validateTargetLocation(targetLoc); err != nil {
		return nil, err
	}
	src, _, err := resolver.ResolveSelection(ctx, sourceLoc, queryOverride, o.buildResolveOptions(opts))
	if err != nil {
		return nil, err
	}

	tgt, prov, err := resolver.ResolveSelection(ctx, targetLoc, "", o.buildResolveOptions(opts))
	if err != nil {
		return nil, err
	}

	var diffs []*SyncDiff

	// If no groups matched but we have a target location,
	// we still want one diff for the potential new group.
	targetGroups := tgt.Groups
	if len(targetGroups) == 0 {
		name := tgt.Location.Query
		if name == "" {
			name = tgt.Location.Resource
		}
		diffs = append(diffs, &SyncDiff{
			TargetName:  name,
			TrackLookup: make(map[string]models.Track),
		})
	}

	for _, group := range targetGroups {
		diff := &SyncDiff{
			TargetName:  group.Name,
			TrackLookup: make(map[string]models.Track),
		}

		allTracks, _ := prov.Tracks().List(ctx, o.buildExecContext(opts), "")
		for _, t := range allTracks {
			diff.TrackLookup[t.ID] = t
		}

		var currentIDs []string
		for _, t := range allTracks {
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
		diffs = append(diffs, diff)
	}

	return diffs, nil
}

func (o *Orchestrator) Fix(ctx context.Context, locStr string, queryOverride string, opts RunOptions, fixOpts FixOptions) (int, error) {
	if err := o.validateLocation(locStr); err != nil {
		return 0, err
	}
	sel, prov, err := resolver.ResolveSelection(ctx, locStr, queryOverride, o.buildResolveOptions(opts))
	if err != nil {
		return 0, err
	}

	if len(sel.Items) == 0 {
		return 0, nil
	}

	pFixOpts := provider.FixOptions{
		Actions: fixOpts.Actions,
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
	if err := o.validateTargetLocation(locStr); err != nil {
		return 0, err
	}
	sel, prov, err := resolver.ResolveSelection(ctx, locStr, queryOverride, o.buildResolveOptions(opts))
	if err != nil {
		return 0, err
	}

	if opts.Apply && !prov.System().Capabilities().CanUpdateMetadata {
		return 0, fmt.Errorf("provider %q does not support metadata updates", prov.Name())
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
	if err := o.validateTargetLocation(locStr); err != nil {
		return models.ResourceGroup{}, err
	}
	sel, prov, err := resolver.ResolveSelection(ctx, locStr, "", o.buildResolveOptions(opts))
	if err != nil {
		return models.ResourceGroup{}, err
	}

	var parent models.ResourceGroup
	if sel.Location.Query != "" && len(sel.Groups) > 0 {
		parent = sel.Groups[0]
	}

	newNode, err := prov.Groups().Create(ctx, o.buildExecContext(opts), parent, name, groupKind, position)
	if err != nil {
		return models.ResourceGroup{}, err
	}

	if fromLoc != "" {
		src, _, err := resolver.ResolveSelection(ctx, fromLoc, "", o.buildResolveOptions(opts))
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

func (o *Orchestrator) Move(ctx context.Context, locStr string, queryOverride string, opts RunOptions, moveTo string, moveFrom string, moveName string) (int, error) {
	if err := o.validateTargetLocation(locStr); err != nil {
		return 0, err
	}
	src, prov, err := resolver.ResolveSelection(ctx, locStr, queryOverride, o.buildResolveOptions(opts))
	if err != nil {
		return 0, err
	}

	count := 0
	if moveName != "" {
		if err := prov.Groups().Update(ctx, o.buildExecContext(opts), src.Groups[0], moveName, nil); err != nil {
			return 0, err
		}
		count = 1
	} else if src.Location.Resource == "tracks" {
		org, _, err := resolver.ResolveSelection(ctx, moveFrom, "", o.buildResolveOptions(opts))
		if err != nil {
			return 0, err
		}
		tgt, _, err := resolver.ResolveSelection(ctx, moveTo, "", o.buildResolveOptions(opts))
		if err != nil {
			return 0, err
		}
		for _, origin := range org.Groups {
			for _, target := range tgt.Groups {
				if _, err := prov.Tracks().Groups().Move(ctx, o.buildExecContext(opts), src.Tracks, origin, target); err != nil {
					return 0, err
				}
			}
		}
		count = len(src.Tracks)
	} else {
		tgt, _, err := resolver.ResolveSelection(ctx, moveTo, "", o.buildResolveOptions(opts))
		if err != nil {
			return 0, err
		}
		targetParent := tgt.Groups[0]
		for _, t := range src.Groups {
			if err := prov.Groups().Update(ctx, o.buildExecContext(opts), t, "", &targetParent); err != nil {
				continue
			}
		}
		count = len(src.Groups)
	}

	if opts.Apply {
		return count, prov.System().Save(ctx, o.buildExecContext(opts), opts.FilePath)
	}
	return count, nil
}

func (o *Orchestrator) Delete(ctx context.Context, locStr string, queryOverride string, opts RunOptions, fromLocs []string, recursive bool) (int, error) {
	if err := o.validateTargetLocation(locStr); err != nil {
		return 0, err
	}
	sel, prov, err := resolver.ResolveSelection(ctx, locStr, queryOverride, o.buildResolveOptions(opts))
	if err != nil {
		return 0, err
	}

	count := len(sel.Items)

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
					return 0, err
				}
			}
		}
	} else {
		for _, fromStr := range fromLocs {
			org, _, err := resolver.ResolveSelection(ctx, fromStr, "", o.buildResolveOptions(opts))
			if err != nil {
				return 0, err
			}
			for _, target := range org.Groups {
				if _, err := prov.Tracks().Groups().Remove(ctx, o.buildExecContext(opts), sel.Tracks, target); err != nil {
					return 0, err
				}
			}
		}
	}

	if opts.Apply {
		return count, prov.System().Save(ctx, o.buildExecContext(opts), opts.FilePath)
	}
	return count, nil
}

func (o *Orchestrator) validateLocation(locStr string) error {
	return o.doValidateLocation(locStr, false)
}

func (o *Orchestrator) validateTargetLocation(locStr string) error {
	return o.doValidateLocation(locStr, true)
}

func (o *Orchestrator) doValidateLocation(locStr string, mustBeWritable bool) error {
	if locStr == "" {
		return nil
	}
	if strings.HasPrefix(locStr, "m3u://") || strings.HasPrefix(locStr, "m3u8://") {
		return nil
	}
	parts := strings.SplitN(locStr, "/", 2)
	if len(parts) < 2 {
		return nil // location string might just be a query or malformed, let resolver handle it
	}
	prov := parts[0]
	resPart := parts[1]
	res := strings.SplitN(resPart, " ", 2)[0]

	// If prov is a UUID, we can't easily validate without config lookup.
	// For now we only validate short names like "rb", "plex".
	if len(prov) == 36 && strings.Contains(prov, "-") {
		return nil
	}

	if !factory.ValidateResource(prov, res, mustBeWritable) {
		// Only error if the provider is actually known. If it's unknown,
		// let the factory.NewProvider return the error later.
		if info, err := factory.GetProviderInfo(prov); err == nil {
			if mustBeWritable {
				return fmt.Errorf("provider %q does not support writing to resource %q", prov, res)
			}
			return fmt.Errorf("provider %q does not support resource %q; supported: %v", prov, res, resourceNames(info.Resources))
		}
	}
	return nil
}

func resourceNames(rs []factory.ResourceInfo) []string {
	var names []string
	for _, r := range rs {
		names = append(names, r.Name)
	}
	return names
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
