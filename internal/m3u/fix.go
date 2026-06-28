package m3u

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FixOptions holds the configuration for the playlist fix operation.
type FixOptions struct {
	Exts           []string
	M3U8           bool
	RemoveOriginal bool
	Force          bool
	OutputPath     string
	Verbose        bool
	DryRun         bool
}

// FixResult holds the outcome of the fix operation, including any missing files found.
type FixResult struct {
	TotalTracks   int
	SkippedTracks []string
	OutputPath    string
}

// FixPlaylist performs the playlist hygiene operations and returns a summary.
func FixPlaylist(inputPath string, opts FixOptions) (*FixResult, error) {
	result := &FixResult{}
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("input file '%s' not found", inputPath)
	}

	outputPath := inputPath
	if opts.OutputPath != "" {
		outputPath = opts.OutputPath
	} else if opts.M3U8 {
		ext := filepath.Ext(inputPath)
		if ext == ".m3u" {
			outputPath = strings.TrimSuffix(inputPath, ".m3u") + ".m3u8"
		} else if ext != ".m3u8" {
			return nil, fmt.Errorf("input file must be .m3u or .m3u8 for M3U8 enrichment")
		}
	}

	if !opts.DryRun && outputPath != inputPath {
		if _, err := os.Stat(outputPath); err == nil && !opts.Force {
			return nil, fmt.Errorf("output file '%s' already exists. Use --force to overwrite", outputPath)
		}
	}

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	var tmpFile *os.File
	if !opts.DryRun {
		outputDir := filepath.Dir(outputPath)
		tmpFile, err = os.CreateTemp(outputDir, "djlt-playlist-*")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp file in %s: %w", outputDir, err)
		}
		defer os.Remove(tmpFile.Name())
	}

	fileDir := filepath.Dir(inputPath)

	if opts.Verbose {
		if opts.DryRun {
			fmt.Printf("Dry run: Analyzing playlist: %s\n", inputPath)
		} else {
			fmt.Printf("Starting playlist fix for: %s\n", inputPath)
		}
	}

	if opts.M3U8 && !opts.DryRun {
		if err := WriteM3U8Header(tmpFile); err != nil {
			return nil, fmt.Errorf("failed to write M3U8 header: %w", err)
		}
	}

	scanner := bufio.NewScanner(inputFile)
	trackCount := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			if opts.M3U8 && strings.HasPrefix(line, "#EXTM3U") {
				continue
			}
			if !opts.DryRun {
				if _, err := fmt.Fprintln(tmpFile, line); err != nil {
					return nil, fmt.Errorf("failed to write line: %w", err)
				}
			}
			continue
		}

		trackCount++
		result.TotalTracks++

		foundPath := ""
		resolvedPath := ""

		// Resolution logic
		if len(opts.Exts) == 0 {
			absPath := line
			if !filepath.IsAbs(line) {
				absPath = filepath.Join(fileDir, line)
			}
			if _, err := os.Stat(absPath); err == nil {
				foundPath = absPath
				resolvedPath = line
			}
		} else {
			for _, ext := range opts.Exts {
				testPath := FormatPath(line, ext)
				absTestPath := testPath
				if !filepath.IsAbs(testPath) {
					absTestPath = filepath.Join(fileDir, testPath)
				}
				if _, err := os.Stat(absTestPath); err == nil {
					foundPath = absTestPath
					resolvedPath = testPath
					break
				}
			}
		}

		// Fallback to original
		if foundPath == "" {
			absPath := line
			if !filepath.IsAbs(line) {
				absPath = filepath.Join(fileDir, line)
			}
			if _, err := os.Stat(absPath); err == nil {
				foundPath = absPath
				resolvedPath = line
			}
		}

		// Skip if not found
		if foundPath == "" {
			result.SkippedTracks = append(result.SkippedTracks, line)
			if opts.Verbose {
				fmt.Printf("[%d] ❌ Skipping (Not found): %s\n", trackCount, line)
			}
			continue
		}

		if opts.Verbose {
			fmt.Printf("[%d] ✔ Resolved: %s\n", trackCount, resolvedPath)
		} else if trackCount%50 == 0 {
			fmt.Printf("Processing tracks... (%d done)\n", trackCount)
		}

		if opts.M3U8 {
			displayName := filepath.Base(resolvedPath)
			if !opts.DryRun {
				if err := WriteM3U8EntryRaw(tmpFile, displayName, resolvedPath, -1); err != nil {
					return nil, fmt.Errorf("failed to write M3U8 entry: %w", err)
				}
			}
		} else {
			if !opts.DryRun {
				if _, err := fmt.Fprintln(tmpFile, resolvedPath); err != nil {
					return nil, err
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	if opts.DryRun {
		result.OutputPath = outputPath
		return result, nil
	}

	tmpFile.Close()

	if err := os.Rename(tmpFile.Name(), outputPath); err != nil {
		return nil, fmt.Errorf("failed to save output file: %w", err)
	}

	result.OutputPath = outputPath
	return result, nil
}
