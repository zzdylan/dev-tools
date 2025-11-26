package app

import (
	"context"
	"fmt"

	"go-tools/backend/config"
	"go-tools/backend/file"
	"go-tools/backend/updater"
	"go-tools/backend/window"
)

// App struct
type App struct {
	ctx            context.Context
	version        string
	windowCtrl     *window.Controller
	fileSaver      *file.Saver
	updater        *updater.Updater
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.windowCtrl = window.NewController(ctx)
	a.fileSaver = file.NewSaver(ctx)
}

// SetVersion 设置应用版本
func (a *App) SetVersion(version string) {
	a.version = version
	// 创建更新检查器
	a.updater = updater.NewUpdater(version, "https://devtools.51godream.com/version.json")
}

// GetVersion 获取应用版本
func (a *App) GetVersion() string {
	return a.version
}

// ==================== 窗口设置 ====================

// LoadWindowSettings 加载窗口设置
func (a *App) LoadWindowSettings() (*config.WindowSettings, error) {
	return config.LoadWindowSettings()
}

// SaveWindowSettings 保存窗口设置
func (a *App) SaveWindowSettings() error {
	width, height := a.windowCtrl.GetSize()
	x, y := a.windowCtrl.GetPosition()
	maximised := a.windowCtrl.IsMaximised()

	settings := &config.WindowSettings{
		Width:     width,
		Height:    height,
		X:         x,
		Y:         y,
		Maximised: maximised,
	}

	return config.SaveWindowSettings(settings)
}

// ResetWindowSize 重置窗口大小到默认值
func (a *App) ResetWindowSize() error {
	return a.windowCtrl.ResetSize()
}

// ==================== 窗口控制 ====================

// MinimizeWindow 最小化窗口
func (a *App) MinimizeWindow() {
	a.windowCtrl.Minimize()
}

// CloseWindow 关闭窗口
func (a *App) CloseWindow() {
	a.windowCtrl.Close()
}

// ToggleFullscreen 切换全屏状态
func (a *App) ToggleFullscreen() {
	a.windowCtrl.ToggleFullscreen()
}

// IsFullscreen 获取全屏状态
func (a *App) IsFullscreen() bool {
	return a.windowCtrl.IsFullscreen()
}

// ==================== 文件保存 ====================

// SaveFile 保存文件
func (a *App) SaveFile(content string, options file.SaveOptions, isBase64 bool) (string, error) {
	return a.fileSaver.Save(content, options, isBase64)
}

// ==================== 更新检查 ====================

// CheckForUpdate 检查更新
func (a *App) CheckForUpdate(owner, repo string) (*updater.UpdateInfo, error) {
	if a.updater == nil {
		return nil, fmt.Errorf("updater not initialized")
	}
	return a.updater.CheckForUpdate(owner, repo)
}

// ==================== 其他 ====================

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetContext returns the app context (for internal use)
func (a *App) GetContext() context.Context {
	return a.ctx
}
