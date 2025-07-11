<template>
  <div class="jwt-decoder">
    <!-- 顶部：标签导航 -->
    <div class="top-header">
      <div class="tab-nav">
        <button class="action-btn" @click="loadSample" title="示例">示例</button>
      </div>
      <div class="tab-actions">
        <button class="clear-btn" @click="clearAll" title="清空">× 清空</button>
      </div>
    </div>

    <!-- 底部：内容区域 -->
    <div class="content-layout">
      <!-- 输入区域 -->
      <div class="input-panel">
        <textarea
          v-model="jwtToken"
          placeholder="粘贴JWT Token..."
          class="text-area"
          @input="autoDecodeIfValid"
        ></textarea>
      </div>

      <!-- 输出区域 -->
      <div class="output-panel">
        <div v-if="error" class="error-message">
          <div class="error-icon">⚠️</div>
          <div class="error-text">{{ error }}</div>
        </div>

        <div v-else-if="jwtData" class="result-content">
          <pre class="json-result">{{ JSON.stringify(jwtData, null, 2) }}</pre>
        </div>

        <div v-else class="placeholder">
          <div class="placeholder-icon">🔍</div>
          <div class="placeholder-text">输入JWT Token开始解析</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useToolsStore } from '../stores/tools'
import { storeToRefs } from 'pinia'

const store = useToolsStore()
const { jwtDecoder } = storeToRefs(store)

const jwtToken = computed({
    get: () => jwtDecoder.value.token,
    set: (val) => {
        jwtDecoder.value.token = val
    }
})

const jwtData = ref<any>(null)
const error = ref('')

// 解析JWT Token
const decodeJwt = () => {
    error.value = ''
    jwtData.value = null

    if (!jwtToken.value.trim()) {
        error.value = '请输入JWT Token'
        return
    }

    try {
        const parts = jwtToken.value.trim().split('.')

        if (parts.length !== 3) {
            error.value = '无效的JWT格式，JWT应该包含3个部分（header.payload.signature）'
            return
        }

        // 解析Header
        const header = JSON.parse(atob(parts[0].replace(/-/g, '+').replace(/_/g, '/')))

        // 解析Payload
        const payload = JSON.parse(atob(parts[1].replace(/-/g, '+').replace(/_/g, '/')))

        // 保存原始签名
        const signature = parts[2]

        jwtData.value = {
            header,
            payload,
            signature
        }

    } catch (err) {
        error.value = 'JWT解析失败：' + (err instanceof Error ? err.message : '未知错误')
    }
}

// 加载示例
const loadSample = () => {
    jwtToken.value = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpbmZvIjp7InV1aWQiOiI1NDEzMjA4ODIyMTYyODc3OTgiLCJzZXNzaW9uX2lkIjoiZjA1NGIxMjQtZTZmNy00NGM1LTk2NmUtMzdiNWVhMDc4Y2M3IiwidGVhbV9pZCI6OX0sImV4cCI6MTc1MDgxNjExNX0.-EkJqr10jm_YmkMc0pmjuIq-hV51dthXOaJ0ClpUWsQ'
    autoDecodeIfValid()
}

// 清空所有
const clearAll = () => {
    jwtToken.value = ''
    jwtData.value = null
    error.value = ''
}

// 自动解析（当输入变化时）
const autoDecodeIfValid = () => {
    if (jwtToken.value.trim() && jwtToken.value.split('.').length === 3) {
        decodeJwt()
    } else if (!jwtToken.value.trim()) {
        jwtData.value = null
        error.value = ''
    }
}

// 监听输入变化
watch(() => jwtToken.value, () => {
    if (jwtToken.value.trim()) {
        // 延迟自动解析，避免频繁解析
        setTimeout(autoDecodeIfValid, 500)
    } else {
        jwtData.value = null
        error.value = ''
    }
})

// 组件挂载时检查是否有内容需要解析
onMounted(() => {
    if (jwtToken.value.trim()) {
        autoDecodeIfValid()
    }
})
</script>

<style scoped>
.jwt-decoder {
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
}

.action-btn {
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

.action-btn:hover {
  background: #e9ecef;
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

.text-area {
  width: 100%;
  height: 100%;
  padding: 12px;
  border: none;
  outline: none;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 11px;
  line-height: 1.4;
  resize: none;
  background: transparent;
  color: #212529;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.text-area::placeholder {
  color: #9ca3af;
  font-size: 11px;
}

.text-area:focus {
  outline: none;
  box-shadow: none;
}

.output-panel {
  padding: 0;
}

.output-panel > * {
  padding: 12px;
}

.error-message {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 8px;
  color: #dc2626;
  margin-bottom: 16px;
}

.error-icon {
  font-size: 20px;
}

.error-text {
  font-size: 11px;
  font-weight: 500;
}

.result-content {
  flex: 1;
  overflow: auto;
}

.json-result {
  background: transparent;
  padding: 0;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 11px;
  line-height: 1.4;
  overflow-x: auto;
  white-space: pre-wrap;
  color: #212529;
  margin: 0;
  border: none;
}

.placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 12px;
  color: #9ca3af;
}

.placeholder-icon {
  font-size: 48px;
  opacity: 0.5;
}

.placeholder-text {
  font-size: 11px;
  font-weight: 500;
}

@media (max-width: 1024px) {
  .content-layout {
    grid-template-columns: 1fr;
    gap: 12px;
  }
}
</style>