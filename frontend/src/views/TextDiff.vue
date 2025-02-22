<template>
  <div class="text-diff">
    <div class="diff-section">
      <div class="input-panels">
        <!-- 左侧文本 -->
        <div class="input-panel">
          <div ref="oldEditor" class="editor-container"></div>
        </div>

        <!-- 右侧文本 -->
        <div class="input-panel">
          <div ref="newEditor" class="editor-container"></div>
        </div>
      </div>

      <!-- 对比结果 -->
      <div class="diff-result">
        <div class="panel-header">
          <div class="panel-title">
            <span class="title-text">差异对比结果</span>
          </div>
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
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useClipboard } from '@vueuse/core'
import { diffLines } from 'diff'
import { EditorView, basicSetup } from 'codemirror'
import { EditorState } from '@codemirror/state'
import { lineNumbers } from '@codemirror/view'

const { copy: copyToClipboard } = useClipboard()
const oldText = ref('')
const newText = ref('')
const diffResult = ref('')
const ignoreWhitespace = ref(false)
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

const copy = async (text: string) => {
  if (!text) {
    ElMessage.warning('没有可复制的内容')
    return
  }
  await copyToClipboard(text)
  ElMessage.success('已复制到剪贴板')
}
</script>

<style scoped>
.text-diff {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.diff-section {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.input-panels {
  display: flex;
  gap: 20px;
}

.input-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  overflow: hidden;
}

.editor-container {
  flex: 1;
  min-height: 200px;
  border: none;
  font-size: 14px;
  font-family: monospace;
}

:deep(.cm-editor) {
  height: 100%;
}

:deep(.cm-editor.cm-focused) {
  outline: none;
  border-color: #60a5fa;
  box-shadow: 0 0 0 2px rgba(96, 165, 250, 0.2);
}

.diff-result {
  padding: 0;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  overflow: hidden;
}

.result-content {
  padding: 12px;
  font-family: monospace;
  font-size: 14px;
  line-height: 1.5;
  white-space: pre;
  background: #f9fafb;
}

:deep(.diff-line) {
  display: flex;
  padding: 0 8px;
}

:deep(.diff-line.added) {
  background-color: #dcfce7;
  color: #166534;
}

:deep(.diff-line.removed) {
  background-color: #fee2e2;
  color: #991b1b;
}

:deep(.line-number) {
  width: 30px;
  margin-right: 16px;
  color: #6b7280;
  user-select: none;
  text-align: right;
  padding-right: 8px;
}

:deep(.line-number.old) {
  border-right: 1px solid #d1d5db;
}

.diff-options {
  display: flex;
  gap: 16px;
}

.option-item {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.checkbox {
  width: 16px;
  height: 16px;
  cursor: pointer;
}

.option-text {
  font-size: 14px;
  color: #374151;
  user-select: none;
}
</style>
