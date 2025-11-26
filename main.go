package main

import (
	"context"
	"embed"
	"runtime"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	runtime2 "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

// 应用版本号
var version = "1.0.3"

func main() {
	// Create an instance of the app structure
	app := NewApp()
	app.SetVersion(version)

	// Create JsonProcessor instance
	jsonProcessor := NewJsonProcessor()
	xmlProcessor := NewXMLProcessor()

	// 检测操作系统
	isMacOS := runtime.GOOS == "darwin"

	// 加载窗口设置
	windowSettings, err := app.LoadWindowSettings()
	if err != nil {
		println("Error loading window settings:", err.Error())
		windowSettings = &WindowSettings{
			Width:     1000,
			Height:    700,
			X:         -1,
			Y:         -1,
			Maximised: false,
		}
	}

	// 确定窗口启动状态
	windowStartState := options.Normal
	if windowSettings.Maximised {
		windowStartState = options.Maximised
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:            "dev tools",
		Width:            windowSettings.Width,
		Height:           windowSettings.Height,
		WindowStartState: windowStartState,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)

			// 启动定时器，每2秒检查并保存窗口设置
			go func() {
				ticker := time.NewTicker(2 * time.Second)
				defer ticker.Stop()
				for range ticker.C {
					if err := app.SaveWindowSettings(); err != nil {
						println("Error saving window settings:", err.Error())
					}
				}
			}()
		},
		OnDomReady: func(ctx context.Context) {
			// 恢复窗口位置（只在保存了有效位置时）
			if windowSettings.X >= 0 && windowSettings.Y >= 0 {
				runtime2.WindowSetPosition(ctx, windowSettings.X, windowSettings.Y)
			}

			// 发送初始窗口状态
			runtime2.EventsEmit(ctx, "window_changed", map[string]interface{}{
				"fullscreen": runtime2.WindowIsFullscreen(ctx),
				"maximised":  runtime2.WindowIsMaximised(ctx),
			})
		},
		OnBeforeClose: func(ctx context.Context) bool {
			// 保存窗口设置
			if err := app.SaveWindowSettings(); err != nil {
				println("Error saving window settings:", err.Error())
			}
			return false
		},
		Bind: []interface{}{
			app,
			jsonProcessor,
			xmlProcessor,
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
