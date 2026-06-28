package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/djerr"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/resolver"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

// HandleError provides user-friendly messages for sentinel provider errors.
func HandleError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, djerr.ErrReadOnly) {
		return fmt.Errorf("operation failed: this provider is read-only")
	}
	if errors.Is(err, djerr.ErrUnsupportedResource) {
		return fmt.Errorf("operation failed: this resource type is not supported by the provider")
	}
	if errors.Is(err, djerr.ErrInvalidParent) {
		return fmt.Errorf("operation failed: cannot create the resource in that location (structural constraint)")
	}

	return err
}

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

func getExecContext() provider.ExecutionContext {
	return provider.ExecutionContext{
		DryRun:  dryRun,
		Verbose: verbose,
	}
}

func getResolveOptions() resolver.ResolveOptions {
	return resolver.ResolveOptions{
		FilePath:      filePath,
		FilterMissing: filterMissing,
		FilterExists:  filterExists,
		DryRun:        dryRun,
		Verbose:       verbose,
	}
}

func ResolveSelection(locStr string, queryOverride string) (*resolver.Selection, error) {
	return resolver.ResolveSelection(locStr, queryOverride, getResolveOptions())
}

func stringsTitle(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
