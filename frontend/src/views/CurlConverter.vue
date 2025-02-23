<template>
  <div class="curl-converter">
    <div class="converter-section">
      <div class="input-group">
        <div class="label">cURL 命令</div>
        <div class="input-with-buttons">
          <textarea
            v-model="curlCommand"
            placeholder="输入 cURL 命令，例如: curl https://api.example.com -H 'Authorization: Bearer token'"
            class="text-area"
          ></textarea>
          <div class="button-group">
            <select v-model="targetLang" class="lang-select">
              <option value="javascript">JavaScript (Fetch)</option>
              <option value="python">Python (requests)</option>
              <option value="go">Go (net/http)</option>
              <option value="go-resty">Go (resty)</option>
              <option value="java">Java</option>
              <option value="php">PHP</option>
              <option value="nodejs">Node.js (Axios)</option>
            </select>
            <button class="tool-btn" @click="convert">转换</button>
            <button class="tool-btn" @click="clear">清空</button>
            <button class="tool-btn" @click="copy(curlCommand)">复制</button>
          </div>
        </div>
      </div>

      <div class="input-group">
        <div class="label">转换结果</div>
        <div class="input-with-buttons">
          <div class="editor-container">
            <MonacoEditor
              ref="monacoEditor"
              :value="convertedCode"
              :options="editorOptions"
              :language="getEditorLanguage"
              theme="vs"
            />
          </div>
          <div class="button-group">
            <button class="tool-btn" @click="copy(convertedCode)">
              复制结果
            </button>
          </div>
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
    const url = cmd.match(/curl\s+['"]?([^'"]\S+)['"]?/)?.[1]
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

  return `fetch("${url}", {
  method: "${method}",
  headers: ${JSON.stringify(headerObj, null, 2)},
  ${data ? `body: '${data}',` : ''}
})
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

  return `import requests

response = requests.${method.toLowerCase()}(
    "${url}",
    headers={
${headerStr}
    }${data ? `,\n    data='${data}'` : ''}
)

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
\t"io/ioutil"
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

\t// 读取响应
\tbody, err := ioutil.ReadAll(resp.Body)
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
\tresp, err := client.R().
\t\tSetHeaders(headers).
${data ? `\t\tSetBody(\`${data}\`).` : ''}\t\t${method.toLowerCase()}("${url}")

\tif err != nil {
\t\tfmt.Printf("请求失败: %s\\n", err)
\t\treturn
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
    .map((h) => `    .addHeader("${h.split(': ')[0]}", "${h.split(': ')[1]}")`)
    .join('\n')

  return `import okhttp3.*;
import java.io.IOException;

public class HttpClient {
    public static void main(String[] args) {
        OkHttpClient client = new OkHttpClient();

        Request request = new Request.Builder()
            .url("${url}")
            .method("${method}", ${
    data
      ? `RequestBody.create(MediaType.parse("application/json"), "${data}")`
      : 'null'
  })
${headerStr}
            .build();

        try (Response response = client.newCall(request).execute()) {
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
    .map((h) => `"${h.split(': ')[0]}: ${h.split(': ')[1]}"`)
    .join(',\n    ')

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

curl_close($curl);

if ($err) {
    echo "cURL Error #:" . $err;
} else {
    echo $response;
}`
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

  return `const axios = require('axios');

const config = {
  method: '${method.toLowerCase()}',
  url: '${url}',
  headers: ${JSON.stringify(headerObj, null, 2)}${
    data ? `,\n  data: ${JSON.stringify(JSON.parse(data), null, 2)}` : ''
  }
};

axios(config)
  .then(function (response) {
    console.log(JSON.stringify(response.data, null, 2));
  })
  .catch(function (error) {
    console.log(error);
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
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.converter-section {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.input-group {
  margin-bottom: 20px;
}

.input-group:last-child {
  margin-bottom: 0;
}

.label {
  font-size: 14px;
  color: #374151;
  margin-bottom: 8px;
}

.input-with-buttons {
  display: flex;
  gap: 12px;
}

.text-area {
  flex: 1;
  min-height: 120px;
  padding: 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  font-family: monospace;
  resize: vertical;
}

.code-area {
  background: #f8f9fa;
}

.text-area:focus {
  outline: none;
  border-color: #60a5fa;
  box-shadow: 0 0 0 2px rgba(96, 165, 250, 0.2);
}

.button-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tool-btn {
  padding: 8px 16px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  background: white;
  color: #374151;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.tool-btn:hover {
  background: #f3f4f6;
  border-color: #9ca3af;
}

.lang-select {
  padding: 8px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  background: white;
  color: #374151;
  font-size: 14px;
  cursor: pointer;
  outline: none;
}

.lang-select:focus {
  border-color: #60a5fa;
  box-shadow: 0 0 0 2px rgba(96, 165, 250, 0.2);
}

.editor-container {
  flex: 1;
  height: 300px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  overflow: hidden;
}
</style>
