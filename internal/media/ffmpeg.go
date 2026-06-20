package media

import (
	"fmt"
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
	// Apply path mappings to source
	source = t.ApplyPathMap(source)

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
