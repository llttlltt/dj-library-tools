package media

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/infra/sys"
)

type Transcoder struct {
	Config *Config
	FS     sys.FileSystem
	Runner sys.Runner
}

func NewTranscoder(cfg *Config) *Transcoder {
	return &Transcoder{
		Config: cfg,
		FS:     sys.OSFileSystem{},
		Runner: sys.RealRunner{},
	}
}

func (t *Transcoder) Transcode(source, dest string) error {
	// 1. Check if FFmpeg exists
	if _, err := t.Runner.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found in PATH. Please install FFmpeg to use transcoding features")
	}

	// 2. Check if destination already exists (Smart Skip)
	if info, err := t.FS.Stat(dest); err == nil {
		if info.Size() > 0 {
			// File exists and is non-empty, we skip for now.
			return nil
		}
	}

	// Apply path mappings to source
	source = t.ApplyPathMap(source)

	// Ensure source exists
	if _, err := t.FS.Stat(source); err != nil {
		return fmt.Errorf("source file not found: %s", source)
	}

	// 3. Resolve command arguments
	args, err := t.resolveCommandArgs(source, dest)
	if err != nil {
		return err
	}

	if output, err := t.Runner.Run(args[0], args[1:]...); err != nil {
		return fmt.Errorf("ffmpeg error: %v, output: %s", err, string(output))
	}

	return nil
}

func (t *Transcoder) resolveCommandArgs(source, dest string) ([]string, error) {
	cmdStr, ok := t.Config.Formats[t.Config.Format]
	if !ok {
		return nil, fmt.Errorf("unsupported format: %s", t.Config.Format)
	}

	// We use a simple replacement for the placeholders, but we MUST NOT split by fields
	// after replacement because that breaks paths with spaces.
	// Instead, we split the template first, then replace tokens in each part.
	rawParts := strings.Fields(cmdStr)
	if len(rawParts) == 0 {
		return nil, fmt.Errorf("invalid command for format %s", t.Config.Format)
	}

	resolvedParts := make([]string, len(rawParts))
	for i, p := range rawParts {
		p = strings.ReplaceAll(p, "$source", source)
		p = strings.ReplaceAll(p, "$dest", dest)
		resolvedParts[i] = p
	}

	return resolvedParts, nil
}
