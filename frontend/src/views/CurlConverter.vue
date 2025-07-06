<template>
  <div class="curl-converter">
    <div class="converter-container">
      <!-- 输入区域 -->
      <div class="input-section">
        <div class="section-header">
          <h3>cURL 命令</h3>
          <div class="header-controls">
            <select v-model="targetLang" class="lang-select">
              <option value="javascript">JavaScript</option>
              <option value="python">Python</option>
              <option value="go">Go</option>
              <option value="go-resty">Go (Resty)</option>
              <option value="java">Java</option>
              <option value="php">PHP</option>
              <option value="nodejs">Node.js</option>
            </select>
            <button class="btn-convert" @click="convert">转换</button>
          </div>
        </div>
        <div class="input-wrapper">
          <textarea
            v-model="curlCommand"
            placeholder="粘贴 cURL 命令..."
            class="curl-input"
          ></textarea>
          <div class="input-actions">
            <button class="btn-clear" @click="clear">清空</button>
            <button class="btn-copy" @click="copy(curlCommand)">复制</button>
          </div>
        </div>
      </div>

      <!-- 输出区域 -->
      <div class="output-section">
        <div class="section-header">
          <h3>转换结果</h3>
          <button class="btn-copy" @click="copy(convertedCode)">复制结果</button>
        </div>
        <div class="output-wrapper">
          <MonacoEditor
            ref="monacoEditor"
            :value="convertedCode"
            :options="editorOptions"
            :language="getEditorLanguage"
            theme="vs"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useClipboard } from '@vueuse/core'
import { useToolsStore } from '../stores/tools'
import { storeToRefs } from 'pinia'
import MonacoEditor from 'monaco-editor-vue3'

const { copy: copyToClipboard } = useClipboard()
const store = useToolsStore()
const { curlConverter } = storeToRefs(store)

const curlCommand = computed({
  get: () => curlConverter.value.curlCommand,
  set: (val) => (curlConverter.value.curlCommand = val),
})

const targetLang = computed({
  get: () => curlConverter.value.targetLang,
  set: (val) => (curlConverter.value.targetLang = val),
})

const convertedCode = computed({
  get: () => curlConverter.value.convertedCode,
  set: (val) => (curlConverter.value.convertedCode = val),
})

const monacoEditor = ref()

const editorOptions = {
  fontSize: 14,
  tabSize: 2,
  minimap: {
    enabled: false,
  },
  scrollBeyondLastLine: true,
  automaticLayout: true,
  wordWrap: 'on',
  lineNumbers: 'on',
  roundedSelection: false,
  renderIndentGuides: true,
  formatOnPaste: false,
  formatOnType: false,
  readOnly: true,
}

const getEditorLanguage = computed(() => {
  switch (targetLang.value) {
    case 'javascript':
    case 'nodejs':
      return 'javascript'
    case 'python':
      return 'python'
    case 'go':
      return 'go'
    case 'java':
      return 'java'
    case 'php':
      return 'php'
    default:
      return 'plaintext'
  }
})

const convert = () => {
  try {
    if (!curlCommand.value.trim()) {
      ElMessage.warning('请输入 cURL 命令')
      return
    }

    // 解析 curl 命令
    const cmd = curlCommand.value.trim()
    // 修复URL匹配正则，排除选项参数
    const url = cmd.match(/curl\s+(?:-\w+\s+\w+\s+)*['"]?([^'"]\S+)['"]?/)?.[1] || 
                cmd.match(/curl\s+.*?['"]?(https?:\/\/[^'"\s]+)['"]?/)?.[1]
    const headers = Array.from(cmd.matchAll(/-H\s+['"]([^'"]+)['"]/g)).map(
      (m) => m[1]
    )
    const method = cmd.match(/-X\s+(\w+)/)?.[1] || 'GET'
    const data = cmd.match(/-d\s+['"]([^'"]+)['"]/)?.[1]

    if (!url) {
      throw new Error('无效的 cURL 命令')
    }

    // 根据目标语言生成代码
    let code = ''
    switch (targetLang.value) {
      case 'javascript':
        code = generateJavaScriptCode(url, headers, method, data)
        break
      case 'python':
        code = generatePythonCode(url, headers, method, data)
        break
      case 'go':
        code = generateGoCode(url, headers, method, data)
        break
      case 'go-resty':
        code = generateGoRestyCode(url, headers, method, data)
        break
      case 'java':
        code = generateJavaCode(url, headers, method, data)
        break
      case 'php':
        code = generatePHPCode(url, headers, method, data)
        break
      case 'nodejs':
        code = generateNodejsCode(url, headers, method, data)
        break
    }

    convertedCode.value = code
    ElMessage.success('转换成功')
  } catch (error) {
    ElMessage.error(
      '转换失败：' + (error instanceof Error ? error.message : '未知错误')
    )
  }
}

const generateJavaScriptCode = (
  url: string,
  headers: string[],
  method: string,
  data?: string
) => {
  const headerObj = headers.reduce((acc, h) => {
    const [key, value] = h.split(': ')
    acc[key] = value
    return acc
  }, {} as Record<string, string>)

  const config = {
    method: method,
    headers: headerObj,
    ...(data && { body: data })
  }

  return `fetch("${url}", ${JSON.stringify(config, null, 2)})
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));`
}

const generatePythonCode = (
  url: string,
  headers: string[],
  method: string,
  data?: string
) => {
  const headerStr = headers
    .map((h) => `    "${h.split(': ')[0]}": "${h.split(': ')[1]}"`)
    .join(',\n')

  // 检查是否是JSON数据
  const isJson = headers.some(h => h.toLowerCase().includes('content-type') && h.toLowerCase().includes('json'))
  const dataParam = data ? (isJson ? `json=${data}` : `data='${data}'`) : ''

  return `import requests

response = requests.${method.toLowerCase()}(
    "${url}",
    headers={
${headerStr}
    }${dataParam ? `,\n    ${dataParam}` : ''}
)

response.raise_for_status()  # 检查响应状态码
print(response.json())`
}

const generateGoCode = (
  url: string,
  headers: string[],
  method: string,
  data?: string
) => {
  const headerStr = headers
    .map((h) => `\t"${h.split(': ')[0]}": "${h.split(': ')[1]}"`)
    .join(',\n')

  return `package main

import (
\t"fmt"
\t"io"
\t"net/http"
\t"strings"
)

func main() {
\t// 创建请求
\tclient := &http.Client{}
${data ? `\tpayload := strings.NewReader(\`${data}\`)` : ''}
\treq, err := http.NewRequest("${method}", "${url}", ${
    data ? 'payload' : 'nil'
  })
\tif err != nil {
\t\tfmt.Printf("创建请求失败: %s\\n", err)
\t\treturn
\t}

\t// 设置请求头
\theaders := map[string]string{
${headerStr}
\t}
\tfor key, value := range headers {
\t\treq.Header.Set(key, value)
\t}

\t// 发送请求
\tresp, err := client.Do(req)
\tif err != nil {
\t\tfmt.Printf("请求失败: %s\\n", err)
\t\treturn
\t}
\tdefer resp.Body.Close()

\t// 检查响应状态码
\tif resp.StatusCode >= 400 {
\t\tfmt.Printf("请求失败，状态码: %d\\n", resp.StatusCode)
\t}

\t// 读取响应
\tbody, err := io.ReadAll(resp.Body)
\tif err != nil {
\t\tfmt.Printf("读取响应失败: %s\\n", err)
\t\treturn
\t}

\tfmt.Println(string(body))
}`
}

const generateGoRestyCode = (
  url: string,
  headers: string[],
  method: string,
  data?: string
) => {
  const headerStr = headers
    .map((h) => `\t"${h.split(': ')[0]}": "${h.split(': ')[1]}"`)
    .join(',\n')

  return `package main

import (
\t"fmt"
\t"github.com/go-resty/resty/v2"
)

func main() {
\t// 创建客户端
\tclient := resty.New()

\t// 设置请求头
\theaders := map[string]string{
${headerStr}
\t}

\t// 发送请求
\treq := client.R().SetHeaders(headers)
${data ? `\treq.SetBody(\`${data}\`)` : ''}
\tresp, err := req.${method.charAt(0).toUpperCase() + method.slice(1).toLowerCase()}("${url}")

\tif err != nil {
\t\tfmt.Printf("请求失败: %s\\n", err)
\t\treturn
\t}

\t// 检查响应状态码
\tif resp.StatusCode() >= 400 {
\t\tfmt.Printf("请求失败，状态码: %d\\n", resp.StatusCode())
\t}

\tfmt.Println(string(resp.Body()))
}`
}

const generateJavaCode = (
  url: string,
  headers: string[],
  method: string,
  data?: string
) => {
  const headerStr = headers
    .map((h) => `            .addHeader("${h.split(': ')[0]}", "${h.split(': ')[1]}")`)
    .join('\n')

  // 检查Content-Type来确定MediaType
  const contentType = headers.find(h => h.toLowerCase().startsWith('content-type'))?.split(': ')[1] || 'application/json'

  return `import okhttp3.*;
import java.io.IOException;

public class HttpClient {
    public static void main(String[] args) {
        OkHttpClient client = new OkHttpClient();

        Request request = new Request.Builder()
            .url("${url}")
            .method("${method}", ${
    data
      ? `RequestBody.create(MediaType.parse("${contentType}"), "${data}")`
      : 'null'
  })
${headerStr}
            .build();

        try (Response response = client.newCall(request).execute()) {
            if (!response.isSuccessful()) {
                System.out.println("请求失败，状态码: " + response.code());
            }
            System.out.println(response.body().string());
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}`
}

const generatePHPCode = (
  url: string,
  headers: string[],
  method: string,
  data?: string
) => {
  const headerStr = headers
    .map((h) => `        "${h.split(': ')[0]}: ${h.split(': ')[1]}"`)
    .join(',\n')

  return `<?php
$curl = curl_init();

curl_setopt_array($curl, array(
    CURLOPT_URL => "${url}",
    CURLOPT_RETURNTRANSFER => true,
    CURLOPT_ENCODING => "",
    CURLOPT_MAXREDIRS => 10,
    CURLOPT_TIMEOUT => 30,
    CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
    CURLOPT_CUSTOMREQUEST => "${method}",
    CURLOPT_HTTPHEADER => array(
${headerStr}
    )${data ? `,\n    CURLOPT_POSTFIELDS => '${data}'` : ''}
));

$response = curl_exec($curl);
$err = curl_error($curl);
$httpCode = curl_getinfo($curl, CURLINFO_HTTP_CODE);

curl_close($curl);

if ($err) {
    echo "cURL Error: " . $err;
} else {
    if ($httpCode >= 400) {
        echo "请求失败，状态码: " . $httpCode . "\\n";
    }
    echo $response;
}
?>`
}

const generateNodejsCode = (
  url: string,
  headers: string[],
  method: string,
  data?: string
) => {
  const headerObj = headers.reduce((acc, h) => {
    const [key, value] = h.split(': ')
    acc[key] = value
    return acc
  }, {} as Record<string, string>)

  // 检查是否是JSON数据
  const isJson = headers.some(h => h.toLowerCase().includes('content-type') && h.toLowerCase().includes('json'))
  let dataValue = ''
  if (data) {
    if (isJson) {
      try {
        dataValue = JSON.stringify(JSON.parse(data), null, 2)
      } catch {
        dataValue = `'${data}'`
      }
    } else {
      dataValue = `'${data}'`
    }
  }

  return `const axios = require('axios');

const config = {
  method: '${method.toLowerCase()}',
  url: '${url}',
  headers: ${JSON.stringify(headerObj, null, 2)}${
    data ? `,\n  data: ${dataValue}` : ''
  }
};

axios(config)
  .then(function (response) {
    console.log('状态码:', response.status);
    console.log('响应数据:', JSON.stringify(response.data, null, 2));
  })
  .catch(function (error) {
    if (error.response) {
      console.log('请求失败，状态码:', error.response.status);
      console.log('错误信息:', error.response.data);
    } else {
      console.log('请求错误:', error.message);
    }
  });`
}

const clear = () => {
  curlCommand.value = ''
  convertedCode.value = ''
  const editor = monacoEditor.value?.editor
  if (editor) {
    editor.setValue('')
  }
  ElMessage.success('已清空')
}

const copy = async (text: string) => {
  if (!text) {
    ElMessage.warning('没有可复制的内容')
    return
  }
  if (text === convertedCode.value) {
    const editor = monacoEditor.value?.editor
    if (editor) {
      text = editor.getValue()
    }
  }
  await copyToClipboard(text)
  ElMessage.success('已复制到剪贴板')
}
</script>

<style scoped>
.curl-converter {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
  background: #f8fafc;
  min-height: 100vh;
}

.converter-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
  height: calc(100vh - 60px);
}

.input-section,
.output-section {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  border: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #e2e8f0;
  background: #f8fafc;
}

.section-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #1e293b;
}

.header-controls {
  display: flex;
  gap: 12px;
  align-items: center;
}

.lang-select {
  padding: 8px 12px;
  padding-right: 32px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  background: white;
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e");
  background-position: right 8px center;
  background-repeat: no-repeat;
  background-size: 16px;
  color: #374151;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  outline: none;
  transition: all 0.2s;
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
}

.lang-select:focus {
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%233b82f6' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e");
}

.lang-select:hover {
  border-color: #9ca3af;
}

.btn-convert {
  padding: 8px 16px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-convert:hover {
  background: #2563eb;
}

.btn-copy {
  padding: 6px 12px;
  background: #f1f5f9;
  color: #475569;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-copy:hover {
  background: #e2e8f0;
}

.btn-clear {
  padding: 6px 12px;
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fecaca;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-clear:hover {
  background: #fee2e2;
}

.input-wrapper,
.output-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  position: relative;
}

.curl-input {
  flex: 1;
  padding: 20px;
  border: none;
  outline: none;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
  resize: none;
  background: transparent;
  color: #1e293b;
}

.curl-input::placeholder {
  color: #94a3b8;
}

.input-actions {
  position: absolute;
  bottom: 12px;
  right: 12px;
  display: flex;
  gap: 8px;
}

.output-wrapper {
  position: relative;
}

.output-wrapper .btn-copy {
  position: absolute;
  top: 12px;
  right: 12px;
  z-index: 10;
}

@media (max-width: 1024px) {
  .converter-container {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .curl-converter {
    padding: 16px;
  }
}
</style>
