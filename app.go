package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// MinimizeWindow minimizes the window
func (a *App) MinimizeWindow() {
	runtime.WindowMinimise(a.ctx)
}

// CloseWindow closes the window
func (a *App) CloseWindow() {
	runtime.Quit(a.ctx)
}

// ToggleFullscreen toggles the window fullscreen state
func (a *App) ToggleFullscreen() {
	fmt.Println("ToggleFullscreen")
	if runtime.WindowIsFullscreen(a.ctx) {
		fmt.Println("Unfullscreen")
		runtime.WindowUnfullscreen(a.ctx)
	} else {
		fmt.Println("Fullscreen")
		runtime.WindowFullscreen(a.ctx)
	}
}

// IsFullscreen returns whether the window is fullscreen
func (a *App) IsFullscreen() bool {
	return runtime.WindowIsFullscreen(a.ctx)
}
