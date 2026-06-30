package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/llttlltt/dj-library-tools/internal/ui/cli"

	// Register Providers
	_ "github.com/llttlltt/dj-library-tools/internal/providers/m3u"
	_ "github.com/llttlltt/dj-library-tools/internal/providers/plex"
	_ "github.com/llttlltt/dj-library-tools/internal/providers/rekordbox"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := cli.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
