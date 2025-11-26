package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// WindowSettings 窗口设置
type WindowSettings struct {
	Width     int  `json:"width"`
	Height    int  `json:"height"`
	X         int  `json:"x"`
	Y         int  `json:"y"`
	Maximised bool `json:"maximised"`
}

// GetConfigPath 获取配置文件路径
func GetConfigPath() (string, error) {
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
func LoadWindowSettings() (*WindowSettings, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		// 文件不存在时返回默认值
		if os.IsNotExist(err) {
			return &WindowSettings{
				Width:     1000,
				Height:    700,
				X:         -1,
				Y:         -1,
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
