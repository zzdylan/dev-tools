<template>
  <div class="url-parser">
    <!-- 顶部：标签导航 -->
    <div class="top-header">
      <div class="tab-nav">
        <button class="tab-btn" :class="{ active: viewMode === 'table' }" @click="viewMode = 'table'">表格视图</button>
        <button class="tab-btn" :class="{ active: viewMode === 'json' }" @click="viewMode = 'json'">JSON视图</button>
      </div>
      <div class="tab-actions">
        <button class="clear-btn" @click="clear" title="清空">× 清空</button>
      </div>
    </div>

    <!-- URL输入 -->
    <textarea 
      v-model="urlText"
      placeholder="输入需要解析的URL或参数...&#10;支持格式：&#10;• https://example.com?name=value&age=18&#10;• ?name=value&age=18&#10;• name=value&age=18"
      class="url-input" 
      @input="parseUrl"
    ></textarea>

    <!-- 结果区域 -->
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
              <tr v-if="Object.keys(parsedParams).length === 0" class="empty-row">
                <td colspan="3" class="empty-cell">暂无数据</td>
              </tr>
              <tr v-else v-for="(value, key, index) in parsedParams" :key="key" class="param-row">
                <td class="index-cell">{{ index + 1 }}</td>
                <td class="key-cell">{{ key }}</td>
                <td class="value-cell">{{ value }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div v-else class="json-view">
        <pre class="json-content">{{ Object.keys(parsedParams).length === 0 ? '暂无数据' : formattedJson }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useToolsStore } from '../stores/tools'
import { storeToRefs } from 'pinia'
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
      // 完整URL解析，只提取查询参数
      const url = new URL(inputText)

      // 只解析查询参数
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

// 组件加载时解析已有的URL
onMounted(() => {
  if (urlText.value) {
    parseUrl()
  }
})

// 监听urlText变化，自动解析
watch(urlText, () => {
  parseUrl()
})
</script>

<style scoped>
.url-parser {
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
  background: #ffffff;
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
  align-items: center;
  gap: 8px;
  padding: 0 12px;
  background: #ffffff;
  height: 100%;
}

.tab-btn {
  padding: 0 10px;
  margin: 0;
  background: #f8f9fa;
  border: none;
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

.tab-btn + .tab-btn {
  border-left: 1px solid #d1d5db;
}

.tab-btn:hover {
  background: #e9ecef;
}

.tab-btn.active {
  background: #ffffff;
  color: #212529;
  font-weight: 500;
}

.tab-btn.active:hover {
  background: #ffffff;
}

.clear-btn {
  padding: 0 10px;
  margin: 0;
  background: #f8f9fa;
  border: none;
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

.url-input {
  width: 100%;
  height: 120px;
  margin-bottom: 16px;
  padding: 12px;
  border: 1px solid #d1d5db;
  border-radius: 0;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 11px;
  line-height: 1.4;
  resize: vertical;
  transition: all 0.2s;
  background: white;
  outline: none;
  flex-shrink: 0;
}

.url-input:focus {
  outline: none;
  border-color: #3b82f6;
}

.url-input::placeholder {
  color: #9ca3af;
  font-size: 11px;
}


.result-content {
  flex: 1;
  background: white;
  border: 1px solid #d1d5db;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.empty-row .empty-cell {
  text-align: center;
  color: #94a3b8;
  font-style: italic;
  padding: 20px;
}

.table-view {
  flex: 1;
  overflow: auto;
}

.table-container {
  overflow: auto;
}

.params-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 11px;
  background: white;
  user-select: text;
  -webkit-user-select: text;
  -moz-user-select: text;
  -ms-user-select: text;
}

.params-table th {
  background: #f8fafc;
  padding: 8px 12px;
  text-align: left;
  font-weight: 600;
  color: #475569;
  border-bottom: 1px solid #d1d5db;
  position: sticky;
  top: 0;
}

.params-table td {
  padding: 8px 12px;
  border-bottom: 1px solid #f1f5f9;
  color: #212529;
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
  color: #212529;
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
  flex: 1;
  overflow: auto;
}

.json-content {
  background: transparent;
  padding: 12px;
  margin: 0;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 11px;
  line-height: 1.4;
  color: #212529;
  white-space: pre-wrap;
  word-wrap: break-word;
  height: 100%;
  box-sizing: border-box;
}


@media (max-width: 768px) {
  .tab-nav {
    flex-direction: column;
  }
  
  .tab-btn {
    border-right: none;
    border-bottom: 1px solid #d1d5db;
  }
  
  .tab-btn:last-child {
    border-bottom: none;
  }
  
  .table-container {
    overflow-x: auto;
  }
}
</style>
