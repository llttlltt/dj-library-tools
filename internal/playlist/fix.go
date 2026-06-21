package playlist

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FixOptions holds the configuration for the playlist fix operation.
type FixOptions struct {
	Ext            string
	M3U8           bool
	RemoveOriginal bool
	Force          bool
}

// FixPlaylist performs the playlist hygiene operations.
func FixPlaylist(inputPath string, opts FixOptions) error {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("input file '%s' not found", inputPath)
	}

	outputPath := inputPath
	if opts.M3U8 {
		ext := filepath.Ext(inputPath)
		if ext == ".m3u" {
			outputPath = strings.TrimSuffix(inputPath, ".m3u") + ".m3u8"
		} else if ext != ".m3u8" {
			return fmt.Errorf("input file must be .m3u or .m3u8 for M3U8 enrichment")
		}
	}

	if outputPath != inputPath {
		if _, err := os.Stat(outputPath); err == nil && !opts.Force {
			return fmt.Errorf("output file '%s' already exists. Use --force to overwrite", outputPath)
		}
	}

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	tmpFile, err := os.CreateTemp("", "djlt-playlist-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	fileDir := filepath.Dir(inputPath)

	if opts.M3U8 {
		if err := WriteM3U8Header(tmpFile); err != nil {
			return fmt.Errorf("failed to write M3U8 header: %w", err)
		}
	}

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			if opts.M3U8 && strings.HasPrefix(line, "#EXTM3U") {
				continue
			}
			if _, err := fmt.Fprintln(tmpFile, line); err != nil {
				return fmt.Errorf("failed to write line: %w", err)
			}
			continue
		}

		targetPath := line
		if opts.Ext != "" {
			targetPath = FormatPath(line, opts.Ext)
		}

		absTargetPath := targetPath
		if !filepath.IsAbs(targetPath) {
			absTargetPath = filepath.Join(fileDir, targetPath)
		}

		if opts.M3U8 {
			if _, err := os.Stat(absTargetPath); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Warning: Track not found: %s\n", absTargetPath)
				if _, err := fmt.Fprintln(tmpFile, targetPath); err != nil {
					return err
				}
			} else {
				meta, err := ExtractMetadata(absTargetPath)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: Could not read metadata for %s: %v\n", absTargetPath, err)
					if _, err := fmt.Fprintln(tmpFile, targetPath); err != nil {
						return err
					}
				} else {
					duration := 0.0 // Placeholder for now
					if err := WriteM3U8Entry(tmpFile, meta, targetPath, duration); err != nil {
						return fmt.Errorf("failed to write M3U8 entry: %w", err)
					}
				}
			}
		} else {
			if _, err := fmt.Fprintln(tmpFile, targetPath); err != nil {
				return fmt.Errorf("failed to write line: %w", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	tmpFile.Close()

	if err := os.Rename(tmpFile.Name(), outputPath); err != nil {
		return fmt.Errorf("failed to save output file: %w", err)
	}

	if opts.RemoveOriginal && inputPath != outputPath {
		if err := os.Remove(inputPath); err != nil {
			return fmt.Errorf("failed to remove original file: %w", err)
		}
	}

	return nil
}
