package main

import (
	"log"

	"g1sti/internal/app"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

func main() {
	application := app.New()

	err := wails.Run(&options.App{
		Title:  "G1STI",
		Width:  1024,
		Height: 720,
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
