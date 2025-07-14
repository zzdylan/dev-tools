<template>
  <div class="json-to-go">
    <!-- 顶部：标签导航 -->
    <div class="top-header">
      <div class="tab-nav">
        <button class="action-btn" @click="loadSample" title="示例">示例</button>
        <div class="switch-container">
          <span>分开结构体</span>
          <el-switch v-model="useFlatten" @change="convertToGo" size="small" />
        </div>
      </div>
      <div class="tab-actions">
        <button class="clear-btn" @click="clearContent" title="清空">× 清空</button>
      </div>
    </div>

    <!-- 底部：内容区域 -->
    <div class="content-layout">
      <!-- JSON输入区域 -->
      <div class="input-panel">
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
      <div class="output-panel">
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
import { ElMessage } from 'element-plus'
// @ts-ignore
import jsonToGo from '../utils/jsonToGo.js'
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

const clearContent = () => {
  code.value = ''
  goResult.value = ''
}

onMounted(() => {
  // 组件加载时不自动加载示例，保持空白状态
})

</script>

<style scoped>
.json-to-go {
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


.switch-container {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 12px;
  margin: 0;
  height: 28px;
  border: none;
  border-left: none;
  border-right: 1px solid #d1d5db;
  background: #f8f9fa;
  font-size: 11px;
  color: #6c757d;
  border-radius: 0;
  transition: all 0.2s;
  box-sizing: border-box;
}

.switch-container:hover {
  background: #e9ecef;
}

.editor-wrapper {
  flex: 1;
  overflow: hidden;
}

:deep(.monaco-editor) {
  flex: 1;
}

@media (max-width: 1024px) {
  .content-layout {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .tab-actions {
    gap: 4px;
  }
  
  .switch-container {
    padding: 0 8px;
    gap: 6px;
    font-size: 10px;
  }
  
  .switch-container span {
    font-size: 9px;
  }
}
</style>