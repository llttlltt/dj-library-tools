package main

import (
	"embed"
	"time"

	"github.com/llttlltt/dj-library-tools/internal/ui/gui"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := gui.NewApp()

	// Perform background update check
	go func() {
		// Wait a few seconds for app startup to settle
		time.Sleep(5 * time.Second)
		info, err := app.CheckForUpdate(false)
		if err == nil && info != nil && info.Available {
			// Notify frontend (could show a small toast)
			// runtime.EventsEmit(app.Context(), "update-available", info)
		}
	}()

	err := wails.Run(&options.App{
		Title:     "DJ Library Tools",
		Width:     1280,
		Height:    800,
		MinWidth:  1024,
		MinHeight: 700,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.Startup,
		Bind:      []interface{}{app},
	})
	if err != nil {
		panic(err)
	}
}
