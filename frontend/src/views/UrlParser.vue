<template>
  <div class="url-parser">
    <div class="parser-section">
      <div class="input-group">
        <div class="label">URL</div>
        <div class="input-with-buttons">
          <textarea
            v-model="urlText"
            placeholder="输入需要解析的URL"
            class="text-area"
            @input="parseUrl"
          ></textarea>
          <div class="button-group">
            <button class="tool-btn" @click="clear">清空</button>
            <button class="tool-btn" @click="copy(urlText)">复制</button>
          </div>
        </div>
      </div>

      <div class="view-tabs">
        <button
          class="tab-btn"
          :class="{ active: viewMode === 'json' }"
          @click="viewMode = 'json'"
        >
          JSON
        </button>
        <button
          class="tab-btn"
          :class="{ active: viewMode === 'table' }"
          @click="viewMode = 'table'"
        >
          表格
        </button>
      </div>

      <div v-if="viewMode === 'json'" class="result-json">
        <pre>{{ formattedJson }}</pre>
      </div>

      <div v-else class="result-table">
        <table>
          <thead>
            <tr>
              <th>-</th>
              <th>参数名</th>
              <th>参数值</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(value, key, index) in parsedParams" :key="key">
              <td>{{ index + 1 }}</td>
              <td>{{ key }}</td>
              <td>{{ value }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useClipboard } from '@vueuse/core'

const { copy: copyToClipboard } = useClipboard()
const urlText = ref('')
const viewMode = ref<'json' | 'table'>('json')
const parsedParams = ref<Record<string, string>>({})

const parseUrl = () => {
  try {
    if (!urlText.value) {
      parsedParams.value = {}
      return
    }

    const url = new URL(urlText.value)
    const params: Record<string, string> = {}

    // 解析基本URL信息
    params.protocol = url.protocol.replace(':', '')
    params.host = url.host
    params.pathname = url.pathname

    // 解析查询参数
    url.searchParams.forEach((value, key) => {
      try {
        // 尝试解码参数值
        params[key] = decodeURIComponent(value)
      } catch {
        params[key] = value
      }
    })

    parsedParams.value = params
  } catch (error) {
    console.error('URL解析错误:', error)
  }
}

const formattedJson = computed(() => {
  return JSON.stringify(parsedParams.value, null, 2)
})

const clear = () => {
  urlText.value = ''
  parsedParams.value = {}
  ElMessage.success('已清空')
}

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
.url-parser {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.parser-section {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.input-group {
  margin-bottom: 20px;
}

.label {
  font-size: 14px;
  color: #374151;
  margin-bottom: 8px;
}

.input-with-buttons {
  display: flex;
  gap: 12px;
}

.text-area {
  flex: 1;
  min-height: 100px;
  padding: 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  font-family: monospace;
  resize: vertical;
}

.text-area:focus {
  outline: none;
  border-color: #60a5fa;
  box-shadow: 0 0 0 2px rgba(96, 165, 250, 0.2);
}

.button-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tool-btn {
  padding: 8px 16px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  background: white;
  color: #374151;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.tool-btn:hover {
  background: #f3f4f6;
  border-color: #9ca3af;
}

.view-tabs {
  margin-bottom: 16px;
  border-bottom: 1px solid #e5e7eb;
}

.tab-btn {
  padding: 8px 16px;
  border: none;
  background: none;
  font-size: 14px;
  color: #6b7280;
  cursor: pointer;
  margin-right: 16px;
  position: relative;
}

.tab-btn.active {
  color: #2563eb;
}

.tab-btn.active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  right: 0;
  height: 2px;
  background: #2563eb;
}

.result-json {
  background: #f8f9fa;
  padding: 16px;
  border-radius: 6px;
  overflow: auto;
}

.result-json pre {
  margin: 0;
  font-family: monospace;
  font-size: 12px;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.result-table {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
}

th,
td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #e5e7eb;
}

th {
  background: #f9fafb;
  font-weight: 500;
  color: #374151;
}

td {
  color: #4b5563;
}

tr:hover td {
  background: #f9fafb;
}
</style>
