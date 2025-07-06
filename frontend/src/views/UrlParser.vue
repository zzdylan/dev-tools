<template>
  <div class="url-parser">
    <!-- URL输入区域 -->
    <div class="input-section">
      <div class="section-header">
        <h3>URL 解析</h3>
        <button v-if="urlText" class="btn-clear" @click="clear">清空</button>
      </div>
      <div class="input-wrapper">
        <textarea 
          v-model="urlText"
          placeholder="输入需要解析的URL或参数...&#10;支持格式：&#10;• https://example.com?name=value&age=18&#10;• ?name=value&age=18&#10;• name=value&age=18"
          class="url-input" 
          @input="parseUrl"
        ></textarea>
      </div>
    </div>

    <!-- 解析结果区域 -->
    <div class="result-section" v-if="urlText && Object.keys(parsedParams).length > 0">
      <div class="section-header">
        <h3>解析结果</h3>
        <div class="view-tabs">
          <button 
            class="tab-btn" 
            :class="{ active: viewMode === 'table' }" 
            @click="viewMode = 'table'"
          >
            表格视图
          </button>
          <button 
            class="tab-btn" 
            :class="{ active: viewMode === 'json' }" 
            @click="viewMode = 'json'"
          >
            JSON视图
          </button>
        </div>
      </div>
      
      <div class="result-content">
        <div v-if="viewMode === 'table'" class="table-view">
          <div class="table-container">
            <table class="params-table">
              <thead>
                <tr>
                  <th class="index-col">#</th>
                  <th class="key-col">参数名</th>
                  <th class="value-col">参数值</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(value, key, index) in parsedParams" :key="key" class="param-row">
                  <td class="index-cell">{{ index + 1 }}</td>
                  <td class="key-cell">{{ key }}</td>
                  <td class="value-cell">{{ value }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div v-else class="json-view">
          <pre class="json-content">{{ formattedJson }}</pre>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div class="empty-state" v-if="!urlText || Object.keys(parsedParams).length === 0">
      <div class="empty-icon">🔗</div>
      <div class="empty-text">输入URL或参数字符串开始解析</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useClipboard } from '@vueuse/core'
import { useToolsStore } from '../stores/tools'
import { storeToRefs } from 'pinia'

const { copy: copyToClipboard } = useClipboard()
const store = useToolsStore()
const { urlParser } = storeToRefs(store)
const urlText = computed({
  get: () => urlParser.value.urlText,
  set: (val) => (urlParser.value.urlText = val),
})
const viewMode = computed({
  get: () => urlParser.value.viewMode,
  set: (val) => (urlParser.value.viewMode = val),
})
const parsedParams = ref<Record<string, string>>({})

const parseUrl = () => {
  try {
    if (!urlText.value) {
      parsedParams.value = {}
      return
    }

    const params: Record<string, string> = {}
    let inputText = urlText.value.trim()

    // 检查是否为完整的URL（包含协议）
    if (inputText.match(/^https?:\/\//)) {
      // 完整URL解析
      const url = new URL(inputText)

      // 解析基本URL信息
      params.protocol = url.protocol.replace(':', '')
      params.host = url.host
      params.hostname = url.hostname
      params.port = url.port || (url.protocol === 'https:' ? '443' : '80')
      params.pathname = url.pathname
      params.hash = url.hash

      // 解析查询参数
      url.searchParams.forEach((value, key) => {
        try {
          params[key] = decodeURIComponent(value)
        } catch {
          params[key] = value
        }
      })
    } else {
      // 处理纯参数或带问号的参数
      let queryString = inputText

      // 如果不以?开头，添加?
      if (!queryString.startsWith('?')) {
        queryString = '?' + queryString
      }

      // 使用URLSearchParams解析查询参数
      const searchParams = new URLSearchParams(queryString)

      searchParams.forEach((value, key) => {
        try {
          params[key] = decodeURIComponent(value)
        } catch {
          params[key] = value
        }
      })
    }

    parsedParams.value = params
  } catch (error) {
    console.error('URL解析错误:', error)
    // 如果解析失败，尝试简单的手动解析
    try {
      const params: Record<string, string> = {}
      let inputText = urlText.value.trim()

      // 移除开头的?
      if (inputText.startsWith('?')) {
        inputText = inputText.substring(1)
      }

      // 按&分割参数
      const pairs = inputText.split('&')
      pairs.forEach(pair => {
        const [key, value] = pair.split('=')
        if (key) {
          try {
            params[decodeURIComponent(key)] = value ? decodeURIComponent(value) : ''
          } catch {
            params[key] = value || ''
          }
        }
      })

      parsedParams.value = params
    } catch (fallbackError) {
      console.error('手动解析也失败:', fallbackError)
      parsedParams.value = {}
    }
  }
}

const formattedJson = computed(() => {
  return JSON.stringify(parsedParams.value, null, 2)
})

const clear = () => {
  urlText.value = ''
  parsedParams.value = {}
}
</script>

<style scoped>
.url-parser {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
  background: #f8fafc;
  min-height: 100vh;
}

.input-section,
.result-section {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  border: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  margin-bottom: 24px;
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

.input-wrapper {
  padding: 24px;
}

.url-input {
  width: 100%;
  min-height: 120px;
  padding: 16px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
  resize: vertical;
  transition: all 0.2s;
  background: #fafbfc;
}

.url-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  background: white;
}

.url-input::placeholder {
  color: #94a3b8;
  line-height: 1.4;
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

.view-tabs {
  display: flex;
  gap: 8px;
}

.tab-btn {
  padding: 8px 16px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  background: white;
  color: #64748b;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.tab-btn:hover {
  background: #f8fafc;
  border-color: #cbd5e1;
}

.tab-btn.active {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.tab-btn.active:hover {
  background: #2563eb;
  border-color: #2563eb;
}

.result-content {
  flex: 1;
}

.table-view {
  padding: 24px;
}

.table-container {
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  overflow: hidden;
}

.params-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
  background: white;
  user-select: text;
  -webkit-user-select: text;
  -moz-user-select: text;
  -ms-user-select: text;
}

.params-table th {
  background: #f8fafc;
  padding: 16px;
  text-align: left;
  font-weight: 600;
  color: #475569;
  border-bottom: 1px solid #e2e8f0;
}

.params-table td {
  padding: 16px;
  border-bottom: 1px solid #f1f5f9;
  color: #1e293b;
}

.index-col {
  width: 60px;
}

.key-col {
  width: 200px;
}

.value-col {
  min-width: 200px;
}

.param-row:hover td {
  background: #f8fafc;
}

.index-cell {
  color: #94a3b8;
  font-weight: 500;
  text-align: center;
}

.key-cell {
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-weight: 600;
  color: #3b82f6;
  user-select: text;
  -webkit-user-select: text;
  -moz-user-select: text;
  -ms-user-select: text;
}

.value-cell {
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  word-break: break-all;
  user-select: text;
  -webkit-user-select: text;
  -moz-user-select: text;
  -ms-user-select: text;
}

.json-view {
  padding: 24px;
}

.json-content {
  background: #f8fafc;
  padding: 20px;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  margin: 0;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.6;
  color: #1e293b;
  white-space: pre-wrap;
  word-wrap: break-word;
  overflow-x: auto;
}

.empty-state {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  border: 1px solid #e2e8f0;
  padding: 60px;
  text-align: center;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.6;
}

.empty-text {
  font-size: 16px;
  color: #64748b;
  font-weight: 500;
}

@media (max-width: 1024px) {
  .url-parser {
    padding: 16px;
  }
  
  .view-tabs {
    flex-direction: column;
    gap: 4px;
  }
  
  .tab-btn {
    text-align: center;
  }
  
  .table-container {
    overflow-x: auto;
  }
  
  .empty-state {
    padding: 40px 20px;
  }
}
</style>
