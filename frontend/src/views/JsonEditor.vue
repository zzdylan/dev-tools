<template>
  <div class="json-editor">
    <div class="tabs-header" ref="tabsHeader" @wheel="handleTabsWheel">
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
        <!-- 模式切换按钮 -->
        <button class="tool-btn mode-btn" @click="toggleMode">
          <span class="tool-icon">{{ currentTab?.compareMode ? '📝' : '🔀' }}</span>
          {{ currentTab?.compareMode ? '普通模式' : '对比模式' }}
        </button>

        <button class="tool-btn config-btn" ref="configBtn" @click="toggleSettings">
          <span class="tool-icon">⚙️</span>
          配置
        </button>

        <!-- 配置面板 -->
        <div v-show="showSettings" class="settings-panel" ref="settingsPanel">
          <label class="setting-item">
            <input type="checkbox" v-model="settings.autoDecodeUnicode" />
            自动解码 Unicode
          </label>
          <label class="setting-item">
            <input type="checkbox" v-model="settings.showArrayIndex" @change="handleArrayIndexSettingChange" />
            显示数组索引
          </label>
        </div>
      </div>

      <!-- 功能按钮放右边 -->
      <div class="tools-group">
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
        <button class="tool-btn" @click="loadSample">
          <span class="tool-icon">📝</span>
          示例
        </button>
        <button class="tool-btn" @click="removeAllTabs">
          <span class="tool-icon">🗂️</span>
          删标签页
        </button>
        <button class="tool-btn" @click="clearContent">
          <span class="tool-icon">🗑️</span>
          清空
        </button>
      </div>
    </div>

    <div class="editor-wrapper">
      <div v-for="id in Object.keys(jsonEditorTabs)" :key="id" class="editor-container" v-show="id === tabId">
        <!-- 编辑模式：单编辑器 -->
        <template v-if="!jsonEditorTabs[id].compareMode">
          <MonacoEditor :ref="(el: any) => { if (el) handleEditorMount(el, id) }" :value="jsonEditorTabs[id].code"
            @change="(val: string) => handleChange(val, id)" :options="options" language="json" theme="vs" />
        </template>

        <!-- 对比模式：Diff Editor -->
        <template v-else>
          <div :ref="(el: any) => { if (el) diffEditorRefs[id] = el }" class="diff-editor-container"></div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, reactive, watch, onMounted, toRaw, onBeforeUnmount } from 'vue'
import MonacoEditor from 'monaco-editor-vue3'
import { ElMessage, ElMessageBox } from 'element-plus'
import { FormatJson, CompressJson } from '../../wailsjs/go/processor/JsonProcessor'
import { onClickOutside } from '@vueuse/core'
import { useToolsStore } from '../stores/tools'
import { storeToRefs } from 'pinia'
import { useRoute, useRouter } from 'vue-router'
import * as monaco from 'monaco-editor'

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
})

// 组件挂载后初始化当前标签页的对比模式（如果需要）
onMounted(() => {
  // 初始化所有处于对比模式的标签页
  Object.keys(store.jsonEditorTabs).forEach(id => {
    const tab = store.jsonEditorTabs[id]
    if (tab?.compareMode && !diffEditorInstances[id]) {
      // 只有当前标签页才需要立即创建
      if (id === tabId.value) {
        nextTick(() => {
          createDiffEditor(id)
        })
      }
    }
  })
})

// 组件卸载时清理定时器
onBeforeUnmount(() => {
  Object.values(decorationTimers).forEach(timer => clearTimeout(timer))
})

// 为每个标签页保存编辑器引用
const editorRefs = reactive<Record<string, any>>({})
// Diff 编辑器容器引用
const diffEditorRefs = reactive<Record<string, HTMLElement>>({})
// Diff 编辑器实例引用
const diffEditorInstances = reactive<Record<string, monaco.editor.IStandaloneDiffEditor>>({})
// 存储装饰器集合
const editorDecorations = reactive<Record<string, string[]>>({})
// 存储去抖定时器
const decorationTimers = reactive<Record<string, NodeJS.Timeout>>({})

// 使用当前标签页的数据
const currentTab = computed(() => store.jsonEditorTabs[tabId.value])

// 当前标签页的编辑器
const getCurrentEditor = () => {
  // 如果是对比模式，返回 Diff Editor 的修改编辑器
  if (currentTab.value?.compareMode) {
    const diffEditor = toRaw(diffEditorInstances[tabId.value])
    if (diffEditor) {
      return {
        editor: diffEditor.getModifiedEditor()
      }
    }
    return null
  }
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
  get: () => store.jsonEditorSettings,
  set: (val) => {
    store.jsonEditorSettings = val
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
  glyphMargin: true,  // 启用左侧图标栏
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
  renderLineHighlight: 'none',
}

// 配置状态
const showSettings = ref(false)
const settingsPanel = ref<HTMLElement | null>(null)
const configBtn = ref<HTMLElement | null>(null)
const tabsHeader = ref<HTMLElement | null>(null)

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
    // 对比模式：格式化两个编辑器
    if (currentTab.value?.compareMode) {
      const diffEditor = toRaw(diffEditorInstances[tabId.value])
      if (!diffEditor) {
        ElMessage.error('编辑器未准备好')
        return
      }

      const model = diffEditor.getModel()
      if (!model) {
        ElMessage.error('获取内容失败')
        return
      }

      // 格式化原始编辑器
      const originalValue = toRaw(model.original).getValue()
      if (originalValue.trim()) {
        const formattedOriginal = await FormatJson(originalValue, settings.value.autoDecodeUnicode)
        toRaw(model.original).setValue(formattedOriginal)
      }

      // 格式化修改编辑器
      const modifiedValue = toRaw(model.modified).getValue()
      if (modifiedValue.trim()) {
        const formattedModified = await FormatJson(modifiedValue, settings.value.autoDecodeUnicode)
        toRaw(model.modified).setValue(formattedModified)
      }

      return
    }

    // 编辑模式：格式化单个编辑器
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
    // 对比模式：压缩两个编辑器
    if (currentTab.value?.compareMode) {
      const diffEditor = toRaw(diffEditorInstances[tabId.value])
      if (!diffEditor) {
        ElMessage.error('编辑器未准备好')
        return
      }

      const model = diffEditor.getModel()
      if (!model) {
        ElMessage.error('获取内容失败')
        return
      }

      // 压缩原始编辑器
      const originalValue = toRaw(model.original).getValue()
      if (originalValue.trim()) {
        const compressedOriginal = await CompressJson(originalValue, settings.value.autoDecodeUnicode)
        toRaw(model.original).setValue(compressedOriginal)
      }

      // 压缩修改编辑器
      const modifiedValue = toRaw(model.modified).getValue()
      if (modifiedValue.trim()) {
        const compressedModified = await CompressJson(modifiedValue, settings.value.autoDecodeUnicode)
        toRaw(model.modified).setValue(compressedModified)
      }

      return
    }

    // 编辑模式：压缩单个编辑器
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

    // 更新数组索引装饰器（去抖处理）
    updateArrayDecorations(id)
  }
}

// 解析JSON并找到所有数组元素的位置
interface ArrayInfo {
  line: number
  index: number
  total: number
}

const findArrayElements = (text: string): ArrayInfo[] => {
  const results: ArrayInfo[] = []

  try {
    const lines = text.split('\n')
    let inString = false
    let escapeNext = false
    let bracketDepth = 0  // [] 深度
    let braceDepth = 0    // {} 深度
    let currentArrayTotal = 0
    let currentArrayIndex = 0
    let elementFoundOnLine = false // 标记当前行是否已找到元素

    for (let i = 0; i < lines.length; i++) {
      const line = lines[i]
      elementFoundOnLine = false

      for (let j = 0; j < line.length; j++) {
        const char = line[j]

        // 处理转义字符
        if (escapeNext) {
          escapeNext = false
          continue
        }
        if (char === '\\') {
          escapeNext = true
          continue
        }

        // 处理字符串状态
        if (char === '"') {
          inString = !inString
          continue
        }

        // 只在非字符串内处理结构字符
        if (!inString) {
          if (char === '{') {
            // 只在数组第一层且不在对象内部时，这才是数组元素的开始
            if (bracketDepth >= 1 && braceDepth === 0 && !elementFoundOnLine) {
              results.push({
                line: i + 1,  // Monaco使用1-based行号
                index: currentArrayIndex,
                total: currentArrayTotal
              })
              currentArrayIndex++
              elementFoundOnLine = true // 防止同一行多次添加
            }
            braceDepth++
          } else if (char === '}') {
            braceDepth--
          } else if (char === '[') {
            bracketDepth++
            if (bracketDepth >= 1) {
              // 进入新数组，计算元素个数并重置braceDepth
              currentArrayTotal = countArrayElements(lines, i)
              currentArrayIndex = 0
              braceDepth = 0  // 重置braceDepth，因为我们只关心数组内部的对象
            }
          } else if (char === ']') {
            bracketDepth--
          }
        }
      }
    }
  } catch (e) {
    console.error('Error finding array elements:', e)
  }

  return results
}

// 计算数组元素个数 - 只计算顶层元素,不包括嵌套对象的字段
const countArrayElements = (lines: string[], startLine: number): number => {
  let bracketDepth = 0
  let braceDepth = 0
  let count = 0
  let inString = false
  let escapeNext = false

  for (let i = startLine; i < lines.length; i++) {
    const line = lines[i]

    for (let j = 0; j < line.length; j++) {
      const char = line[j]

      // 处理转义
      if (escapeNext) {
        escapeNext = false
        continue
      }
      if (char === '\\') {
        escapeNext = true
        continue
      }

      // 处理字符串
      if (char === '"') {
        inString = !inString
        continue
      }

      if (!inString) {
        if (char === '[') {
          bracketDepth++
        } else if (char === ']') {
          bracketDepth--
          if (bracketDepth === 0) {
            return count
          }
        } else if (char === '{') {
          // 只在数组第一层且不在对象内部时才是数组元素
          if (bracketDepth === 1 && braceDepth === 0) {
            count++
          }
          braceDepth++
        } else if (char === '}') {
          braceDepth--
        }
      }
    }
  }

  return count
}


// 更新数组装饰器（带去抖）
const updateArrayDecorations = (id: string) => {
  // 清除之前的定时器
  if (decorationTimers[id]) {
    clearTimeout(decorationTimers[id])
  }

  // 设置新的定时器（300ms 去抖）
  decorationTimers[id] = setTimeout(() => {
    applyArrayDecorations(id)
  }, 300)
}

// 动态注入CSS样式
const injectArrayIndexStyles = (arrayInfos: ArrayInfo[]) => {
  // 移除旧的样式
  const oldStyle = document.getElementById('array-index-dynamic-styles')
  if (oldStyle) {
    oldStyle.remove()
  }

  if (arrayInfos.length === 0) return

  // 创建新的样式
  const style = document.createElement('style')
  style.id = 'array-index-dynamic-styles'

  let css = ''
  arrayInfos.forEach(info => {
    const className = `array-index-marker-${info.line}`
    css += `
      .monaco-editor .${className}::after {
        content: ' [${info.index + 1}/${info.total}]';
        color: #9ca3af;
        font-size: 10px;
        font-weight: 600;
        opacity: 0.75;
        margin-left: 2px;
        font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, monospace;
        user-select: none;
        pointer-events: none;
      }
    `
  })

  style.textContent = css
  document.head.appendChild(style)
}

// 应用数组装饰器 - 使用 inlineClassName
const applyArrayDecorations = (id: string) => {
  const editor = editorRefs[id]?.editor
  if (!editor) return

  const model = editor.getModel()
  if (!model) return

  // 移除旧的装饰器
  if (editorDecorations[id]) {
    editorDecorations[id] = editor.deltaDecorations(editorDecorations[id], [])
  }

  // 如果配置为不显示数组索引，移除样式并返回
  if (!settings.value.showArrayIndex) {
    const oldStyle = document.getElementById('array-index-dynamic-styles')
    if (oldStyle) {
      oldStyle.remove()
    }
    return
  }

  const text = model.getValue()
  const arrayInfos = findArrayElements(text)

  // 注入动态样式
  injectArrayIndexStyles(arrayInfos)

  // 创建装饰器配置 - 给 { 字符添加class
  const decorations = arrayInfos.map(info => {
    const lineContent = model.getLineContent(info.line)
    const braceIndex = lineContent.indexOf('{')
    const startCol = braceIndex >= 0 ? braceIndex + 1 : 1
    const endCol = startCol + 1

    return {
      range: new monaco.Range(info.line, startCol, info.line, endCol),
      options: {
        inlineClassName: `array-index-marker-${info.line}`
      }
    }
  })

  // 应用装饰器
  editorDecorations[id] = editor.deltaDecorations([], decorations)
}

// 处理数组索引显示设置变化
const handleArrayIndexSettingChange = () => {
  // 为所有标签页重新应用装饰器（根据新的设置）
  Object.keys(store.jsonEditorTabs).forEach(id => {
    if (editorRefs[id]?.editor) {
      applyArrayDecorations(id)
    }
  })
}

// 处理编辑器挂载
const handleEditorMount = (el: any, id: string) => {
  if (el) {
    editorRefs[id] = el
    // 编辑器挂载后，如果有内容则应用装饰器
    nextTick(() => {
      if (store.jsonEditorTabs[id]?.code) {
        applyArrayDecorations(id)
      }
    })
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


const removeAllTabs = async () => {
  try {
    // 获取所有标签页ID
    const allTabIds = Object.keys(store.jsonEditorTabs)
    const nonDefaultTabs = allTabIds.filter(id => id !== 'default')
    
    // 构建确认消息
    let message = '确定要删除所有标签页并清空内容吗？'
    if (nonDefaultTabs.length > 0) {
      message += `\n\n将删除 ${nonDefaultTabs.length} 个标签页，并清空默认标签页的内容。`
    } else {
      message += '\n\n将清空默认标签页的内容。'
    }
    
    // 显示确认对话框
    await ElMessageBox.confirm(
      message,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
        dangerouslyUseHTMLString: false
      }
    )
    
    // 删除所有非默认标签页
    allTabIds.forEach(id => {
      if (id !== 'default') {
        delete store.jsonEditorTabs[id]
        delete editorRefs[id]
      }
    })
    
    // 清空默认标签页内容
    if (store.jsonEditorTabs.default) {
      store.jsonEditorTabs.default.code = ''
    }
    
    // 清空历史记录
    tabHistory.value = ['default']
    
    // 跳转到默认标签页
    router.push({ name: 'JsonEditorTab', params: { id: 'default' } })
    
    ElMessage.success('已删除所有标签页并清空内容')
  } catch (err: any) {
    // 用户取消操作时不显示错误
    if (err === 'cancel') {
      return
    }
    ElMessage.error('删除标签页失败：' + (err.message || err))
  }
}

const clearContent = () => {
  try {
    // 对比模式：清空两个编辑器
    if (currentTab.value?.compareMode) {
      const diffEditor = toRaw(diffEditorInstances[tabId.value])
      if (!diffEditor) {
        ElMessage.error('编辑器未准备好')
        return
      }

      const model = diffEditor.getModel()
      if (!model) {
        ElMessage.error('获取内容失败')
        return
      }

      toRaw(model.original).setValue('')
      toRaw(model.modified).setValue('')
      error.value = ''
      ElMessage.success('已清空')
      return
    }

    // 编辑模式：清空单个编辑器
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

    const sampleStr2 = `{
  "user": {
    "id": 456,
    "name": "Jane Smith",
    "email": "jane.smith@example.com"
  },
  "products": [
    {
      "id": "p1",
      "name": "Product A",
      "price": 19.99
    },
    {
      "id": "p4",
      "name": "Product D",
      "price": 49.99
    }
  ],
  "order": {
    "orderId": "xyz789",
    "date": "2023-08-20",
    "items": [
      {
        "productId": "p1",
        "quantity": 1
      },
      {
        "productId": "p4",
        "quantity": 3
      }
    ]
  }
}`

    // 对比模式：加载两个不同的示例
    if (currentTab.value?.compareMode) {
      const diffEditor = toRaw(diffEditorInstances[tabId.value])
      if (!diffEditor) {
        ElMessage.error('编辑器未准备好')
        return
      }

      const model = diffEditor.getModel()
      if (!model) {
        ElMessage.error('获取内容失败')
        return
      }

      toRaw(model.original).setValue(sampleStr)
      toRaw(model.modified).setValue(sampleStr2)
      return
    }

    // 编辑模式：加载单个示例
    const currentEditor = getCurrentEditor()
    if (!currentEditor?.editor) {
      ElMessage.error('编辑器未准备好')
      return
    }

    const model = currentEditor.editor.getModel()
    if (model) {
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

// 处理标签页区域的滚轮事件
const handleTabsWheel = (event: WheelEvent) => {
  if (tabsHeader.value) {
    event.preventDefault()
    // 水平滚动，deltaY是垂直滚动量，我们将其转换为水平滚动
    tabsHeader.value.scrollLeft += event.deltaY
  }
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
    // 释放编辑器引用（不调用dispose，避免卡死）
    delete editorRefs[id]
    delete editorRefs[id + '_compare']

    // 删除 Diff Editor 引用（不调用dispose，避免卡死）
    delete diffEditorInstances[id]
    delete diffEditorRefs[id]

    // 从存储中删除标签页
    delete store.jsonEditorTabs[id]

    // 清理历史记录中已删除的标签页引用
    tabHistory.value = tabHistory.value.filter(tabId =>
      tabId !== id && store.jsonEditorTabs[tabId]
    )

    console.log('Tab removed:', id)
  })
}

// 切换编辑/对比模式
const toggleMode = () => {
  if (currentTab.value) {
    currentTab.value.compareMode = !currentTab.value.compareMode

    // 切换到对比模式时，创建 Diff Editor
    if (currentTab.value.compareMode) {
      nextTick(() => {
        createDiffEditor(tabId.value)
      })
    } else {
      // 切换回编辑模式时，只删除引用，不调用 dispose()（避免卡死）
      // DOM 会被 v-if 自动移除
      delete diffEditorInstances[tabId.value]
      delete diffEditorRefs[tabId.value]
    }
  }
}

// 创建 Diff Editor
const createDiffEditor = (id: string) => {
  console.log('createDiffEditor called for tab:', id)
  const container = diffEditorRefs[id]
  console.log('Container element:', container)

  if (!container) {
    console.error('Container not found for tab:', id)
    return
  }

  // 如果已存在实例，直接返回
  if (diffEditorInstances[id]) {
    console.log('Diff editor already exists for tab:', id)
    return
  }

  try {
    console.log('Creating diff editor instance...')
    // 创建 Diff Editor
    const diffEditor = monaco.editor.createDiffEditor(container, {
      fontSize: 12,
      automaticLayout: true,
      renderSideBySide: true, // 并排显示
      enableSplitViewResizing: true, // 允许调整大小
      readOnly: false,
      minimap: { enabled: false },
      scrollBeyondLastLine: false,
      wordWrap: 'on',
      lineNumbers: 'on',
      renderIndicators: true, // 显示差异指示器
      ignoreTrimWhitespace: false, // 不忽略空白差异
      renderOverviewRuler: false, // 隐藏概览标尺（右侧的颜色条）
      scrollbar: {
        vertical: 'visible',
        horizontal: 'visible',
        verticalScrollbarSize: 8,
        horizontalScrollbarSize: 8,
        alwaysConsumeMouseWheel: true,
      },
      originalEditable: true, // 允许编辑原始编辑器
    })

    // 创建 model
    const originalModel = monaco.editor.createModel(
      store.jsonEditorTabs[id]?.code || '',
      'json'
    )
    const modifiedModel = monaco.editor.createModel(
      store.jsonEditorTabs[id]?.compareCode || '',
      'json'
    )

    // 设置 model
    diffEditor.setModel({
      original: originalModel,
      modified: modifiedModel
    })

    // 监听修改编辑器的内容变化
    modifiedModel.onDidChangeContent(() => {
      if (store.jsonEditorTabs[id]) {
        store.jsonEditorTabs[id].compareCode = modifiedModel.getValue()
      }
    })

    // 监听原始编辑器的内容变化
    originalModel.onDidChangeContent(() => {
      if (store.jsonEditorTabs[id]) {
        store.jsonEditorTabs[id].code = originalModel.getValue()
      }
    })

    diffEditorInstances[id] = diffEditor
    console.log('Diff editor created successfully for tab:', id)
  } catch (e) {
    console.error('Error creating diff editor:', e)
  }
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
  /* 启用滚轮滚动 */
  scroll-behavior: smooth;
}

.tabs-nav {
  display: flex;
  align-items: stretch;
  min-width: max-content;
  background: #f8f9fa;
  flex-shrink: 0;
}

.tab-item {
  padding: 0 8px;
  margin: 0;
  background: #f8f9fa;
  border: none;
  color: #6c757d;
  font-size: 10px;
  cursor: pointer;
  width: 60px;
  min-width: 60px;
  max-width: none;
  height: 100%;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 2px;
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
  padding: 1px 3px;
  border-radius: 2px;
  font-size: 9px;
  margin-left: 2px;
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

/* 自定义滚动条样式 */
.tabs-header::-webkit-scrollbar {
  height: 4px;
}

.tabs-header::-webkit-scrollbar-track {
  background: #f8f9fa;
}

.tabs-header::-webkit-scrollbar-thumb {
  background: #d1d5db;
  border-radius: 2px;
}

.tabs-header::-webkit-scrollbar-thumb:hover {
  background: #9ca3af;
}

/* Diff Editor 容器样式 */
.diff-editor-container {
  width: 100%;
  height: 100%;
  position: relative;
  overflow: hidden;
}

.diff-editor-container :deep(.monaco-diff-editor) {
  height: 100% !important;
}

.diff-editor-container :deep(.monaco-editor) {
  height: 100% !important;
}

/* 隐藏 Diff Editor 的 overview ruler（右侧颜色条）*/
.diff-editor-container :deep(.decorationsOverviewRuler) {
  display: none !important;
}

/* 自定义 Diff Editor 滚动条样式 */
.diff-editor-container :deep(.monaco-scrollable-element > .scrollbar > .slider) {
  background: rgba(100, 100, 100, 0.4) !important;
}

.diff-editor-container :deep(.monaco-scrollable-element > .scrollbar.vertical) {
  width: 8px !important;
}

.diff-editor-container :deep(.monaco-scrollable-element > .scrollbar.horizontal) {
  height: 8px !important;
}

.diff-editor-container :deep(.monaco-scrollable-element > .scrollbar > .slider:hover) {
  background: rgba(100, 100, 100, 0.7) !important;
}

/* 激活状态的按钮样式 */
.tool-btn.active {
  background: #e0f2fe;
  border-color: #0ea5e9;
  color: #0369a1;
}

.tool-btn.active:hover {
  background: #bae6fd;
  border-color: #0284c7;
}

/* 数组索引装饰器样式 - 在代码行内显示 */
:deep(.array-index-inline) {
  color: #9ca3af !important;
  font-size: 10px !important;
  font-weight: 600 !important;
  opacity: 0.75 !important;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, monospace !important;
  font-style: normal !important;
  user-select: none !important;
  pointer-events: none !important;
}
</style>
