package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create JsonProcessor instance
	jsonProcessor := NewJsonProcessor()
	xmlProcessor := NewXMLProcessor()

	// Create application with options
	err := wails.Run(&options.App{
		Title: "dev tools",
		Width:  1000,
		Height: 700,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			jsonProcessor,
			xmlProcessor,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
