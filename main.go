package main

import (
	"context"
	"embed"
	"runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"go-tools/backend/app"
	"go-tools/backend/processor"
)

//go:embed all:frontend/dist
var assets embed.FS

// 应用版本号
var version = "1.0.7"

func main() {
	// Create an instance of the app structure
	application := app.NewApp()
	application.SetVersion(version)

	// Create processor instances
	jsonProcessor := processor.NewJsonProcessor()
	xmlProcessor := processor.NewXMLProcessor()
	charlesGenerator := processor.NewCharlesGenerator()

	// 检测操作系统
	isMacOS := runtime.GOOS == "darwin"

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "dev tools",
		Width:  1000,
		Height: 700,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup: func(ctx context.Context) {
			application.Startup(ctx)
		},
		Bind: []interface{}{
			application,
			jsonProcessor,
			xmlProcessor,
			charlesGenerator,
		},
		Frameless: !isMacOS, // macOS使用有边框窗口，其他平台无边框
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
		Windows: &windows.Options{
			WebviewIsTransparent:              false,
			WindowIsTranslucent:               false,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
