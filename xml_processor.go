package main

import (
	"bytes"
	"encoding/xml"
	"io"
	"strings"
)

type XMLProcessor struct{}

func NewXMLProcessor() *XMLProcessor {
	return &XMLProcessor{}
}

func (x *XMLProcessor) FormatXML(input string) (string, error) {
	var buf bytes.Buffer
	decoder := xml.NewDecoder(strings.NewReader(strings.TrimSpace(input)))
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "  ")

	// 处理第一个元素
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		// 跳过空白字符
		if charData, ok := token.(xml.CharData); ok {
			if len(strings.TrimSpace(string(charData))) == 0 {
				continue
			}
		}

		// 处理 XML 声明
		if procInst, ok := token.(xml.ProcInst); ok && procInst.Target == "xml" {
			encoder.EncodeToken(token)
			buf.WriteString("\n")
			continue
		}

		if err := encoder.EncodeToken(token); err != nil {
			return "", err
		}
	}

	if err := encoder.Flush(); err != nil {
		return "", err
	}

	// 移除多余的空行
	lines := strings.Split(buf.String(), "\n")
	var result []string
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line != "" || (i > 0 && i < len(lines)-1) {
			result = append(result, lines[i])
		}
	}

	// 去掉编码器自动添加的换行符
	finalResult := strings.Join(result, "\n")
	return strings.TrimSuffix(finalResult, "\n"), nil
}

func (x *XMLProcessor) CompressXML(input string) (string, error) {
	var buf bytes.Buffer
	decoder := xml.NewDecoder(strings.NewReader(strings.TrimSpace(input)))
	encoder := xml.NewEncoder(&buf)

	// 禁用自动换行
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		// 跳过空白字符
		if _, ok := token.(xml.CharData); ok {
			str := strings.TrimSpace(string(token.(xml.CharData)))
			if str == "" {
				continue
			}
		}

		if err := encoder.EncodeToken(token); err != nil {
			return "", err
		}
	}

	if err := encoder.Flush(); err != nil {
		return "", err
	}

	// 移除多余的空白字符
	result := strings.ReplaceAll(buf.String(), ">  <", "><")
	result = strings.ReplaceAll(result, "> <", "><")
	// 去掉编码器自动添加的换行符
	return strings.TrimSuffix(result, "\n"), nil
}
