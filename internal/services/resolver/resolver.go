package resolver

import (
	"context"
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/core/location"
	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/llttlltt/dj-library-tools/internal/core/query"
	provider "github.com/llttlltt/dj-library-tools/internal/providers"
	"github.com/llttlltt/dj-library-tools/internal/providers/factory"
)

// Selection represents a resolved set of resources from a provider.
type Selection = provider.Selection

// ResolveOptions holds configuration for selection resolution.
type ResolveOptions struct {
	FilePath             string
	RekordboxPrimaryPath string
	Apply                bool
	Verbose              bool
	Host                 string
	Port                 int
	Token                string
	Feedback             provider.Feedback
}

// ResolveSelection resolves a location string into a Selection.
func ResolveSelection(ctx context.Context, locStr string, queryOverride string, opts ResolveOptions) (*Selection, provider.Provider, error) {
	if locStr == "" {
		return &Selection{}, nil, nil
	}

	filePath := opts.FilePath
	host := opts.Host
	port := opts.Port
	token := opts.Token

	var loc location.Location
	if strings.HasPrefix(locStr, "m3u://") || strings.HasPrefix(locStr, "m3u8://") {
		scheme := "m3u"
		if strings.HasPrefix(locStr, "m3u8://") {
			scheme = "m3u8"
		}

		// Split query if present in locStr and not overridden
		parts := strings.SplitN(locStr, " ", 2)
		uri := parts[0]
		if len(parts) > 1 && queryOverride == "" {
			queryOverride = parts[1]
		}

		prefix := scheme + "://"
		pathAndResource := strings.TrimPrefix(uri, prefix)

		loc.Provider = scheme
		loc.Query = queryOverride
		loc.Resource = "tracks"
		filePath = pathAndResource

		if strings.HasSuffix(pathAndResource, "/tracks") {
			loc.Resource = "tracks"
			filePath = strings.TrimSuffix(pathAndResource, "/tracks")
		} else if strings.HasSuffix(pathAndResource, "/playlists") {
			loc.Resource = "playlists"
			filePath = strings.TrimSuffix(pathAndResource, "/playlists")
		}
	} else {
		loc = location.ParseLocation(locStr, queryOverride)
		if loc.Resource == "" {
			return nil, nil, fmt.Errorf("resource must be specified in location %q", locStr)
		}

		// If provider is a UUID, load the Connection and use its config
		if len(loc.Provider) == 36 && strings.Contains(loc.Provider, "-") {
			conn, err := config.FindConnectionByID(loc.Provider)
			if err == nil {
				loc.Provider = conn.Provider // replace UUID with "rb", "plex" etc.

				// Override options with Connection config
				if fp, ok := conn.Config["file_path"]; ok {
					filePath = fp
				}
				if h, ok := conn.Config["host"]; ok {
					host = h
				}
				if p, ok := conn.Config["port"]; ok {
					fmt.Sscanf(p, "%d", &port)
				}
				if t, ok := conn.Config["token"]; ok {
					token = t
				}
			}
		}
	}

	if filePath == "" && loc.Provider == "rb" {
		filePath = opts.RekordboxPrimaryPath
	}

	factoryOpts := factory.ProviderOptions{
		FilePath: filePath,
		Host:     host,
		Port:     port,
		Token:    token,
	}

	prov, err := factory.NewProvider(loc.Provider, factoryOpts)
	if err != nil {
		return nil, nil, err
	}

	sel := &Selection{Location: loc}
	execCtx := provider.ExecutionContext{
		Apply:    opts.Apply,
		Verbose:  opts.Verbose,
		Feedback: opts.Feedback,
	}
	if execCtx.Feedback == nil {
		execCtx.Feedback = provider.NoopFeedback{}
	}

	var items []models.Resource
	if loc.Resource == "tracks" {
		tracks, err := prov.Tracks().List(ctx, execCtx, loc.Query)
		if err != nil {
			return nil, nil, err
		}
		for _, t := range tracks {
			items = append(items, t)
		}
	} else if loc.Resource == "playlists" || loc.Resource == "folders" {
		kindFilter := "kind:=playlist"
		if loc.Resource == "folders" {
			kindFilter = "kind:=folder"
		}
		effectiveQuery := loc.Query
		if effectiveQuery != "" {
			effectiveQuery = effectiveQuery + " " + kindFilter
		} else {
			effectiveQuery = kindFilter
		}
		groups, err := prov.Groups().List(ctx, execCtx, effectiveQuery)
		if err != nil {
			return nil, nil, err
		}
		for _, g := range groups {
			items = append(items, g)
		}
	} else {
		return nil, nil, fmt.Errorf("unsupported resource type: %s", loc.Resource)
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
