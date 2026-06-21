package media

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Transcoder struct {
	Config *Config
}

func NewTranscoder(cfg *Config) *Transcoder {
	return &Transcoder{Config: cfg}
}

func (t *Transcoder) Transcode(source, dest string) error {
	// 1. Check if FFmpeg exists
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found in PATH. Please install FFmpeg to use transcoding features")
	}

	// 2. Check if destination already exists (Smart Skip)
	if info, err := os.Stat(dest); err == nil {
		if info.Size() > 0 {
			// File exists and is non-empty, we skip for now.
			// In the future, we could check bitrates/metadata to decide if we need to re-transcode.
			return nil
		}
	}

	// Apply path mappings to source
	source = t.ApplyPathMap(source)

	// Ensure source exists
	if _, err := os.Stat(source); err != nil {
		return fmt.Errorf("source file not found: %s", source)
	}

	cmdStr, ok := t.Config.Formats[t.Config.Format]
	if !ok {
		return fmt.Errorf("unsupported format: %s", t.Config.Format)
	}

	// Replace variables in the ffmpeg command string
	cmdStr = strings.ReplaceAll(cmdStr, "$source", source)
	cmdStr = strings.ReplaceAll(cmdStr, "$dest", dest)

	// Split the command string into parts for exec.Command
	// Note: Simple splitting by space. For complex commands with quoted spaces, 
	// we might need a more robust shell-word splitter.
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return fmt.Errorf("invalid command for format %s", t.Config.Format)
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ffmpeg error: %v, output: %s", err, string(output))
	}

	return nil
}
