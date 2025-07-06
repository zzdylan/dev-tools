<template>
  <div class="text-diff">
    <!-- 输入区域 -->
    <div class="input-section">
      <div class="input-panels">
        <!-- 左侧文本 -->
        <div class="input-panel">
          <div class="panel-header">
            <h3>原始文本</h3>
            <button v-if="oldText" class="btn-clear" @click="clearOldText">清空</button>
          </div>
          <div ref="oldEditor" class="editor-container"></div>
        </div>

        <!-- 右侧文本 -->
        <div class="input-panel">
          <div class="panel-header">
            <h3>新文本</h3>
            <button v-if="newText" class="btn-clear" @click="clearNewText">清空</button>
          </div>
          <div ref="newEditor" class="editor-container"></div>
        </div>
      </div>
    </div>

    <!-- 对比结果 -->
    <div class="diff-result" v-if="(oldText || newText) && diffResult !== '无差异'">
      <div class="result-header">
        <h3>差异对比结果</h3>
        <div class="diff-options">
          <label class="option-item">
            <input
              type="checkbox"
              v-model="ignoreWhitespace"
              class="checkbox"
            />
            <span class="option-text">忽略空白字符</span>
          </label>
        </div>
      </div>
      <div class="result-content" v-html="diffResult"></div>
    </div>

    <!-- 无差异提示 -->
    <div class="no-diff" v-if="(oldText || newText) && diffResult === '无差异'">
      <div class="no-diff-icon">✓</div>
      <div class="no-diff-text">文本内容相同，没有差异</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useClipboard } from '@vueuse/core'
import { diffLines } from 'diff'
import { EditorView, basicSetup } from 'codemirror'
import { EditorState } from '@codemirror/state'
import { lineNumbers } from '@codemirror/view'
import { useToolsStore } from '../stores/tools'
import { storeToRefs } from 'pinia'

const { copy: copyToClipboard } = useClipboard()
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

const clearOldText = () => {
  oldText.value = ''
  if (oldEditorView) {
    oldEditorView.dispatch({
      changes: { from: 0, to: oldEditorView.state.doc.length, insert: '' }
    })
  }
}

const clearNewText = () => {
  newText.value = ''
  if (newEditorView) {
    newEditorView.dispatch({
      changes: { from: 0, to: newEditorView.state.doc.length, insert: '' }
    })
  }
}
</script>

<style scoped>
.text-diff {
  padding: 16px;
  max-width: 1400px;
  margin: 0 auto;
  background: #f8fafc;
  min-height: 100vh;
}

.input-section {
  margin-bottom: 16px;
}

.input-panels {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.input-panel {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  border: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.panel-header,
.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #e2e8f0;
  background: #f8fafc;
}

.panel-header h3,
.result-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #1e293b;
}

.btn-clear {
  padding: 4px 8px;
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fecaca;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-clear:hover {
  background: #fee2e2;
}

.editor-container {
  flex: 1;
  min-height: 200px;
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
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  border: 1px solid #e2e8f0;
  overflow: hidden;
}

.result-content {
  padding: 0;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
  background: white;
  max-height: 300px;
  overflow-y: auto;
}

.no-diff {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  border: 1px solid #e2e8f0;
  padding: 30px;
  text-align: center;
}

.no-diff-icon {
  font-size: 36px;
  color: #10b981;
  margin-bottom: 12px;
}

.no-diff-text {
  font-size: 14px;
  color: #64748b;
  font-weight: 500;
}

:deep(.diff-line) {
  display: flex;
  padding: 6px 12px;
  border-bottom: 1px solid #f1f5f9;
  transition: background-color 0.2s;
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
  width: 36px;
  margin-right: 12px;
  color: #94a3b8;
  user-select: none;
  text-align: right;
  padding-right: 8px;
  font-size: 11px;
  font-weight: 500;
}

:deep(.line-number.old) {
  border-right: 1px solid #e2e8f0;
}

.diff-options {
  display: flex;
  gap: 12px;
}

.option-item {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  padding: 3px 6px;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.option-item:hover {
  background-color: #f1f5f9;
}

.checkbox {
  width: 14px;
  height: 14px;
  cursor: pointer;
  accent-color: #3b82f6;
}

.option-text {
  font-size: 12px;
  color: #475569;
  user-select: none;
  font-weight: 500;
}

@media (max-width: 1024px) {
  .input-panels {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .text-diff {
    padding: 16px;
  }
  
  .editor-container {
    min-height: 150px;
  }
}
</style>
