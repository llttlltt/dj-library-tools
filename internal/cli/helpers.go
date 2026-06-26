package cli

import (
	"fmt"

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

// stringsTitle is a simple helper since strings.Title is deprecated.
func stringsTitle(s string) string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("%c%s", s[0]-32, s[1:])
}
