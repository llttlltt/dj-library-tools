package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func EnforceExtension(path, extension string) string {
	base := strings.TrimSuffix(path, filepath.Ext(path))
	if extension != "" && !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}
	return base + extension
}

func CheckFileOverwrite(path string, forceOverwrite bool) error {
	if _, err := os.Stat(path); err == nil {
		if !forceOverwrite {
			return fmt.Errorf("file '%s' already exists. Use the --force (-f) flag to overwrite", path)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check file status '%s': %w", path, err)
	}
	return nil
}

func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}
