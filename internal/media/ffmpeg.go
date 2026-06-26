package media

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/sys"
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

	cmdStr, ok := t.Config.Formats[t.Config.Format]
	if !ok {
		return fmt.Errorf("unsupported format: %s", t.Config.Format)
	}

	// Replace variables in the ffmpeg command string
	cmdStr = strings.ReplaceAll(cmdStr, "$source", source)
	cmdStr = strings.ReplaceAll(cmdStr, "$dest", dest)

	// Split the command string into parts
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return fmt.Errorf("invalid command for format %s", t.Config.Format)
	}

	if output, err := t.Runner.Run(parts[0], parts[1:]...); err != nil {
		return fmt.Errorf("ffmpeg error: %v, output: %s", err, string(output))
	}

	return nil
}
