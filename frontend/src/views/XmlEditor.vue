<template>
  <div class="xml-editor">
    <el-tabs v-model="activeTabName" type="border-card" addable closable @tab-add="createTab" @tab-remove="closeTab"
      @tab-change="handleTabChange" class="editor-tabs-container">
      <el-tab-pane v-for="(tab, id) in xmlEditorTabs" :key="id" :name="String(id)" :closable="String(id) !== 'default'"
        :label="String(id) === 'default' ? 'XML 编辑器' : `XML 编辑器 ${Object.keys(xmlEditorTabs).indexOf(String(id))}`">
      </el-tab-pane>
    </el-tabs>

    <div class="toolbar">
      <div class="config-wrapper">
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
        <button class="tool-btn" @click="loadSample">
          <span class="tool-icon">📝</span>
          示例
        </button>
        <button class="tool-btn" @click="formatXml">
          <span class="tool-icon">✏️</span>
          格式化
        </button>
        <button class="tool-btn" @click="compressXml">
          <span class="tool-icon">📦</span>
          压缩
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
      <div v-for="id in Object.keys(xmlEditorTabs)" :key="id" class="editor-container" v-show="id === tabId">
        <MonacoEditor :ref="(el: any) => { if (el) editorRefs[id] = el }" :value="xmlEditorTabs[id].code"
          @change="(val: string) => handleChange(val, id)" :options="options" language="xml" theme="vs" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, nextTick, reactive, watch } from 'vue'
import MonacoEditor from 'monaco-editor-vue3'
import { useClipboard } from '@vueuse/core'
import { ElMessage } from 'element-plus'
import { FormatXML, CompressXML } from '../../wailsjs/go/main/XmlProcessor'
import { onClickOutside } from '@vueuse/core'
import { useToolsStore } from '../stores/tools'
import { storeToRefs } from 'pinia'
import { useRoute, useRouter } from 'vue-router'

const { copy } = useClipboard()
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
const currentTab = computed(() => store.xmlEditorTabs[tabId.value])

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

const options = {
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
  autoIndent: 'advanced',
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

const handleChange = (value: string, id: string) => {
  if (store.xmlEditorTabs[id]) {
    store.xmlEditorTabs[id].code = value
  }
}

const loadSample = () => {
  try {
    const currentEditor = getCurrentEditor()
    if (!currentEditor?.editor) {
      ElMessage.error('编辑器未准备好')
      return
    }

    const model = currentEditor.editor.getModel()
    if (model) {
      const sampleXml = `<classroom>
  <course>Introduction to Computer Science</course>
  <instructor>Dr. Smith</instructor>
  <students>
    <student>
      <student_id>001</student_id>
      <n>Emily Johnson</n>
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
      <n>Michael Smith</n>
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
      code.value = sampleXml // 直接更新 store
      model.setValue(sampleXml)
    }
  } catch (e) {
    console.error('加载示例失败:', e)
    ElMessage.error('加载示例失败')
  }
}

const formatXml = async () => {
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
    // 释放编辑器实例
    delete editorRefs[id]
    // 从存储中删除标签页
    delete store.xmlEditorTabs[id]

    // 清理历史记录中已删除的标签页引用
    tabHistory.value = tabHistory.value.filter(tabId =>
      tabId !== id && store.xmlEditorTabs[tabId]
    )
  })
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
  flex: 0 0 48px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 16px;
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

/* 确保滚动条可见和可用 */
:deep(.monaco-editor .scrollbar) {
  width: 14px !important;
  right: 0 !important;
}

:deep(.monaco-editor .overflow-guard) {
  width: 100% !important;
  height: 100% !important;
}

:deep(.monaco-scrollable-element > .scrollbar > .slider) {
  background: rgba(100, 100, 100, 0.4) !important;
  width: 10px !important;
  left: 2px !important;
  border-radius: 5px !important;
}

:deep(.monaco-scrollable-element > .scrollbar.vertical) {
  width: 14px !important;
}

:deep(.monaco-scrollable-element > .scrollbar.vertical > .slider) {
  width: 10px !important;
  left: 2px !important;
}

:deep(.monaco-scrollable-element > .scrollbar.horizontal) {
  height: 14px !important;
  bottom: 0 !important;
}

:deep(.monaco-scrollable-element > .scrollbar.horizontal > .slider) {
  height: 10px !important;
  top: 2px !important;
}

:deep(.monaco-scrollable-element > .scrollbar.visible) {
  opacity: 1 !important;
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

.editor-tabs-container {
  flex-shrink: 0;
}

:deep(.el-tabs__header) {
  margin: 0;
  padding: 0;
  background: #f5f5f5;
  border-bottom: 1px solid #dcdcdc;
}

:deep(.el-tabs__content) {
  display: none;
}

:deep(.el-tabs--border-card) {
  border: none;
  box-shadow: none;
}

:deep(.el-tabs--border-card > .el-tabs__header) {
  border: none;
  background: transparent;
}

:deep(.el-tabs--border-card > .el-tabs__header .el-tabs__nav-wrap) {
  padding: 0;
}

:deep(.el-tabs--border-card > .el-tabs__header .el-tabs__nav) {
  border: none;
}

:deep(.el-tabs--border-card > .el-tabs__header .el-tabs__item) {
  border: 1px solid #dcdcdc;
  border-bottom: none;
  margin: 0;
  padding: 12px 20px;
  border-radius: 0;
  background: #f5f5f5;
  color: #333333;
  font-size: 13px;
  line-height: 1;
  transition: all 0.15s ease;
  min-width: 120px;
  text-align: center;
  position: relative;
}

:deep(.el-tabs--border-card > .el-tabs__header .el-tabs__item:not(:first-child)) {
  border-left: none;
}

:deep(.el-tabs--border-card > .el-tabs__header .el-tabs__item:hover:not(.is-active)) {
  background: #eeeeee;
  color: #000000;
}

:deep(.el-tabs--border-card > .el-tabs__header .el-tabs__item.is-active) {
  background: #ffffff;
  color: #1976d2;
  font-weight: normal;
  border-bottom: 1px solid #ffffff;
  z-index: 1;
}

:deep(.el-tabs__nav-next),
:deep(.el-tabs__nav-prev) {
  background: #f5f5f5;
  border: 1px solid #dcdcdc;
  border-radius: 0;
  color: #666666;
  line-height: 1;
}

:deep(.el-tabs__nav-next:hover),
:deep(.el-tabs__nav-prev:hover) {
  background: #eeeeee;
  color: #333333;
}

:deep(.el-tabs__new-tab) {
  background: #f5f5f5;
  border: 1px solid #dcdcdc;
  border-left: none;
  border-radius: 0;
  color: #666666;
  width: 40px;
  height: 45px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: normal;
  transition: all 0.15s ease;
  margin: 0;
}

:deep(.el-tabs__new-tab:hover) {
  background: #eeeeee;
  color: #333333;
}

:deep(.el-icon-close) {
  transition: all 0.15s ease;
  padding: 4px;
  margin-left: 8px;
  border-radius: 2px;
  font-size: 12px;
}

:deep(.el-icon-close:hover) {
  background: #e0e0e0;
  color: #333333;
}
</style>
