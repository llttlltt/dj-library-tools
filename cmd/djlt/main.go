package main

import (
	"os"

	"github.com/llttlltt/dj-library-tools/internal/cli"

	// Register Providers
	_ "github.com/llttlltt/dj-library-tools/internal/providers/m3u"
	_ "github.com/llttlltt/dj-library-tools/internal/providers/plex"
	_ "github.com/llttlltt/dj-library-tools/internal/providers/rekordbox"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
