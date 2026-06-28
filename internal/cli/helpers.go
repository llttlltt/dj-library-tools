package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/provider/factory"
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
	Tracks   []models.Track
	Nodes    []models.ResourceGroup
	Location utils.Location
	Provider provider.Provider
}

func getExecContext() provider.ExecutionContext {
	return provider.ExecutionContext{
		DryRun:  dryRun,
		Verbose: verbose,
	}
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
	
	opts := factory.ProviderOptions{
		FilePath: filePath,
		Config:   cfg,
	}

	if strings.Contains(loc.Provider, "rb") || strings.Contains(loc.Provider, "rekordbox") {
	}

	prov, err := factory.NewProvider(loc.Provider, opts)
	if err != nil {
		return nil, err
	}

	sel := &Selection{Location: loc, Provider: prov}
	ctx := getExecContext()
	
	items, err := prov.GetResources(ctx, loc.Resource, loc.Query)
	if err != nil {
		return nil, err
	}
	sel.Items = items

	for _, item := range items {
		if t, ok := item.(models.Track); ok {
			sel.Tracks = append(sel.Tracks, t)
		} else if g, ok := item.(models.ResourceGroup); ok {
			sel.Nodes = append(sel.Nodes, g)
		}
	}

	return sel, nil
}

func stringsTitle(s string) string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("%c%s", s[0]-32, s[1:])
}
