package main

import (
	"bytes"
	"context"
	"encoding/json"
)

// JsonProcessor 处理 JSON 相关的功能
type JsonProcessor struct {
	ctx context.Context
}

func NewJsonProcessor() *JsonProcessor {
	return &JsonProcessor{}
}

// FormatJson 格式化 JSON 字符串
func (j *JsonProcessor) FormatJson(jsonStr string, autoDecodeUnicode bool) string {
	if !autoDecodeUnicode {
		// 如果不需要解码，使用 bytes.Buffer 来格式化，保持原有顺序
		var out bytes.Buffer
		err := json.Indent(&out, []byte(jsonStr), "", "  ")
		if err != nil {
			return err.Error()
		}
		return out.String()
	}

	// 需要解码时，使用 map 来保持顺序
	var obj map[string]interface{}
	d := json.NewDecoder(bytes.NewReader([]byte(jsonStr)))
	d.UseNumber() // 使用 Number 类型来保持数字精度
	if err := d.Decode(&obj); err != nil {
		return err.Error()
	}

	// 使用 bytes.Buffer 来格式化
	var out bytes.Buffer
	enc := json.NewEncoder(&out)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false) // 不转义 HTML 字符
	if err := enc.Encode(obj); err != nil {
		return err.Error()
	}

	return out.String()
}

// CompressJson 压缩 JSON 字符串
func (j *JsonProcessor) CompressJson(jsonStr string, autoDecodeUnicode bool) string {
	if !autoDecodeUnicode {
		// 如果不需要解码，使用 bytes.Buffer 来压缩，保持原有顺序
		var out bytes.Buffer
		err := json.Compact(&out, []byte(jsonStr))
		if err != nil {
			return err.Error()
		}
		return out.String()
	}

	// 需要解码时，使用 map 来保持顺序
	var obj map[string]interface{}
	d := json.NewDecoder(bytes.NewReader([]byte(jsonStr)))
	d.UseNumber() // 使用 Number 类型来保持数字精度
	if err := d.Decode(&obj); err != nil {
		return err.Error()
	}

	// 使用 bytes.Buffer 来压缩
	var out bytes.Buffer
	enc := json.NewEncoder(&out)
	enc.SetEscapeHTML(false) // 不转义 HTML 字符
	if err := enc.Encode(obj); err != nil {
		return err.Error()
	}

	return out.String()
}
