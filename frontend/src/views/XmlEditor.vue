<template>
  <div class="xml-editor">
    <div class="toolbar">
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
      <div class="editor-container">
        <MonacoEditor
          ref="monacoEditor"
          :value="code"
          @change="handleEditorChange"
          :options="options"
          language="xml"
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
import { FormatXML, CompressXML } from '../../wailsjs/go/main/XmlProcessor'
import { useToolsStore } from '../stores/tools'
import { storeToRefs } from 'pinia'

const { copy } = useClipboard()
const store = useToolsStore()
const { xmlEditor } = storeToRefs(store)
const code = computed({
  get: () => xmlEditor.value.code,
  set: (val) => (xmlEditor.value.code = val),
})
const monacoEditor = ref()

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
  language: 'xml',
}

const handleEditorChange = (value: string) => {
  code.value = value
}

const loadSample = () => {
  try {
    const editor = monacoEditor.value?.editor
    if (!editor) {
      ElMessage.error('编辑器未准备好')
      return
    }

    const model = editor.getModel()
    if (model) {
      const sampleXml = `<classroom><course>Introduction to Computer Science</course><instructor>Dr. Smith</instructor><students><student><student_id>001</student_id><name>Emily Johnson</name><age>19</age><gender>Female</gender><grades><grade subject="Math">A</grade><grade subject="Programming">B+</grade><grade subject="English">A-</grade></grades></student><student><student_id>002</student_id><name>Michael Smith</name><age>20</age><gender>Male</gender><grades><grade subject="Math">B</grade><grade subject="Programming">A</grade><grade subject="English">B</grade></grades></student></students></classroom>`
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

    model.setValue('')
  } catch (e) {
    ElMessage.error('清空失败')
  }
}
</script>

<style scoped>
.xml-editor {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.toolbar {
  height: 48px;
  display: flex;
  align-items: center;
  padding: 0 16px;
  border-bottom: 1px solid #eaecef;
  position: relative;
}

.tools-group {
  display: flex;
  gap: 8px;
}

.tool-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  background: #fff;
  color: #374151;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.tool-btn:hover {
  background: #f3f4f6;
  border-color: #bbc0c4;
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

:deep(.monaco-editor .scrollbar) {
  width: 14px !important;
  right: 0 !important;
}
</style>
