<template>
  <div class="xml-editor">
    <div class="tabs-header" ref="tabsHeader" @wheel="handleTabsWheel">
      <div class="tabs-nav">
        <div v-for="id in Object.keys(xmlEditorTabs)" :key="id" class="tab-item" :class="{ active: id === tabId }"
          @click="handleTabChange(id)">
          {{ id === 'default' ? '标签1' : `标签${Object.keys(xmlEditorTabs).indexOf(id) + 1}` }}
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
        <div v-show="showSettings" class="settings-panel" ref="settingsPanel" v-click-outside="closeSettings">
          <label class="setting-item">
            <input type="checkbox" v-model="settings.autoFormat" />
            自动格式化
          </label>
        </div>
      </div>

      <div class="tools-group">
        <button class="tool-btn" @click="formatXml">
          <span class="tool-icon">✏️</span>
          格式化
        </button>
        <button class="tool-btn" @click="compressXml">
          <span class="tool-icon">📦</span>
          压缩
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
      <div v-for="id in Object.keys(xmlEditorTabs)" :key="id" class="editor-container" v-show="id === tabId">
        <!-- 编辑模式：单编辑器 -->
        <template v-if="!xmlEditorTabs[id].compareMode">
          <MonacoEditor :ref="(el: any) => { if (el) editorRefs[id] = el }" :value="xmlEditorTabs[id].code"
            @change="(val: string) => handleChange(val, id)" :options="options" language="xml" theme="vs" />
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
import { ref, computed, nextTick, reactive, watch, onMounted } from 'vue'
import MonacoEditor from 'monaco-editor-vue3'
import { ElMessage, ElMessageBox } from 'element-plus'
import { FormatXML, CompressXML } from '../../wailsjs/go/main/XmlProcessor'
import { onClickOutside } from '@vueuse/core'
import { useToolsStore } from '../stores/tools'
import { storeToRefs } from 'pinia'
import { useRoute, useRouter } from 'vue-router'
import * as monaco from 'monaco-editor'
const store = useToolsStore()
const { xmlEditorTabs } = storeToRefs(store)
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
  store.setCurrentXmlEditorTab(newTabId)

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
  Object.keys(store.xmlEditorTabs).forEach(id => {
    const tab = store.xmlEditorTabs[id]
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

// 为每个标签页保存编辑器引用
const editorRefs = reactive<Record<string, any>>({})
// Diff 编辑器容器引用
const diffEditorRefs = reactive<Record<string, HTMLElement>>({})
// Diff 编辑器实例引用
const diffEditorInstances = reactive<Record<string, monaco.editor.IStandaloneDiffEditor>>({})

// 使用当前标签页的数据
const currentTab = computed(() => store.xmlEditorTabs[tabId.value])

// 当前标签页的编辑器
const getCurrentEditor = () => {
  // 如果是对比模式，返回 Diff Editor 的修改编辑器
  if (currentTab.value?.compareMode) {
    const diffEditor = diffEditorInstances[tabId.value]
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
  get: () => currentTab.value?.settings ?? { autoFormat: false },
  set: (val) => {
    if (currentTab.value) {
      currentTab.value.settings = val
    }
  },
})

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
  autoIndent: 'advanced',
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

const handleChange = (value: string, id: string) => {
  if (store.xmlEditorTabs[id]) {
    store.xmlEditorTabs[id].code = value
  }
}

const loadSample = () => {
  try {
    const sampleXml1 = `<classroom>
  <course>Introduction to Computer Science</course>
  <instructor>Dr. Smith</instructor>
  <students>
    <student>
      <student_id>001</student_id>
      <name>Emily Johnson</name>
      <age>19</age>
      <gender>Female</gender>
      <grades>
        <grade subject="Math">A</grade>
        <grade subject="Programming">B+</grade>
        <grade subject="English">A-</grade>
      </grades>
    </student>
    <student>
      <student_id>002</student_id>
      <name>Michael Smith</name>
      <age>20</age>
      <gender>Male</gender>
      <grades>
        <grade subject="Math">B</grade>
        <grade subject="Programming">A</grade>
        <grade subject="English">B</grade>
      </grades>
    </student>
  </students>
</classroom>`

    const sampleXml2 = `<classroom>
  <course>Advanced Computer Science</course>
  <instructor>Dr. Johnson</instructor>
  <students>
    <student>
      <student_id>001</student_id>
      <name>Emily Johnson</name>
      <age>20</age>
      <gender>Female</gender>
      <grades>
        <grade subject="Math">A+</grade>
        <grade subject="Programming">A</grade>
        <grade subject="English">A</grade>
      </grades>
    </student>
    <student>
      <student_id>003</student_id>
      <name>Sarah Williams</name>
      <age>19</age>
      <gender>Female</gender>
      <grades>
        <grade subject="Math">A</grade>
        <grade subject="Programming">A-</grade>
        <grade subject="English">B+</grade>
      </grades>
    </student>
  </students>
</classroom>`

    // 对比模式：加载两个不同的示例
    if (currentTab.value?.compareMode) {
      const diffEditor = diffEditorInstances[tabId.value]
      if (!diffEditor) {
        ElMessage.error('编辑器未准备好')
        return
      }

      const model = diffEditor.getModel()
      if (!model) {
        ElMessage.error('获取内容失败')
        return
      }

      model.original.setValue(sampleXml1)
      model.modified.setValue(sampleXml2)
      ElMessage.success('已加载示例数据')
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
      code.value = sampleXml1 // 直接更新 store
      model.setValue(sampleXml1)
    }
  } catch (e) {
    console.error('加载示例失败:', e)
    ElMessage.error('加载示例失败')
  }
}

const formatXml = async () => {
  try {
    // 对比模式：格式化两个编辑器
    if (currentTab.value?.compareMode) {
      const diffEditor = diffEditorInstances[tabId.value]
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
      const originalValue = model.original.getValue()
      if (originalValue.trim()) {
        const formattedOriginal = await FormatXML(originalValue)
        model.original.setValue(formattedOriginal)
      }

      // 格式化修改编辑器
      const modifiedValue = model.modified.getValue()
      if (modifiedValue.trim()) {
        const formattedModified = await FormatXML(modifiedValue)
        model.modified.setValue(formattedModified)
      }

      ElMessage.success('格式化完成')
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
      ElMessage.error('请输入 XML 内容')
      return
    }

    const formatted = await FormatXML(value)
    model.setValue(formatted)
  } catch (err: any) {
    ElMessage.error('格式化失败：' + (err.message || err))
  }
}

const compressXml = async () => {
  try {
    // 对比模式：压缩两个编辑器
    if (currentTab.value?.compareMode) {
      const diffEditor = diffEditorInstances[tabId.value]
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
      const originalValue = model.original.getValue()
      if (originalValue.trim()) {
        const compressedOriginal = await CompressXML(originalValue)
        model.original.setValue(compressedOriginal)
      }

      // 压缩修改编辑器
      const modifiedValue = model.modified.getValue()
      if (modifiedValue.trim()) {
        const compressedModified = await CompressXML(modifiedValue)
        model.modified.setValue(compressedModified)
      }

      ElMessage.success('压缩完成')
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
      ElMessage.error('请输入 XML 内容')
      return
    }

    const compressed = await CompressXML(value)
    model.setValue(compressed)
  } catch (err: any) {
    ElMessage.error('压缩失败：' + (err.message || err))
  }
}

const removeAllTabs = async () => {
  try {
    // 获取所有标签页ID
    const allTabIds = Object.keys(store.xmlEditorTabs)
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
        delete store.xmlEditorTabs[id]
        delete editorRefs[id]
      }
    })
    
    // 清空默认标签页内容
    if (store.xmlEditorTabs.default) {
      store.xmlEditorTabs.default.code = ''
    }
    
    // 清空历史记录
    tabHistory.value = ['default']
    
    // 跳转到默认标签页
    router.push({ name: 'XmlEditorTab', params: { id: 'default' } })
    
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
      const diffEditor = diffEditorInstances[tabId.value]
      if (!diffEditor) {
        ElMessage.error('编辑器未准备好')
        return
      }

      const model = diffEditor.getModel()
      if (!model) {
        ElMessage.error('获取内容失败')
        return
      }

      model.original.setValue('')
      model.modified.setValue('')
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

    model.setValue('')
    code.value = ''
  } catch (e) {
    console.error('清空失败:', e)
    ElMessage.error('清空失败')
  }
}

// 处理标签页切换
const handleTabChange = (tabName: string) => {
  router.push({ name: 'XmlEditorTab', params: { id: tabName } })
}

const createTab = () => {
  const newId = store.createXmlEditorTab()
  router.push({ name: 'XmlEditorTab', params: { id: newId } })
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
      if (historyTabId !== id && store.xmlEditorTabs[historyTabId]) {
        nextTabId = historyTabId
        break
      }
    }

    // 如果历史记录中没有可用的标签页，则查找其他可用标签页
    if (nextTabId === 'default' && !store.xmlEditorTabs['default']) {
      const availableTabs = Object.keys(store.xmlEditorTabs).filter(tabId => tabId !== id)
      if (availableTabs.length > 0) {
        nextTabId = availableTabs[0]
      }
    }

    // 切换到下一个标签页
    router.push({ name: 'XmlEditorTab', params: { id: nextTabId } })
  }

  nextTick(() => {
    // 释放编辑器引用（不调用dispose，避免卡死）
    delete editorRefs[id]

    // 删除 Diff Editor 引用（不调用dispose，避免卡死）
    delete diffEditorInstances[id]
    delete diffEditorRefs[id]

    // 从存储中删除标签页
    delete store.xmlEditorTabs[id]

    // 清理历史记录中已删除的标签页引用
    tabHistory.value = tabHistory.value.filter(tabId =>
      tabId !== id && store.xmlEditorTabs[tabId]
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
      store.xmlEditorTabs[id]?.code || '',
      'xml'
    )
    const modifiedModel = monaco.editor.createModel(
      store.xmlEditorTabs[id]?.compareCode || '',
      'xml'
    )

    // 设置 model
    diffEditor.setModel({
      original: originalModel,
      modified: modifiedModel
    })

    // 监听修改编辑器的内容变化
    modifiedModel.onDidChangeContent(() => {
      if (store.xmlEditorTabs[id]) {
        store.xmlEditorTabs[id].compareCode = modifiedModel.getValue()
      }
    })

    // 监听原始编辑器的内容变化
    originalModel.onDidChangeContent(() => {
      if (store.xmlEditorTabs[id]) {
        store.xmlEditorTabs[id].code = originalModel.getValue()
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
.xml-editor {
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
  gap: 12px;
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
  border-left: 1px solid #d1d5db;
  border-right: 1px solid #d1d5db;
}

.tab-item + .tab-item {
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
</style>
