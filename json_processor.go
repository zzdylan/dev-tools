package main

import (
	"bytes"
	"container/list"
	"context"
	"encoding/json"
	"strconv"
	"strings"
)

// JsonProcessor 处理 JSON 相关的功能
type JsonProcessor struct {
	ctx context.Context
}

func NewJsonProcessor() *JsonProcessor {
	return &JsonProcessor{}
}

// 自定义的 JSON 解码器，用于保持对象键的顺序
type OrderedJSONDecoder struct {
	order *list.List
	data  map[string]interface{}
}

func NewOrderedJSONDecoder() *OrderedJSONDecoder {
	return &OrderedJSONDecoder{
		order: list.New(),
		data:  make(map[string]interface{}),
	}
}

func (d *OrderedJSONDecoder) Decode(data []byte) error {
	// 先用标准解码器解析数据
	if err := json.Unmarshal(data, &d.data); err != nil {
		return err
	}

	// 使用原始数据来获取键的顺序
	dec := json.NewDecoder(bytes.NewReader(data))

	// 跳过开始的对象标记
	if _, err := dec.Token(); err != nil {
		return err
	}

	// 读取所有键
	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			return err
		}

		if key, ok := t.(string); ok {
			d.order.PushBack(key)
			// 跳过值
			if _, err := dec.Token(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *OrderedJSONDecoder) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('{')

	first := true
	for e := d.order.Front(); e != nil; e = e.Next() {
		key := e.Value.(string)
		value := d.data[key]

		if !first {
			buf.WriteByte(',')
		}
		first = false

		// 写入键
		keyJSON, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		buf.Write(keyJSON)
		buf.WriteByte(':')

		// 写入值
		valueJSON, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		buf.Write(valueJSON)
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// decodeUnicode 解码 Unicode 转义序列
func decodeUnicode(s string) string {
	var result strings.Builder
	chars := []rune(s)
	for i := 0; i < len(chars); i++ {
		// 检查是否是非转义的 Unicode 序列（前面不是反斜杠）
		if i+5 < len(chars) && chars[i] == '\\' && chars[i+1] == 'u' &&
			(i == 0 || chars[i-1] != '\\') {
			// 尝试解析 4 位十六进制数
			hex := string(chars[i+2 : i+6])
			// 验证是否是有效的十六进制字符
			if isValidHex(hex) {
				num, err := parseInt(hex, 16, 16)
				if err == nil {
					result.WriteRune(rune(num))
					i += 5 // 跳过已处理的字符
					continue
				}
			}
		}
		result.WriteRune(chars[i])
	}
	return result.String()
}

// isValidHex 检查是否是有效的十六进制字符串
func isValidHex(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// parseInt 解析十六进制字符串
func parseInt(s string, base int, bitSize int) (int64, error) {
	return strconv.ParseInt(s, base, bitSize)
}

// FormatJson 格式化 JSON 字符串
func (j *JsonProcessor) FormatJson(jsonStr string, autoDecodeUnicode bool) string {
	// 使用 RawMessage 来保持原始 JSON 的顺序
	var raw json.RawMessage
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		return err.Error()
	}

	// 如果需要解码 Unicode
	if autoDecodeUnicode {
		jsonStr = decodeUnicode(string(raw))
		if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
			return err.Error()
		}
	}

	// 使用 bytes.Buffer 来格式化
	var out bytes.Buffer
	enc := json.NewEncoder(&out)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false) // 不转义 HTML 字符
	if err := enc.Encode(raw); err != nil {
		return err.Error()
	}

	return out.String()
}

// CompressJson 压缩 JSON 字符串
func (j *JsonProcessor) CompressJson(jsonStr string, autoDecodeUnicode bool) string {
	// 使用 RawMessage 来保持原始 JSON 的顺序
	var raw json.RawMessage
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		return err.Error()
	}

	// 如果需要解码 Unicode
	if autoDecodeUnicode {
		jsonStr = decodeUnicode(string(raw))
		if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
			return err.Error()
		}
	}

	// 使用 bytes.Buffer 来压缩
	var out bytes.Buffer
	enc := json.NewEncoder(&out)
	enc.SetEscapeHTML(false) // 不转义 HTML 字符
	if err := enc.Encode(raw); err != nil {
		return err.Error()
	}

	return out.String()
}
