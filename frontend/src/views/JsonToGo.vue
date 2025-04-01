<template>
  <div class="json-to-go">
    <div class="toolbar">
      <div class="tools-group">
        <button class="tool-btn" @click="loadSample">
          <span class="tool-icon">📝</span>
          示例
        </button>
        <!-- <button class="tool-btn" @click="convertToGo">
          <span class="tool-icon">🔄</span>
          转换
        </button> -->
        <button class="tool-btn" @click="copyToClipboard(goResult)">
          <span class="tool-icon">📋</span>
          复制Go
        </button>
        <button class="tool-btn" @click="clearContent">
          <span class="tool-icon">🗑️</span>
          清空
        </button>
        <div class="switch-container">
          <span>分开结构体:</span>
          <el-switch v-model="useFlatten" @change="convertToGo" />
        </div>
      </div>
    </div>

    <div class="editors-container">
      <div class="editor-container">
        <div class="editor-title">JSON</div>
        <MonacoEditor
          ref="jsonEditor"
          :value="code"
          @change="handleJsonChange"
          :options="jsonOptions"
          language="json"
          theme="vs"
        />
      </div>
      <div class="editor-container">
        <div class="editor-title">Go 结构体</div>
        <MonacoEditor
          ref="goEditor"
          :value="goResult"
          :options="goOptions"
          language="go"
          theme="vs"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import MonacoEditor from 'monaco-editor-vue3'
import { useClipboard } from '@vueuse/core'
import { ElMessage } from 'element-plus'
// @ts-ignore
import jsonToGo from '../utils/jsonToGo.js'

const { copy } = useClipboard()
const code = ref('')
const goResult = ref('')
const jsonEditor = ref()
const goEditor = ref()
const useFlatten = ref(false) // 默认使用内嵌模式

const jsonOptions = {
  fontSize: 12,
  tabSize: 2,
  minimap: { enabled: false },
  scrollBeyondLastLine: true,
  automaticLayout: true,
  wordWrap: 'on',
  lineNumbers: 'on',
}

const goOptions = {
  ...jsonOptions,
  readOnly: true
}

const handleJsonChange = (value: string) => {
  code.value = value
  try {
    if (value.trim()) {
      // 尝试转换，如果是有效的JSON
      JSON.parse(value)
      convertToGo()
    }
  } catch (e) {
    // JSON无效，不自动转换
  }
}

const convertToGo = () => {
  try {
    if (!code.value.trim()) {
      ElMessage.warning('请输入JSON内容')
      return
    }
    
    // 解析JSON，确保是有效的
    JSON.parse(code.value)
    
    // 将JSON转换为Go结构体，根据选项决定是否使用内联结构体
    const result = jsonToGo(code.value, 'MainType', useFlatten.value)
    if (result.error) {
      ElMessage.error('转换失败: ' + result.error)
      return
    }
    
    // 使用转换结果
    goResult.value = result.go
  } catch (error: any) {
    ElMessage.error('转换失败: ' + (error.message || String(error)))
  }
}

const loadSample = () => {
  code.value = JSON.stringify({
    name: "测试用户",
    age: 28,
    is_active: true,
    height: 175.5,
    address: {
      street: "示例街道",
      city: "示例城市",
      postal_code: "100000"
    },
    tags: ["开发", "测试", "部署"],
    projects: [
      {
        id: 1,
        title: "项目A",
        completed: true
      },
      {
        id: 2,
        title: "项目B",
        completed: false
      }
    ],
    last_login: null
  }, null, 2)
  
  convertToGo()
}

const copyToClipboard = (text: string) => {
  if (!text.trim()) {
    ElMessage.warning('没有内容可复制')
    return
  }
  
  copy(text)
  ElMessage.success('已复制到剪贴板')
}

const clearContent = () => {
  code.value = ''
  goResult.value = ''
}

onMounted(() => {
  // 组件加载时自动加载示例
  loadSample()
})
</script>

<style scoped>
.json-to-go {
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
}

.toolbar {
  display: flex;
  justify-content: flex-end;
  padding: 8px;
  background-color: #f5f5f5;
  border-bottom: 1px solid #ddd;
}

.tools-group {
  display: flex;
  gap: 8px;
  align-items: center;
}

.tool-btn {
  display: flex;
  align-items: center;
  padding: 5px 10px;
  border: none;
  border-radius: 4px;
  background-color: #fff;
  color: #333;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.2s;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.tool-btn:hover {
  background-color: #f0f0f0;
}

.tool-icon {
  margin-right: 4px;
}

.switch-container {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-left: 12px;
  font-size: 12px;
}

.editors-container {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.editor-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  border-right: 1px solid #ddd;
  overflow: hidden;
}

.editor-container:last-child {
  border-right: none;
}

.editor-title {
  padding: 8px;
  background-color: #f5f5f5;
  border-bottom: 1px solid #ddd;
  font-size: 14px;
  font-weight: 500;
}

:deep(.monaco-editor) {
  flex: 1;
}
</style> 