package resolver

import (
	"fmt"
	"os"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/provider/factory"
	"github.com/llttlltt/dj-library-tools/internal/query"
	"github.com/llttlltt/dj-library-tools/internal/utils"
)

// Selection represents a resolved set of resources from a provider.
type Selection struct {
	Items    []models.Resource
	Tracks   []models.Track
	Groups   []models.ResourceGroup
	Location utils.Location
	Provider provider.Provider
}

// ResolveOptions holds configuration for selection resolution.
type ResolveOptions struct {
	FilePath      string
	FilterMissing bool
	FilterExists  bool
	Apply         bool
	Verbose       bool
}

// ResolveSelection resolves a location string into a Selection.
func ResolveSelection(locStr string, queryOverride string, opts ResolveOptions) (*Selection, error) {
	if locStr == "" {
		return &Selection{}, nil
	}
	loc := utils.ParseLocation(locStr, queryOverride)
	if loc.Resource == "" {
		return nil, fmt.Errorf("resource must be specified in location %q", locStr)
	}

	cfg, _ := config.LoadAppConfig()

	filePath := opts.FilePath
	if filePath == "" && loc.Provider == "rb" {
		filePath = cfg.Rekordbox.PrimaryFilePath
	}

	factoryOpts := factory.ProviderOptions{
		FilePath: filePath,
		Config:   cfg,
	}

	prov, err := factory.NewProvider(loc.Provider, factoryOpts)
	if err != nil {
		return nil, err
	}

	sel := &Selection{Location: loc, Provider: prov}
	ctx := provider.ExecutionContext{
		Apply:   opts.Apply,
		Verbose: opts.Verbose,
	}
	
	var items []models.Resource
	if loc.Resource == "tracks" {
		tracks, err := prov.Tracks().List(ctx, loc.Query)
		if err != nil { return nil, err }
		for _, t := range tracks { items = append(items, t) }
	} else if loc.Resource == "playlists" || loc.Resource == "folders" {
		groups, err := prov.Groups().List(ctx, loc.Query)
		if err != nil { return nil, err }
		for _, g := range groups { items = append(items, g) }
	} else {
		return nil, fmt.Errorf("unsupported resource type: %s", loc.Resource)
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
		validatePathCapabilities(sel)
	}

	return sel, nil
}

func validatePathCapabilities(sel *Selection) {
	q := query.NewParser().Parse(sel.Location.Query)
	caps := sel.Provider.System().Capabilities()

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
						fmt.Fprintf(os.Stderr, "⚠️  Warning: The current provider (%s) does not support %q. Results for %q may be incomplete or empty.\n", 
							sel.Provider.Name(), collection, v.Field)
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
