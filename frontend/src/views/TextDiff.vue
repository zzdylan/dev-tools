<template>
  <div class="text-diff">
    <!-- 顶部：工具栏 -->
    <div class="top-header">
      <div class="tab-nav">
        <button class="action-btn" @click="loadSample" title="示例">示例</button>
      </div>
      <div class="tab-actions">
        <label class="option-item">
          <input type="checkbox" v-model="ignoreWhitespace" class="checkbox" />
          <span class="option-text">忽略空白</span>
        </label>
        <button class="clear-btn" @click="clearAll" title="清空">× 清空</button>
      </div>
    </div>

    <!-- 底部：Diff Editor -->
    <div class="editor-wrapper">
      <div ref="diffEditorContainer" class="diff-editor-container"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue'
import { useToolsStore } from '../stores/tools'
import { storeToRefs } from 'pinia'
import * as monaco from 'monaco-editor'

const store = useToolsStore()
const { textDiff } = storeToRefs(store)
const oldText = computed({
  get: () => textDiff.value.oldText,
  set: (val) => (textDiff.value.oldText = val),
})
const newText = computed({
  get: () => textDiff.value.newText,
  set: (val) => (textDiff.value.newText = val),
})
const ignoreWhitespace = computed({
  get: () => textDiff.value.ignoreWhitespace,
  set: (val) => (textDiff.value.ignoreWhitespace = val),
})

const diffEditorContainer = ref<HTMLElement>()
let diffEditor: monaco.editor.IStandaloneDiffEditor | null = null
let originalModel: monaco.editor.ITextModel | null = null
let modifiedModel: monaco.editor.ITextModel | null = null

onMounted(() => {
  if (!diffEditorContainer.value) return

  // 创建 Diff Editor
  diffEditor = monaco.editor.createDiffEditor(diffEditorContainer.value, {
    fontSize: 12,
    automaticLayout: true,
    renderSideBySide: true,
    enableSplitViewResizing: true,
    readOnly: false,
    minimap: { enabled: false },
    scrollBeyondLastLine: false,
    wordWrap: 'on',
    lineNumbers: 'on',
    renderIndicators: true,
    ignoreTrimWhitespace: ignoreWhitespace.value,
    renderOverviewRuler: false,
    scrollbar: {
      vertical: 'visible',
      horizontal: 'visible',
      verticalScrollbarSize: 8,
      horizontalScrollbarSize: 8,
      alwaysConsumeMouseWheel: true,
    },
    originalEditable: true,
  })

  // 创建 model
  originalModel = monaco.editor.createModel(oldText.value || '', 'plaintext')
  modifiedModel = monaco.editor.createModel(newText.value || '', 'plaintext')

  // 设置 model
  diffEditor.setModel({
    original: originalModel,
    modified: modifiedModel
  })

  // 监听修改编辑器的内容变化
  modifiedModel.onDidChangeContent(() => {
    if (modifiedModel) {
      newText.value = modifiedModel.getValue()
    }
  })

  // 监听原始编辑器的内容变化
  originalModel.onDidChangeContent(() => {
    if (originalModel) {
      oldText.value = originalModel.getValue()
    }
  })
})

// 监听 ignoreWhitespace 变化，更新编辑器选项
watch(ignoreWhitespace, (newValue) => {
  if (diffEditor) {
    diffEditor.updateOptions({
      ignoreTrimWhitespace: newValue
    })
  }
})

const clearAll = () => {
  if (originalModel && modifiedModel) {
    originalModel.setValue('')
    modifiedModel.setValue('')
  }
}

const loadSample = () => {
  const sampleOld = `Hello World
This is line 2
This is line 3
Some content here
More content
Final line`

  const sampleNew = `Hello World
This is modified line 2
This is line 3
Some different content here
More content
New additional line
Final line`

  if (originalModel && modifiedModel) {
    originalModel.setValue(sampleOld)
    modifiedModel.setValue(sampleNew)
  }
}
</script>

<style scoped>
.text-diff {
  height: 100%;
  width: 100%;
  margin: 0 !important;
  padding: 0 !important;
  background: white;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.top-header {
  display: flex;
  justify-content: space-between;
  align-items: stretch;
  padding: 0;
  margin: 0;
  background: #f8f9fa;
  height: 28px;
  flex-shrink: 0;
}

.tab-nav {
  display: flex;
  align-items: stretch;
  background: #f8f9fa;
}

.tab-actions {
  display: flex;
  align-items: stretch;
  background: #f8f9fa;
}

.action-btn {
  padding: 0 10px;
  margin: 0;
  background: #f8f9fa;
  border: none;
  border-left: 1px solid #d1d5db;
  border-right: 1px solid #d1d5db;
  font-size: 10px;
  color: #6c757d;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 45px;
  height: 100%;
  box-sizing: border-box;
}

.action-btn:hover {
  background: #e9ecef;
}

.clear-btn {
  padding: 0 10px;
  margin: 0;
  background: #f8f9fa;
  border-top: none;
  border-bottom: none;
  border-left: 1px solid #d1d5db;
  border-right: 1px solid #d1d5db;
  font-size: 10px;
  color: #6c757d;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 45px;
  height: 100%;
  box-sizing: border-box;
}

.clear-btn:hover {
  background: #e9ecef;
}

.option-item {
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  padding: 0 8px;
  transition: background-color 0.2s;
  height: 100%;
}

.option-item:hover {
  background-color: #e9ecef;
}

.checkbox {
  width: 12px;
  height: 12px;
  cursor: pointer;
  accent-color: #3b82f6;
}

.option-text {
  font-size: 10px;
  color: #6c757d;
  user-select: none;
  font-weight: 500;
}

.editor-wrapper {
  flex: 1;
  position: relative;
  overflow: hidden;
}

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
