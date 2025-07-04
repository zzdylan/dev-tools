<template>
  <div class="json-to-go">
    <div class="converter-container">
      <!-- JSON输入区域 -->
      <div class="input-section">
        <div class="section-header">
          <h3>JSON 数据</h3>
          <div class="header-controls">
            <button class="btn-sample" @click="loadSample">示例</button>
            <button class="btn-clear" @click="clearContent">清空</button>
          </div>
        </div>
        <div class="editor-wrapper">
          <MonacoEditor 
            ref="jsonEditor" 
            :value="code" 
            @change="handleJsonChange" 
            :options="jsonOptions" 
            language="json"
            theme="vs" 
          />
        </div>
      </div>

      <!-- Go结构体输出区域 -->
      <div class="output-section">
        <div class="section-header">
          <h3>Go 结构体</h3>
          <div class="header-controls">
            <div class="switch-container">
              <span>分开结构体</span>
              <el-switch v-model="useFlatten" @change="convertToGo" size="small" />
            </div>
            <button class="btn-copy" @click="copyToClipboard(goResult)">复制结果</button>
          </div>
        </div>
        <div class="editor-wrapper">
          <MonacoEditor 
            ref="goEditor" 
            :value="goResult" 
            @change="handleGoChange" 
            :options="goOptions" 
            language="go"
            theme="vs" 
          />
        </div>
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
  readOnly: false
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

const handleGoChange = (value: string) => {
  goResult.value = value
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
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
  background: #f8fafc;
  min-height: 100vh;
}

.converter-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
  height: calc(100vh - 80px);
}

.input-section,
.output-section {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  border: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #e2e8f0;
  background: #f8fafc;
}

.section-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #1e293b;
}

.header-controls {
  display: flex;
  gap: 12px;
  align-items: center;
}

.btn-sample {
  padding: 8px 16px;
  background: #8b5cf6;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-sample:hover {
  background: #7c3aed;
}

.btn-clear {
  padding: 6px 12px;
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fecaca;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-clear:hover {
  background: #fee2e2;
}

.btn-copy {
  padding: 6px 12px;
  background: #f1f5f9;
  color: #475569;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-copy:hover {
  background: #e2e8f0;
}

.switch-container {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #64748b;
}

.editor-wrapper {
  flex: 1;
  border-radius: 0 0 12px 12px;
  overflow: hidden;
}

:deep(.monaco-editor) {
  flex: 1;
}

@media (max-width: 1024px) {
  .converter-container {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .json-to-go {
    padding: 16px;
  }
}
</style>