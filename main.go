package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"

	"g1sti/internal/app"
	"g1sti/internal/paths"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	application := app.New()

	if logPath, err := paths.AppLogPath(); err == nil {
		if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err == nil {
			if file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644); err == nil {
				defer file.Close()
				log.SetOutput(file)
			}
		}
	}

	err := wails.Run(&options.App{
		Title:  "G1STI",
		Width:  1024,
		Height: 720,
		Assets: assets,
		Bind: []interface{}{
			application,
		},
		OnStartup:  application.Startup,
		OnShutdown: application.Shutdown,
		Windows: &windows.Options{
			WebviewIsTransparent: false,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
