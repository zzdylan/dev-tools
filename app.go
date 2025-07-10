package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

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

// FileFilter 文件过滤器
type FileFilter struct {
	DisplayName string `json:"displayName"`
	Pattern     string `json:"pattern"`
}

// SaveFileOptions 保存文件选项
type SaveFileOptions struct {
	Title           string       `json:"title"`
	DefaultFilename string       `json:"defaultFilename"`
	Filters         []FileFilter `json:"filters"`
}

// SaveFile 通用保存文件方法，支持base64数据和普通文本
func (a *App) SaveFile(content string, options SaveFileOptions, isBase64 bool) (string, error) {
	// 转换过滤器格式
	var filters []runtime.FileFilter
	for _, filter := range options.Filters {
		filters = append(filters, runtime.FileFilter{
			DisplayName: filter.DisplayName,
			Pattern:     filter.Pattern,
		})
	}

	// 弹出保存对话框
	fileName, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: options.DefaultFilename,
		Title:          options.Title,
		Filters:        filters,
	})
	
	if err != nil {
		return "", fmt.Errorf("用户取消保存")
	}
	
	if fileName == "" {
		return "", fmt.Errorf("用户取消保存")
	}

	var data []byte
	if isBase64 {
		// 解码base64数据
		data, err = base64.StdEncoding.DecodeString(content)
		if err != nil {
			return "", fmt.Errorf("解码base64数据失败: %v", err)
		}
	} else {
		// 普通文本
		data = []byte(content)
	}

	// 写入文件
	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		return "", fmt.Errorf("写入文件失败: %v", err)
	}

	return fileName, nil
}
