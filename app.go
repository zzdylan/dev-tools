package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx     context.Context
	version string
}

// WindowSettings 窗口设置
type WindowSettings struct {
	Width      int  `json:"width"`
	Height     int  `json:"height"`
	X          int  `json:"x"`
	Y          int  `json:"y"`
	Maximised  bool `json:"maximised"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// getConfigPath 获取配置文件路径
func (a *App) getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(homeDir, ".dev-tools")
	// 确保配置目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(configDir, "window.json"), nil
}

// LoadWindowSettings 加载窗口设置
func (a *App) LoadWindowSettings() (*WindowSettings, error) {
	configPath, err := a.getConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		// 文件不存在时返回默认值
		if os.IsNotExist(err) {
			return &WindowSettings{
				Width:  1000,
				Height: 700,
				X:      -1,
				Y:      -1,
				Maximised: false,
			}, nil
		}
		return nil, err
	}

	var settings WindowSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

// SaveWindowSettings 保存窗口设置
func (a *App) SaveWindowSettings() error {
	configPath, err := a.getConfigPath()
	if err != nil {
		return err
	}

	width, height := runtime.WindowGetSize(a.ctx)
	x, y := runtime.WindowGetPosition(a.ctx)
	maximised := runtime.WindowIsMaximised(a.ctx)

	settings := WindowSettings{
		Width:     width,
		Height:    height,
		X:         x,
		Y:         y,
		Maximised: maximised,
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// ResetWindowSize 重置窗口大小到默认值
func (a *App) ResetWindowSize() error {
	// 设置窗口大小为默认值
	runtime.WindowSetSize(a.ctx, 1000, 700)

	// 居中窗口
	runtime.WindowCenter(a.ctx)

	// 如果窗口是最大化的，先取消最大化
	if runtime.WindowIsMaximised(a.ctx) {
		runtime.WindowUnmaximise(a.ctx)
	}

	return nil
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
		Title:           options.Title,
		Filters:         filters,
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

// SetVersion 设置应用版本
func (a *App) SetVersion(version string) {
	if !strings.HasPrefix(version, "v") {
		a.version = "v" + version
	} else {
		a.version = version
	}
}

// GetVersion 获取应用版本
func (a *App) GetVersion() string {
	return a.version
}

// GitHubRelease GitHub Release 结构
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	HTMLURL string `json:"html_url"`
}

// CustomVersionInfo 自定义服务器版本信息结构
type CustomVersionInfo struct {
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	DownloadURLs map[string]string `json:"downloadUrls"` // 支持多平台下载链接
}

// UpdateInfo 更新信息
type UpdateInfo struct {
	CurrentVersion string `json:"currentVersion"`
	LatestVersion  string `json:"latestVersion"`
	HasUpdate      bool   `json:"hasUpdate"`
	Description    string `json:"description"`
	DownloadURL    string `json:"downloadUrl"`
}

// CheckForUpdate 检查更新（优先使用自定义服务器，失败时降级到 GitHub API）
func (a *App) CheckForUpdate(owner, repo string) (*UpdateInfo, error) {
	// 优先尝试从自定义服务器获取
	customURL := "https://devtools.51godream.com/version.json"
	updateInfo, err := a.checkCustomServer(customURL)
	if err == nil {
		return updateInfo, nil
	}

	// 自定义服务器失败，降级到 GitHub API
	println("自定义服务器获取失败，降级到 GitHub API:", err.Error())
	return a.checkGitHubAPI(owner, repo)
}

// checkCustomServer 从自定义服务器检查更新
func (a *App) checkCustomServer(url string) (*UpdateInfo, error) {
	// 发送请求，设置超时时间
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求自定义服务器失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("自定义服务器返回错误: %d", resp.StatusCode)
	}

	// 解析响应
	var versionInfo CustomVersionInfo
	if err := json.NewDecoder(resp.Body).Decode(&versionInfo); err != nil {
		return nil, fmt.Errorf("解析自定义服务器响应失败: %v", err)
	}

	// 比较版本
	currentVer := strings.TrimPrefix(a.version, "v")
	latestVer := strings.TrimPrefix(versionInfo.Version, "v")

	hasUpdate := compareVersion(latestVer, currentVer) > 0

	// 根据当前平台选择下载地址
	downloadURL := ""
	platformKey := goruntime.GOOS // "darwin", "windows", "linux"
	if url, ok := versionInfo.DownloadURLs[platformKey]; ok {
		downloadURL = url
	}
	// 如果没有找到对应平台的下载地址，downloadURL 将为空字符串
	// 前端会处理这种情况并降级到 GitHub

	return &UpdateInfo{
		CurrentVersion: a.version,
		LatestVersion:  versionInfo.Version,
		HasUpdate:      hasUpdate,
		Description:    versionInfo.Description,
		DownloadURL:    downloadURL,
	}, nil
}

// checkGitHubAPI 从 GitHub API 检查更新
func (a *App) checkGitHubAPI(owner, repo string) (*UpdateInfo, error) {
	// GitHub API URL
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

	// 发送请求，设置超时时间
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求 GitHub API 失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API 返回错误: %d", resp.StatusCode)
	}

	// 解析响应
	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("解析 GitHub API 响应失败: %v", err)
	}

	// 比较版本
	currentVer := strings.TrimPrefix(a.version, "v")
	latestVer := strings.TrimPrefix(release.TagName, "v")

	hasUpdate := compareVersion(latestVer, currentVer) > 0

	return &UpdateInfo{
		CurrentVersion: a.version,
		LatestVersion:  release.TagName,
		HasUpdate:      hasUpdate,
		Description:    release.Body,
		DownloadURL:    release.HTMLURL,
	}, nil
}

// compareVersion 比较版本号
// 返回: 1 表示 v1 > v2, -1 表示 v1 < v2, 0 表示相等
func compareVersion(v1, v2 string) int {
	// 移除 'v' 前缀
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	// 分割版本号
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	// 比较每个部分
	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var n1, n2 int
		if i < len(parts1) {
			fmt.Sscanf(parts1[i], "%d", &n1)
		}
		if i < len(parts2) {
			fmt.Sscanf(parts2[i], "%d", &n2)
		}

		if n1 > n2 {
			return 1
		} else if n1 < n2 {
			return -1
		}
	}

	return 0
}
