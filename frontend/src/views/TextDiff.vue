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

    <!-- 底部：内容区域 -->
    <div class="content-layout">
      <!-- 左侧：原始文本 -->
      <div class="input-panel">
        <div class="panel-label">原始文本</div>
        <div ref="oldEditor" class="editor-container"></div>
      </div>

      <!-- 右侧：新文本和对比结果 -->
      <div class="output-panel">
        <div class="panel-label">新文本</div>
        <div ref="newEditor" class="editor-container"></div>
        
        <!-- 对比结果 -->
        <div class="diff-result" v-if="(oldText || newText) && diffResult !== '无差异'">
          <div class="result-content" v-html="diffResult"></div>
        </div>
        
        <!-- 无差异提示 -->
        <div class="no-diff" v-if="(oldText || newText) && diffResult === '无差异'">
          <div class="no-diff-icon">✓</div>
          <div class="no-diff-text">文本内容相同</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue'
import { diffLines } from 'diff'
import { EditorView, basicSetup } from 'codemirror'
import { EditorState } from '@codemirror/state'
import { lineNumbers } from '@codemirror/view'
import { useToolsStore } from '../stores/tools'
import { storeToRefs } from 'pinia'

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
const diffResult = ref('')
const oldEditor = ref<HTMLDivElement>()
const newEditor = ref<HTMLDivElement>()
let oldEditorView: EditorView
let newEditorView: EditorView

onMounted(() => {
  // 创建左侧编辑器
  oldEditorView = new EditorView({
    state: EditorState.create({
      doc: oldText.value,
      extensions: [
        basicSetup,
        lineNumbers(),
        EditorView.updateListener.of((update) => {
          if (update.docChanged) {
            oldText.value = update.state.doc.toString()
          }
        }),
      ],
    }),
    parent: oldEditor.value!,
  })

  // 创建右侧编辑器
  newEditorView = new EditorView({
    state: EditorState.create({
      doc: newText.value,
      extensions: [
        basicSetup,
        lineNumbers(),
        EditorView.updateListener.of((update) => {
          if (update.docChanged) {
            newText.value = update.state.doc.toString()
          }
        }),
      ],
    }),
    parent: newEditor.value!,
  })
})

const generateDiff = () => {
  const diff = diffLines(oldText.value, newText.value, {
    ignoreWhitespace: ignoreWhitespace.value,
    newlineIsToken: true,
  })
  let html = ''
  let oldLineNumber = 1
  let newLineNumber = 1

  diff.forEach((part) => {
    const lines = part.value.split('\n').filter((line) => line !== '')
    lines.forEach((line) => {
      if (part.added) {
        html += `<div class="diff-line added"><span class="line-number old">  </span><span class="line-number new">+${newLineNumber}</span>${line}</div>`
        newLineNumber++
      } else if (part.removed) {
        html += `<div class="diff-line removed"><span class="line-number old">-${oldLineNumber}</span><span class="line-number new">  </span>${line}</div>`
        oldLineNumber++
      } else {
        html += `<div class="diff-line"><span class="line-number old">${oldLineNumber}</span><span class="line-number new">${newLineNumber}</span>${line}</div>`
        oldLineNumber++
        newLineNumber++
      }
    })
  })

  diffResult.value = html || '无差异'
}

watch([oldText, newText, ignoreWhitespace], generateDiff, { immediate: true })

const clearAll = () => {
  oldText.value = ''
  newText.value = ''
  if (oldEditorView) {
    oldEditorView.dispatch({
      changes: { from: 0, to: oldEditorView.state.doc.length, insert: '' }
    })
  }
  if (newEditorView) {
    newEditorView.dispatch({
      changes: { from: 0, to: newEditorView.state.doc.length, insert: '' }
    })
  }
}

const loadSample = () => {
  const sampleOld = 'Hello World\nThis is line 2\nThis is line 3'
  const sampleNew = 'Hello World\nThis is modified line 2\nThis is line 3\nThis is new line 4'
  
  oldText.value = sampleOld
  newText.value = sampleNew
  
  if (oldEditorView) {
    oldEditorView.dispatch({
      changes: { from: 0, to: oldEditorView.state.doc.length, insert: sampleOld }
    })
  }
  if (newEditorView) {
    newEditorView.dispatch({
      changes: { from: 0, to: newEditorView.state.doc.length, insert: sampleNew }
    })
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
  gap: 8px;
  background: #f8f9fa;
}

.tab-btn {
  padding: 0 10px;
  background: #f8f9fa;
  border: none;
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
}

.tab-btn:last-child {
  border-right: none;
}

.tab-btn:hover {
  background: #e9ecef;
}

.tab-btn.active {
  background: #ffffff;
  color: #212529;
  font-weight: 500;
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

.content-layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  flex: 1;
  align-items: stretch;
}

.input-panel,
.output-panel {
  border: 1px solid #d1d5db;
  background: #ffffff;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.panel-label {
  padding: 8px 12px;
  background: #f8f9fa;
  border-bottom: 1px solid #d1d5db;
  font-size: 10px;
  color: #6c757d;
  font-weight: 500;
}

.option-item {
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  padding: 2px 4px;
  transition: background-color 0.2s;
}

.option-item:hover {
  background-color: #f1f5f9;
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

.editor-container {
  flex: 1;
  min-height: 150px;
}

:deep(.cm-editor) {
  height: 100%;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 13px;
}

:deep(.cm-editor.cm-focused) {
  outline: none;
}

:deep(.cm-content) {
  padding: 14px;
}

.diff-result {
  border-top: 1px solid #d1d5db;
  background: #ffffff;
  flex: 1;
  overflow: hidden;
}

.result-content {
  padding: 0;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 11px;
  line-height: 1.4;
  background: white;
  height: 100%;
  overflow-y: auto;
}

.no-diff {
  border-top: 1px solid #d1d5db;
  background: #f8f9fa;
  padding: 20px;
  text-align: center;
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.no-diff-icon {
  font-size: 24px;
  color: #10b981;
  margin-bottom: 8px;
}

.no-diff-text {
  font-size: 11px;
  color: #6c757d;
  font-weight: 500;
}

:deep(.diff-line) {
  display: flex;
  padding: 4px 8px;
  border-bottom: 1px solid #f1f5f9;
  transition: background-color 0.2s;
  font-size: 11px;
}

:deep(.diff-line:hover) {
  background-color: #f8fafc;
}

:deep(.diff-line.added) {
  background-color: #ecfdf5;
  border-left: 3px solid #10b981;
  color: #065f46;
}

:deep(.diff-line.added:hover) {
  background-color: #d1fae5;
}

:deep(.diff-line.removed) {
  background-color: #fef2f2;
  border-left: 3px solid #ef4444;
  color: #991b1b;
}

:deep(.diff-line.removed:hover) {
  background-color: #fee2e2;
}

:deep(.line-number) {
  width: 32px;
  margin-right: 8px;
  color: #94a3b8;
  user-select: none;
  text-align: right;
  padding-right: 6px;
  font-size: 10px;
  font-weight: 500;
}

:deep(.line-number.old) {
  border-right: 1px solid #e2e8f0;
}


@media (max-width: 1024px) {
  .content-layout {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .editor-container {
    min-height: 120px;
  }
}
</style>
