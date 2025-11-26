package file

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Filter 文件过滤器
type Filter struct {
	DisplayName string `json:"displayName"`
	Pattern     string `json:"pattern"`
}

// SaveOptions 保存文件选项
type SaveOptions struct {
	Title           string   `json:"title"`
	DefaultFilename string   `json:"defaultFilename"`
	Filters         []Filter `json:"filters"`
}

// Saver 文件保存器
type Saver struct {
	ctx context.Context
}

// NewSaver 创建文件保存器
func NewSaver(ctx context.Context) *Saver {
	return &Saver{ctx: ctx}
}

// Save 通用保存文件方法，支持base64数据和普通文本
func (s *Saver) Save(content string, options SaveOptions, isBase64 bool) (string, error) {
	// 转换过滤器格式
	var filters []runtime.FileFilter
	for _, filter := range options.Filters {
		filters = append(filters, runtime.FileFilter{
			DisplayName: filter.DisplayName,
			Pattern:     filter.Pattern,
		})
	}

	// 弹出保存对话框
	fileName, err := runtime.SaveFileDialog(s.ctx, runtime.SaveDialogOptions{
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
