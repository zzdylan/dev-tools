package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

// WindowSettings 窗口设置
type WindowSettings struct {
	Width     int  `json:"width"`
	Height    int  `json:"height"`
	X         int  `json:"x"`
	Y         int  `json:"y"`
	Maximised bool `json:"maximised"`
}

// GetUserDataDir 获取应用数据目录
func GetUserDataDir(appName string) string {
	// 获取用户主目录
	homeDir, _ := os.UserHomeDir()

	// 定义应用的数据目录（根据平台调整路径）
	var appDataDir string
	switch runtime.GOOS {
	case "windows":
		appDataDir = filepath.Join(homeDir, "AppData", "Roaming", appName)
	case "darwin":
		appDataDir = filepath.Join(homeDir, "Library", "Application Support", appName)
	default:
		// Linux 使用 XDG 标准
		appDataDir = filepath.Join(homeDir, ".config", appName)
	}

	// 如果目录不存在，则创建
	if _, err := os.Stat(appDataDir); os.IsNotExist(err) {
		os.MkdirAll(appDataDir, 0755)
	}

	return appDataDir
}

// GetConfigPath 获取配置文件路径
func GetConfigPath() (string, error) {
	configDir := GetUserDataDir("dev-tools")
	return filepath.Join(configDir, "window.json"), nil
}

// LoadWindowSettings 加载窗口设置
func LoadWindowSettings() (*WindowSettings, error) {
	// 默认窗口设置
	defaultSettings := &WindowSettings{
		Width:     1000,
		Height:    700,
		X:         -1,
		Y:         -1,
		Maximised: false,
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return defaultSettings, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		// 任何读取错误都返回默认值
		return defaultSettings, nil
	}

	var settings WindowSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		// JSON 解析失败也返回默认值
		return defaultSettings, nil
	}

	// 验证加载的设置是否有效
	if settings.Width <= 0 || settings.Height <= 0 {
		settings.Width = defaultSettings.Width
		settings.Height = defaultSettings.Height
	}

	return &settings, nil
}

// SaveWindowSettings 保存窗口设置
func SaveWindowSettings(settings *WindowSettings) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
