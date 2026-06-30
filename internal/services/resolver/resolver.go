package resolver

import (
	"fmt"
	"os"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/core/location"
	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/llttlltt/dj-library-tools/internal/providers"
	"github.com/llttlltt/dj-library-tools/internal/providers/factory"
	"github.com/llttlltt/dj-library-tools/internal/core/query"
)

// Selection represents a resolved set of resources from a provider.
type Selection = provider.Selection

// ResolveOptions holds configuration for selection resolution.
type ResolveOptions struct {
	FilePath             string
	RekordboxPrimaryPath string
	FilterMissing        bool
	FilterExists         bool
	Apply                bool
	Verbose              bool
	Feedback             provider.Feedback
}

// ResolveSelection resolves a location string into a Selection.
func ResolveSelection(locStr string, queryOverride string, opts ResolveOptions) (*Selection, provider.Provider, error) {
	if locStr == "" {
		return &Selection{}, nil, nil
	}
	loc := location.ParseLocation(locStr, queryOverride)
	if loc.Resource == "" {
		return nil, nil, fmt.Errorf("resource must be specified in location %q", locStr)
	}

	filePath := opts.FilePath
	if filePath == "" && loc.Provider == "rb" {
		filePath = opts.RekordboxPrimaryPath
	}

	factoryOpts := factory.ProviderOptions{
		FilePath: filePath,
	}

	prov, err := factory.NewProvider(loc.Provider, factoryOpts)
	if err != nil {
		return nil, nil, err
	}

	sel := &Selection{Location: loc}
	ctx := provider.ExecutionContext{
		Apply:    opts.Apply,
		Verbose:  opts.Verbose,
		Feedback: opts.Feedback,
	}
	if ctx.Feedback == nil {
		ctx.Feedback = provider.NoopFeedback{}
	}
	
	var items []models.Resource
	if loc.Resource == "tracks" {
		tracks, err := prov.Tracks().List(ctx, loc.Query)
		if err != nil { return nil, nil, err }
		for _, t := range tracks { items = append(items, t) }
	} else if loc.Resource == "playlists" || loc.Resource == "folders" {
		groups, err := prov.Groups().List(ctx, loc.Query)
		if err != nil { return nil, nil, err }
		for _, g := range groups { items = append(items, g) }
	} else {
		return nil, nil, fmt.Errorf("unsupported resource type: %s", loc.Resource)
	}

	// Apply physical health filtering if flags are set
	if opts.FilterMissing || opts.FilterExists {
		var filtered []models.Resource
		for _, item := range items {
			if t, ok := item.(models.Track); ok {
				_, statErr := os.Stat(t.Location)
				missing := os.IsNotExist(statErr)
				if opts.FilterMissing && !missing { continue }
				if opts.FilterExists && missing { continue }
			}
			filtered = append(filtered, item)
		}
		items = filtered
	}

	sel.Items = items
	for _, item := range items {
		if t, ok := item.(models.Track); ok {
			sel.Tracks = append(sel.Tracks, t)
		} else if g, ok := item.(models.ResourceGroup); ok {
			sel.Groups = append(sel.Groups, g)
		}
	}

	// Validate path-based query capabilities
	if loc.Query != "" {
		validatePathCapabilities(sel, prov, opts.Feedback)
	}

	return sel, prov, nil
}

func validatePathCapabilities(sel *Selection, prov provider.Provider, fb provider.Feedback) {
	if fb == nil {
		fb = provider.NoopFeedback{}
	}
	q := query.NewParser().Parse(sel.Location.Query)
	caps := prov.System().Capabilities()

	var check func(expr query.Expression)
	check = func(expr query.Expression) {
		switch v := expr.(type) {
		case query.Comparison:
			if strings.ContainsAny(v.Field, "./-") {
				collection := strings.Split(v.Field, ".")[0]
				collection = strings.Split(collection, "/")[0]
				collection = strings.Split(collection, "-")[0]

				if reqCap, ok := models.CollectionCapabilities[collection]; ok {
					hasCap := false
					switch reqCap {
					case models.CapCues:
						hasCap = caps.SupportsCues
					case models.CapBeatgrids:
						hasCap = caps.SupportsBeatgrids
					}
					if !hasCap {
						fb.OnWarning(fmt.Sprintf("The current provider (%s) does not support %q. Results for %q may be incomplete or empty.", 
							prov.Name(), collection, v.Field))
					}
				}
			}
		case query.Logical:
			check(v.Left)
			check(v.Right)
		case query.Not:
			check(v.Expr)
		}
	}
	if q.Root != nil {
		check(q.Root)
	}
}
