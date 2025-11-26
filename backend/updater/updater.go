package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	goruntime "runtime"
	"strings"
	"time"
)

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

// Updater 更新检查器
type Updater struct {
	currentVersion string
	customURL      string
}

// NewUpdater 创建更新检查器
func NewUpdater(currentVersion string, customURL string) *Updater {
	// 确保版本号有 v 前缀
	if !strings.HasPrefix(currentVersion, "v") {
		currentVersion = "v" + currentVersion
	}
	return &Updater{
		currentVersion: currentVersion,
		customURL:      customURL,
	}
}

// CheckForUpdate 检查更新（优先使用自定义服务器，失败时降级到 GitHub API）
func (u *Updater) CheckForUpdate(owner, repo string) (*UpdateInfo, error) {
	// 优先尝试从自定义服务器获取
	if u.customURL != "" {
		updateInfo, err := u.checkCustomServer()
		if err == nil {
			return updateInfo, nil
		}
		// 自定义服务器失败，降级到 GitHub API
		println("自定义服务器获取失败，降级到 GitHub API:", err.Error())
	}

	return u.checkGitHubAPI(owner, repo)
}

// checkCustomServer 从自定义服务器检查更新
func (u *Updater) checkCustomServer() (*UpdateInfo, error) {
	// 发送请求，设置超时时间
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(u.customURL)
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
	currentVer := strings.TrimPrefix(u.currentVersion, "v")
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
		CurrentVersion: u.currentVersion,
		LatestVersion:  versionInfo.Version,
		HasUpdate:      hasUpdate,
		Description:    versionInfo.Description,
		DownloadURL:    downloadURL,
	}, nil
}

// checkGitHubAPI 从 GitHub API 检查更新
func (u *Updater) checkGitHubAPI(owner, repo string) (*UpdateInfo, error) {
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
	currentVer := strings.TrimPrefix(u.currentVersion, "v")
	latestVer := strings.TrimPrefix(release.TagName, "v")

	hasUpdate := compareVersion(latestVer, currentVer) > 0

	return &UpdateInfo{
		CurrentVersion: u.currentVersion,
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
