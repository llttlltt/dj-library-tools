package main

import (
	"os"

	"github.com/llttlltt/dj-library-tools/internal/cli"

	// Register Providers
	_ "github.com/llttlltt/dj-library-tools/internal/provider/m3u"
	_ "github.com/llttlltt/dj-library-tools/internal/provider/plex"
	_ "github.com/llttlltt/dj-library-tools/internal/provider/rb"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
