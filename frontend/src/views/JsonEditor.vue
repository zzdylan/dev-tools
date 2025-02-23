<template>
  <div class="json-editor">
    <div class="toolbar">
      <div class="config-wrapper">
        <button
          class="tool-btn config-btn"
          ref="configBtn"
          @click="toggleSettings"
        >
          <span class="tool-icon">⚙️</span>
          配置
        </button>

        <!-- 配置面板 -->
        <div
          v-show="showSettings"
          class="settings-panel"
          ref="settingsPanel"
          v-click-outside="closeSettings"
        >
          <label class="setting-item">
            <input type="checkbox" v-model="settings.autoDecodeUnicode" />
            自动解码 Unicode
          </label>
          <label class="setting-item">
            <input type="checkbox" v-model="settings.removeEscapes" />
            去转义时删除换行符和制表符
          </label>
        </div>
      </div>

      <!-- 功能按钮放右边 -->
      <div class="tools-group">
        <button class="tool-btn" @click="loadSample">
          <span class="tool-icon">📝</span>
          示例
        </button>
        <button class="tool-btn" @click="formatJson">
          <span class="tool-icon">✏️</span>
          格式化
        </button>
        <button class="tool-btn" @click="compressJson">
          <span class="tool-icon">📦</span>
          压缩
        </button>
        <button class="tool-btn" @click="escapeJson">
          <span class="tool-icon">🔒</span>
          转义
        </button>
        <button class="tool-btn" @click="unescapeJson">
          <span class="tool-icon">🔓</span>
          去转义
        </button>
        <button class="tool-btn" @click="copyToClipboard">
          <span class="tool-icon">📋</span>
          复制
        </button>
        <button class="tool-btn" @click="clearContent">
          <span class="tool-icon">🗑️</span>
          清空
        </button>
      </div>
    </div>

    <div class="editor-wrapper">
      <div class="editor-container">
        <MonacoEditor
          ref="monacoEditor"
          :value="code"
          @change="handleChange"
          :options="options"
          language="json"
          theme="vs"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import MonacoEditor from 'monaco-editor-vue3'
import { useClipboard } from '@vueuse/core'
import { ElMessage } from 'element-plus'
import { FormatJson, CompressJson } from '../../wailsjs/go/main/JsonProcessor'
import { onClickOutside } from '@vueuse/core'
import { useToolsStore } from '../stores/tools'
import { storeToRefs } from 'pinia'

const { copy } = useClipboard()
const store = useToolsStore()
const { jsonEditor } = storeToRefs(store)
const code = computed({
  get: () => jsonEditor.value.code,
  set: (val) => (jsonEditor.value.code = val),
})
const settings = computed({
  get: () => jsonEditor.value.settings,
  set: (val) => (jsonEditor.value.settings = val),
})

const error = ref('')

const options = {
  fontSize: 12,
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
  autoIndent: 'none',
  detectIndentation: false,
  insertSpaces: true,
  trimAutoWhitespace: false,
  folding: true,
  foldingStrategy: 'indentation',
  scrollbar: {
    vertical: 'visible',
    horizontal: 'visible',
    verticalScrollbarSize: 14,
    horizontalScrollbarSize: 14,
    alwaysConsumeMouseWheel: true,
  },
  lineDecorationsWidth: 0,
  lineNumbersMinChars: 0,
  glyphMargin: false,
  renderLineHighlight: 'none',
}

const monacoEditor = ref()

// 配置状态
const showSettings = ref(false)
const settingsPanel = ref<HTMLElement | null>(null)

const configBtn = ref<HTMLElement | null>(null)

// 计算面板位置
const panelStyle = computed(() => {
  if (!configBtn.value) return {}
  const btnRect = configBtn.value.getBoundingClientRect()
  return {
    position: 'absolute',
    top: '100%',
    left: '0',
    marginTop: '4px',
  }
})

// 点击外部关闭配置面板
onClickOutside(settingsPanel, () => {
  showSettings.value = false
})

// 切换配置面板
const toggleSettings = () => {
  showSettings.value = !showSettings.value
}

const formatJson = async () => {
  try {
    const editor = monacoEditor.value?.editor
    if (!editor) {
      ElMessage.error('编辑器未准备好')
      return
    }

    const model = editor.getModel()
    if (!model) {
      ElMessage.error('获取内容失败')
      return
    }

    const value = model.getValue()
    if (!value.trim()) {
      ElMessage.error('请输入 JSON 内容')
      return
    }

    const formatted = await FormatJson(value, settings.value.autoDecodeUnicode)
    model.setValue(formatted)
  } catch (err: any) {
    ElMessage.error('格式化失败：' + (err.message || err))
  }
}

const compressJson = async () => {
  try {
    const editor = monacoEditor.value?.editor
    if (!editor) {
      ElMessage.error('编辑器未准备好')
      return
    }

    const model = editor.getModel()
    if (!model) {
      ElMessage.error('获取内容失败')
      return
    }

    const value = model.getValue()
    if (!value.trim()) {
      ElMessage.error('请输入 JSON 内容')
      return
    }

    const compressed = await CompressJson(
      value,
      settings.value.autoDecodeUnicode
    )
    model.setValue(compressed)
  } catch (err: any) {
    ElMessage.error('压缩失败：' + (err.message || err))
  }
}

// 添加新的转义和去转义函数
const escapeJson = async () => {
  try {
    const editor = monacoEditor.value?.editor
    if (!editor) {
      ElMessage.error('编辑器未准备好')
      return
    }

    const model = editor.getModel()
    if (model) {
      const value = model.getValue()
      const result = escapeString(value)
      model.setValue(result)
      error.value = ''
    }
  } catch (e) {
    console.error('转义错误:', e)
    error.value = e instanceof Error ? e.message : '未知错误'
    ElMessage.error(`转义失败: ${error.value}`)
  }
}

const unescapeJson = async () => {
  try {
    const editor = monacoEditor.value?.editor
    if (!editor) {
      ElMessage.error('编辑器未准备好')
      return
    }

    const model = editor.getModel()
    if (model) {
      const value = model.getValue()
      const result = unescapeString(value)
      model.setValue(result)
      error.value = ''
    }
  } catch (e) {
    console.error('去转义错误:', e)
    error.value = e instanceof Error ? e.message : '未知错误'
    ElMessage.error(`去转义失败: ${error.value}`)
  }
}

// 修改 processString 函数，只处理 Unicode 解码
const processString = (str: string): string => {
  if (settings.value.autoDecodeUnicode) {
    return str.replace(/\\u[\dA-F]{4}/gi, (match) =>
      String.fromCharCode(parseInt(match.replace(/\\u/g, ''), 16))
    )
  }
  return str
}

// 修改 processJsonStrings 函数的条件判断
const processJsonStrings = (value: any, isEscape?: boolean): any => {
  if (typeof value === 'string') {
    if (isEscape === undefined) {
      return processString(value) // 用于格式化和压缩
    }
    return isEscape ? escapeString(value) : unescapeString(value) // 用于转义和去转义
  }
  if (Array.isArray(value)) {
    return value.map((item) => processJsonStrings(item))
  }
  if (typeof value === 'object' && value !== null) {
    const result: Record<string, any> = {}
    for (const [key, val] of Object.entries(value)) {
      result[key] = processJsonStrings(val)
    }
    return result
  }
  return value
}

// 修改转义函数，只处理引号和反斜杠
const escapeString = (str: string): string => {
  return str
    .replace(/\\/g, '\\\\') // 先处理反斜杠
    .replace(/"/g, '\\"') // 再处理双引号
}

// 修改去转义函数，根据设置处理换行符和制表符
const unescapeString = (str: string): string => {
  let result = str
    .replace(/\\"/g, '"') // 先处理双引号
    .replace(/\\\\/g, '\\') // 再处理反斜杠

  if (settings.value.removeEscapes) {
    result = result
      .replace(/\\n/g, '\n') // 换行符
      .replace(/\\t/g, '\t') // 制表符
      .replace(/\\r/g, '\r') // 回车符
  }

  return result
}

const handleChange = (value: string) => {
  code.value = value
  validateJson()
}

const validateJson = () => {
  if (!code.value.trim()) {
    error.value = ''
    return
  }
  try {
    JSON.parse(code.value)
    error.value = ''
  } catch (e) {
    if (e instanceof Error) {
      error.value = e.message
    }
  }
}

const copyToClipboard = async () => {
  try {
    const editor = monacoEditor.value?.editor
    if (!editor) {
      ElMessage.error('编辑器未准备好')
      return
    }

    const model = editor.getModel()
    if (!model) {
      ElMessage.error('获取内容失败')
      return
    }

    const content = model.getValue()
    if (!content) {
      ElMessage.error('内容为空')
      return
    }

    await copy(content)
  } catch (e) {
    console.error('复制失败:', e)
    ElMessage.error('复制失败')
  }
}

const clearContent = () => {
  try {
    const editor = monacoEditor.value?.editor
    if (!editor) {
      ElMessage.error('编辑器未准备好')
      return
    }

    const model = editor.getModel()
    if (!model) {
      ElMessage.error('获取内容失败')
      return
    }

    if (!model.getValue()) {
      ElMessage.error('内容已经为空')
      return
    }

    model.setValue('')
    error.value = ''
  } catch (e) {
    console.error('清空失败:', e)
    ElMessage.error('清空失败')
  }
}

// 加载示例数据
const loadSample = () => {
  try {
    const editor = monacoEditor.value?.editor
    if (!editor) {
      ElMessage.error('编辑器未准备好')
      return
    }

    const model = editor.getModel()
    if (model) {
      // 使用原始字符串模板来保持转义序列
      const sampleStr = `{
  "user": {
    "id": 123,
    "name": "John Doe",
    "email": "john.doe@example.com"
  },
  "products": [
    {
      "id": "p1",
      "name": "Product A",
      "price": 19.99
    },
    {
      "id": "p2",
      "name": "Product B",
      "price": 29.99
    },
    {
      "id": "p3",
      "name": "Product \\u4E2D\\u6587",
      "price": 39.99
    }
  ],
  "order": {
    "orderId": "abc123",
    "date": "2023-08-18",
    "items": [
      {
        "productId": "p1",
        "quantity": 2
      },
      {
        "productId": "p3",
        "quantity": 1
      }
    ]
  }
}`
      code.value = sampleStr // 直接更新 store
      model.setValue(sampleStr)
    }
  } catch (e) {
    console.error('加载示例失败:', e)
    ElMessage.error('加载示例失败')
  }
}

const closeSettings = () => {
  showSettings.value = false
}
</script>

<style scoped>
.json-editor {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.toolbar {
  flex: 0 0 48px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 16px;
  /* background: #f8f9fa; */
  border-bottom: 1px solid #eaecef;
  position: relative;
}

.tools-group {
  display: flex;
  gap: 12px;
}

.tool-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 32px;
  padding: 0 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  background: #fff;
  color: #24292f;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.tool-btn:hover {
  background: #f6f8fa;
  border-color: #bbc0c4;
}

.config-wrapper {
  position: relative;
}

.config-btn {
  background: #f3f4f6;
}

.editor-wrapper {
  flex: 1;
  position: relative;
}

.editor-container {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

:deep(.monaco-editor) {
  height: 100% !important;
}

:deep(.monaco-scrollable-element) {
  height: 100% !important;
}

/* 确保滚动条可见和可用 */
:deep(.monaco-editor .scrollbar) {
  width: 14px !important;
  right: 0 !important;
}

.settings-panel {
  position: absolute;
  top: 100%;
  left: 0;
  margin-top: 4px;
  background: #fff;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  padding: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  min-width: 200px;
}

.setting-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 0;
  font-size: 13px;
  color: #24292f;
  white-space: nowrap;
}

.setting-item input[type='checkbox'] {
  margin: 0;
  width: 14px;
  height: 14px;
}
</style>
