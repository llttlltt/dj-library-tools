package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/utils"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

// BulkAction is a function that performs an action on a target and returns the number of items affected.
type BulkAction func(targetName string, items []string) (bool, int)

// RunBulkOperation encapsulates the 'Loop Targets -> Apply with Progress Bar & Verbose logging' pattern.
func RunBulkOperation(verb string, targetNames []string, itemIDs []string, action BulkAction) {
	preposition := "to"
	if verb == "remove" {
		preposition = "from"
	}
	if dryRun {
		for _, name := range targetNames {
			fmt.Printf("[Dry Run] Would %s %d tracks %s playlist %q\n", verb, len(itemIDs), preposition, name)
		}
		return
	}

	p := mpb.New(mpb.WithWidth(64))
	for _, name := range targetNames {
		bar := p.AddBar(int64(len(itemIDs)),
			mpb.PrependDecorators(
				decor.Name(fmt.Sprintf("%s in %q", stringsTitle(verb), name), decor.WCSyncSpaceR),
				decor.CountersNoUnit("%d / %d", decor.WCSyncSpace),
			),
			mpb.AppendDecorators(decor.Percentage(decor.WCSyncSpace)),
		)

		if verbose {
			fmt.Printf("%s %d items in %q...\n", stringsTitle(verb), len(itemIDs), name)
		}

		// We process in chunks to show progress
		chunkSize := 10
		if len(itemIDs) < chunkSize {
			chunkSize = len(itemIDs)
		}

		totalAffected := 0
		for i := 0; i < len(itemIDs); i += chunkSize {
			end := i + chunkSize
			if end > len(itemIDs) {
				end = len(itemIDs)
			}
			chunk := itemIDs[i:end]
			found, affected := action(name, chunk)
			if !found {
				fmt.Printf("Warning: target %q not found during %s\n", name, verb)
				bar.Abort(false)
				break
			}
			totalAffected += affected
			bar.IncrBy(len(chunk))
		}
		p.Wait()
		fmt.Printf("%s %d items in %q\n", stringsTitle(verb), totalAffected, name)
	}
}

// Selection represents a resolved set of resources from a provider.
type Selection struct {
	Items    []models.Resource
	Tracks   []models.Track // Convenience helpers
	Nodes    []models.Node  // Convenience helpers
	Location utils.Location
	Provider provider.Provider
}

// ResolveSelection resolves a location string into a Selection.
func ResolveSelection(locStr string, queryOverride string) (*Selection, error) {
	if locStr == "" {
		return &Selection{}, nil
	}
	loc := utils.ParseLocation(locStr, queryOverride)
	if loc.Resource == "" {
		return nil, fmt.Errorf("resource must be specified in location %q (e.g. %s/tracks)", locStr, loc.Provider)
	}

	cfg, _ := config.LoadAppConfig()
	var rbXML, xmlPath, _ = loadXMLFunc()

	prov, err := provider.NewProvider(loc.Provider, rbXML, xmlPath, cfg)
	if err != nil {
		return nil, err
	}

	// For M3U, the "Resource" part of the location is actually the file path.
	filePath := loc.Resource
	isM3U := loc.Provider == "m3u" || loc.Provider == "m3u8"
	if isM3U {
		// If the resource was something like "test.m3u8/tracks", we should strip /tracks
		if strings.HasSuffix(filePath, "/tracks") {
			filePath = strings.TrimSuffix(filePath, "/tracks")
			loc.Resource = "tracks"
		} else {
			loc.Resource = "playlists"
		}

		m3uProv, err := provider.NewM3UProvider(filePath)
		if err != nil {
			return nil, err
		}
		prov = m3uProv
	}

	sel := &Selection{Location: loc, Provider: prov}
	if isM3U {
		if loc.Resource == "playlists" {
			nodes, _ := prov.GetPlaylists("")
			sel.Nodes = nodes
			for _, n := range nodes {
				sel.Items = append(sel.Items, n)
			}
		}
	}

	if loc.Resource == "tracks" || isM3U {
		tracks, err := prov.GetTracks(loc.Query)
		if err != nil {
			return nil, err
		}
		sel.Tracks = tracks
		// Only add to sel.Items if we haven't already added nodes
		if len(sel.Items) == 0 {
			for _, t := range tracks {
				sel.Items = append(sel.Items, t)
			}
		}
	} else {
		var nodes []models.Node
		var err error
		if loc.Resource == "folders" {
			nodes, err = prov.GetFolders(loc.Query)
		} else {
			nodes, err = prov.GetPlaylists(loc.Query)
		}
		if err != nil {
			return nil, err
		}
		sel.Nodes = nodes
		for _, n := range nodes {
			sel.Items = append(sel.Items, n)
		}
	}

	return sel, nil
}

// stringsTitle is a simple helper since strings.Title is deprecated.
func stringsTitle(s string) string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("%c%s", s[0]-32, s[1:])
}
