<template>
  <div class="json-editor">
    <div class="tabs-header">
      <div class="tabs-nav">
        <div v-for="id in Object.keys(jsonEditorTabs)" :key="id" class="tab-item" :class="{ active: id === tabId }"
          @click="handleTabChange(id)">
          {{ id === 'default' ? '标签1' : `标签${Object.keys(jsonEditorTabs).indexOf(id) + 1}` }}
          <span v-if="id !== 'default'" class="close-btn" @click.stop="closeTab(id)">
            ×
          </span>
        </div>
        <div class="add-tab" @click="createTab">+</div>
      </div>
    </div>

    <div class="toolbar">
      <div class="config-wrapper">
        <button class="tool-btn config-btn" ref="configBtn" @click="toggleSettings">
          <span class="tool-icon">⚙️</span>
          配置
        </button>

        <!-- 配置面板 -->
        <div v-show="showSettings" class="settings-panel" ref="settingsPanel" v-click-outside="closeSettings">
          <label class="setting-item">
            <input type="checkbox" v-model="settings.autoDecodeUnicode" />
            自动解码 Unicode
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
      <div v-for="id in Object.keys(jsonEditorTabs)" :key="id" class="editor-container" v-show="id === tabId">
        <MonacoEditor :ref="(el: any) => { if (el) editorRefs[id] = el }" :value="jsonEditorTabs[id].code"
          @change="(val: string) => handleChange(val, id)" :options="options" language="json" theme="vs" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, reactive, watch } from 'vue'
import MonacoEditor from 'monaco-editor-vue3'
import { useClipboard } from '@vueuse/core'
import { ElMessage } from 'element-plus'
import { FormatJson, CompressJson } from '../../wailsjs/go/main/JsonProcessor'
import { onClickOutside } from '@vueuse/core'
import { useToolsStore } from '../stores/tools'
import { storeToRefs } from 'pinia'
import { useRoute, useRouter } from 'vue-router'

const { copy } = useClipboard()
const store = useToolsStore()
const { jsonEditorTabs } = storeToRefs(store)
const route = useRoute()
const router = useRouter()
const tabId = computed(() => route.params.id as string)

// 活动标签页名称，用于 ElementPlus Tabs 组件
const activeTabName = ref(tabId.value)

// 标签页访问历史，用于删除标签页时回到上一个访问的标签页
const tabHistory = ref<string[]>([tabId.value])

// 监听路由变化，更新活动标签页和历史记录
watch(tabId, (newTabId, oldTabId) => {
  activeTabName.value = newTabId
  
  // 更新store中的当前标签页
  store.setCurrentJsonEditorTab(newTabId)

  // 更新访问历史
  if (oldTabId && oldTabId !== newTabId) {
    // 移除历史中的当前标签页（如果存在）
    const index = tabHistory.value.indexOf(newTabId)
    if (index > -1) {
      tabHistory.value.splice(index, 1)
    }
    // 将新标签页添加到历史记录的开头
    tabHistory.value.unshift(newTabId)

    // 限制历史记录长度，保持最近的10个
    if (tabHistory.value.length > 10) {
      tabHistory.value = tabHistory.value.slice(0, 10)
    }
  }
}, { immediate: true })

// 为每个标签页保存编辑器引用
const editorRefs = reactive<Record<string, any>>({})

// 使用当前标签页的数据
const currentTab = computed(() => store.jsonEditorTabs[tabId.value])

// 当前标签页的编辑器
const getCurrentEditor = () => {
  return editorRefs[tabId.value]
}

const code = computed({
  get: () => currentTab.value?.code ?? '',
  set: (val) => {
    if (currentTab.value) {
      currentTab.value.code = val
    }
  },
})

const settings = computed({
  get: () =>
    currentTab.value?.settings ?? {
      autoDecodeUnicode: false,
      removeEscapes: false,
    },
  set: (val) => {
    if (currentTab.value) {
      currentTab.value.settings = val
    }
  },
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
    verticalScrollbarSize: 8,
    horizontalScrollbarSize: 8,
    alwaysConsumeMouseWheel: true,
  },
  lineDecorationsWidth: 0,
  lineNumbersMinChars: 0,
  glyphMargin: false,
  renderLineHighlight: 'none',
}

// 配置状态
const showSettings = ref(false)
const settingsPanel = ref<HTMLElement | null>(null)
const configBtn = ref<HTMLElement | null>(null)

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
    const currentEditor = getCurrentEditor()
    if (!currentEditor?.editor) {
      ElMessage.error('编辑器未准备好')
      return
    }

    const model = currentEditor.editor.getModel()
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

    // 创建撤销停止点，使格式化成为独立的撤销单元
    currentEditor.editor.pushUndoStop()

    const fullRange = model.getFullModelRange()
    currentEditor.editor.executeEdits('formatJson', [{
      range: fullRange,
      text: formatted
    }])

    currentEditor.editor.pushUndoStop()
  } catch (err: any) {
    ElMessage.error('格式化失败：' + (err.message || err))
  }
}

const compressJson = async () => {
  try {
    const currentEditor = getCurrentEditor()
    if (!currentEditor?.editor) {
      ElMessage.error('编辑器未准备好')
      return
    }

    const model = currentEditor.editor.getModel()
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

    // 创建撤销停止点，使压缩成为独立的撤销单元
    currentEditor.editor.pushUndoStop()

    const fullRange = model.getFullModelRange()
    currentEditor.editor.executeEdits('compressJson', [{
      range: fullRange,
      text: compressed
    }])

    currentEditor.editor.pushUndoStop()
  } catch (err: any) {
    ElMessage.error('压缩失败：' + (err.message || err))
  }
}

// 添加新的转义和去转义函数
const escapeJson = async () => {
  try {
    const currentEditor = getCurrentEditor()
    if (!currentEditor?.editor) {
      ElMessage.error('编辑器未准备好')
      return
    }

    const model = currentEditor.editor.getModel()
    if (model) {
      const value = model.getValue()
      const result = escapeString(value)

      // 使用 executeEdits 来保持撤销栈
      const fullRange = model.getFullModelRange()
      currentEditor.editor.executeEdits('escapeJson', [{
        range: fullRange,
        text: result
      }])
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
    const currentEditor = getCurrentEditor()
    if (!currentEditor?.editor) {
      ElMessage.error('编辑器未准备好')
      return
    }

    const model = currentEditor.editor.getModel()
    if (model) {
      const value = model.getValue()
      const result = unescapeString(value)

      // 使用 executeEdits 来保持撤销栈
      const fullRange = model.getFullModelRange()
      currentEditor.editor.executeEdits('unescapeJson', [{
        range: fullRange,
        text: result
      }])
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

// 使用原生 JSON 方法进行转义和去转义 - 最可靠的方案
const escapeString = (str: string): string => {
  // 使用 JSON.stringify 然后去掉首尾引号
  return JSON.stringify(str).slice(1, -1)
}

// 标准 JSON 去转义函数
const unescapeString = (str: string): string => {
  try {
    // 使用 JSON.parse 进行去转义，需要加上引号
    return JSON.parse('"' + str + '"')
  } catch (error) {
    // 如果解析失败，返回原字符串
    console.warn('去转义失败，返回原字符串:', error)
    return str
  }
}

const handleChange = (value: string, id: string) => {
  // 避免重复更新导致清除撤销栈
  if (store.jsonEditorTabs[id] && store.jsonEditorTabs[id].code !== value) {
    store.jsonEditorTabs[id].code = value

    // 验证JSON但不阻断编辑器操作
    if (id === tabId.value) {
      // 异步验证，避免阻塞编辑器操作
      nextTick(() => {
        validateJson(value)
      })
    }
  }
}

const validateJson = (content: string = '') => {
  const valueToValidate = content || code.value

  if (!valueToValidate.trim()) {
    error.value = ''
    return
  }

  try {
    JSON.parse(valueToValidate)
    error.value = ''
  } catch (e) {
    if (e instanceof Error) {
      error.value = e.message
    }
  }
}

const copyToClipboard = async () => {
  try {
    const currentEditor = getCurrentEditor()
    if (!currentEditor?.editor) {
      ElMessage.error('编辑器未准备好')
      return
    }

    const model = currentEditor.editor.getModel()
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
    ElMessage.success('复制成功')
  } catch (e) {
    console.error('复制失败:', e)
    ElMessage.error('复制失败')
  }
}

const clearContent = () => {
  try {
    const currentEditor = getCurrentEditor()
    if (!currentEditor?.editor) {
      ElMessage.error('编辑器未准备好')
      return
    }

    const model = currentEditor.editor.getModel()
    if (!model) {
      ElMessage.error('获取内容失败')
      return
    }

    if (!model.getValue()) {
      ElMessage.error('内容已经为空')
      return
    }

    // 使用 executeEdits 来保持撤销栈
    const fullRange = model.getFullModelRange()
    currentEditor.editor.executeEdits('clearContent', [{
      range: fullRange,
      text: ''
    }])
    error.value = ''
  } catch (e) {
    console.error('清空失败:', e)
    ElMessage.error('清空失败')
  }
}

// 加载示例数据
const loadSample = () => {
  try {
    const currentEditor = getCurrentEditor()
    if (!currentEditor?.editor) {
      ElMessage.error('编辑器未准备好')
      return
    }

    const model = currentEditor.editor.getModel()
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
      // 先创建撤销停止点，然后执行操作
      // 这样撤销时会直接回到加载示例之前的状态
      currentEditor.editor.pushUndoStop()

      const fullRange = model.getFullModelRange()
      currentEditor.editor.executeEdits('loadSample', [{
        range: fullRange,
        text: sampleStr
      }])

      // 在操作后也创建停止点，使这个操作成为独立的撤销单元
      currentEditor.editor.pushUndoStop()
    }
  } catch (e) {
    console.error('加载示例失败:', e)
    ElMessage.error('加载示例失败')
  }
}

const closeSettings = () => {
  showSettings.value = false
}

// 处理标签页切换
const handleTabChange = (tabName: string) => {
  router.push({ name: 'JsonEditorTab', params: { id: tabName } })
}

const createTab = () => {
  const newId = store.createJsonEditorTab()
  router.push({ name: 'JsonEditorTab', params: { id: newId } })
}

const closeTab = (targetName: string | number) => {
  const id = String(targetName)

  // 如果要关闭的是当前标签页，需要找到下一个要切换的标签页
  if (id === tabId.value) {
    // 从历史记录中移除当前要关闭的标签页
    const historyIndex = tabHistory.value.indexOf(id)
    if (historyIndex > -1) {
      tabHistory.value.splice(historyIndex, 1)
    }

    // 寻找下一个可用的标签页
    let nextTabId = 'default'

    // 首先尝试从历史记录中找到最近访问的可用标签页
    for (const historyTabId of tabHistory.value) {
      if (historyTabId !== id && store.jsonEditorTabs[historyTabId]) {
        nextTabId = historyTabId
        break
      }
    }

    // 如果历史记录中没有可用的标签页，则查找其他可用标签页
    if (nextTabId === 'default' && !store.jsonEditorTabs['default']) {
      const availableTabs = Object.keys(store.jsonEditorTabs).filter(tabId => tabId !== id)
      if (availableTabs.length > 0) {
        nextTabId = availableTabs[0]
      }
    }

    // 切换到下一个标签页
    router.push({ name: 'JsonEditorTab', params: { id: nextTabId } })
  }

  nextTick(() => {
    // 释放编辑器实例
    delete editorRefs[id]
    // 从存储中删除标签页
    delete store.jsonEditorTabs[id]

    // 清理历史记录中已删除的标签页引用
    tabHistory.value = tabHistory.value.filter(tabId =>
      tabId !== id && store.jsonEditorTabs[tabId]
    )
  })
}
</script>

<style scoped>
.json-editor {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-width: 0;
}

.toolbar {
  flex: 0 0 36px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 8px;
  border-bottom: 1px solid #eaecef;
  position: relative;
  min-width: 0;
}

.tools-group {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  min-width: 0;
}

.tool-btn {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  height: 24px;
  padding: 0 6px;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  background: #fff;
  color: #24292f;
  font-size: 11px;
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
  overflow: hidden;
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

/* 隐藏滚动条上的光标位置指示器 */
:deep(.monaco-editor .decorationsOverviewRuler) {
  display: none !important;
}

.settings-panel {
  position: absolute;
  top: 100%;
  left: 0;
  margin-top: 4px;
  background: #fff;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  padding: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  min-width: 180px;
}

.setting-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 2px 0;
  font-size: 12px;
  color: #24292f;
  white-space: nowrap;
}

.setting-item input[type='checkbox'] {
  margin: 0;
  width: 14px;
  height: 14px;
}

.tabs-header {
  flex-shrink: 0;
  background: #f8f9fa;
  border-bottom: 1px solid #dcdcdc;
  overflow-x: auto;
  overflow-y: hidden;
  height: 30px;
  display: flex;
  align-items: stretch;
}

.tabs-nav {
  display: flex;
  align-items: stretch;
  min-width: max-content;
  background: #f8f9fa;
}

.tab-item {
  padding: 0 10px;
  margin: 0;
  background: #f8f9fa;
  border: none;
  color: #6c757d;
  font-size: 10px;
  cursor: pointer;
  min-width: 70px;
  max-width: 110px;
  height: 100%;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 3px;
  white-space: nowrap;
  flex-shrink: 0;
  box-sizing: border-box;
  transition: all 0.2s;
}

.tab-item {
  border-left: 1px solid #d1d5db;
  border-right: 1px solid #d1d5db;
}

.tab-item + .tab-item {
  border-left: none;
}

.add-tab {
  border-left: none;
}

.tab-item:hover:not(.active) {
  background: #e9ecef;
}

.tab-item.active {
  background: #ffffff;
  color: #212529;
  font-weight: 500;
  z-index: 1;
  position: relative;
}

.close-btn {
  padding: 1px 4px;
  border-radius: 2px;
  font-size: 10px;
  margin-left: 4px;
}

.close-btn:hover {
  background: #e0e0e0;
  color: #333333;
}

.add-tab {
  padding: 0 10px;
  margin: 0;
  background: #f8f9fa;
  border: none;
  border-right: 1px solid #d1d5db;
  color: #6c757d;
  font-size: 10px;
  cursor: pointer;
  width: 28px;
  min-width: 28px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-sizing: border-box;
  transition: all 0.2s;
}

.add-tab:hover {
  background: #e9ecef;
}
</style>
