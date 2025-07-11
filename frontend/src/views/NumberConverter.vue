<template>
  <div class="number-converter">
    <!-- 顶部：标签导航 -->
    <div class="top-header">
      <div class="tab-nav">
        <button class="tab-btn" :class="{ active: inputBase === '2' }" @click="setInputBase('2')">二进制</button>
        <button class="tab-btn" :class="{ active: inputBase === '8' }" @click="setInputBase('8')">八进制</button>
        <button class="tab-btn" :class="{ active: inputBase === '10' }" @click="setInputBase('10')">十进制</button>
        <button class="tab-btn" :class="{ active: inputBase === '16' }" @click="setInputBase('16')">十六进制</button>
        <button class="action-btn" @click="loadSample" title="示例">示例</button>
      </div>
      <div class="tab-actions">
        <button class="clear-btn" @click="clear" title="清空">× 清空</button>
      </div>
    </div>

    <!-- 底部：内容区域 -->
    <div class="content-layout">
      <!-- 输入区域 -->
      <div class="input-panel">
        <textarea
          v-model="inputValue"
          :placeholder="getPlaceholder()"
          class="text-area"
          @input="autoConvert"
        ></textarea>
      </div>

      <!-- 输出区域 -->
      <div class="output-panel">
        <div class="results-container">
          <div class="result-row">
            <span class="result-label">二进制：</span>
            <span class="result-value">{{ binary }}</span>
          </div>
          <div class="result-row">
            <span class="result-label">八进制：</span>
            <span class="result-value">{{ octal }}</span>
          </div>
          <div class="result-row">
            <span class="result-label">十进制：</span>
            <span class="result-value">{{ decimal }}</span>
          </div>
          <div class="result-row">
            <span class="result-label">十六进制：</span>
            <span class="result-value">{{ hex }}</span>
          </div>
        </div>
      </div>
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
const { numberConverter } = storeToRefs(store)

const inputValue = computed({
  get: () => numberConverter.value.inputValue,
  set: (val) => (numberConverter.value.inputValue = val),
})

const inputBase = computed({
  get: () => numberConverter.value.inputBase,
  set: (val) => (numberConverter.value.inputBase = val),
})

const setInputBase = (base: string) => {
  inputBase.value = base
  if (inputValue.value.trim()) {
    autoConvert()
  }
}

const getPlaceholder = () => {
  const placeholders = {
    '2': '输入二进制数字 (例如: 1010)',
    '8': '输入八进制数字 (例如: 755)',
    '10': '输入十进制数字 (例如: 123)',
    '16': '输入十六进制数字 (例如: FF)'
  }
  return placeholders[inputBase.value as keyof typeof placeholders] || '输入数字...'
}

const loadSample = () => {
  const samples = {
    '2': '1010110',
    '8': '777',
    '10': '255',
    '16': 'FF'
  }
  inputValue.value = samples[inputBase.value as keyof typeof samples] || '255'
  autoConvert()
}

// 将输入转换为十进制数
const parseInput = (value: string): number => {
  value = value.trim()
  // 移除常见的前缀
  if (value.startsWith('0x') || value.startsWith('0X')) {
    value = value.slice(2)
  } else if (value.startsWith('0b') || value.startsWith('0B')) {
    value = value.slice(2)
  } else if (value.startsWith('0o') || value.startsWith('0O')) {
    value = value.slice(2)
  }
  
  return parseInt(value, parseInt(inputBase.value))
}

const decimalValue = computed(() => {
  if (!inputValue.value.trim()) return NaN
  try {
    const result = parseInput(inputValue.value)
    return isNaN(result) ? NaN : result
  } catch {
    return NaN
  }
})

const binary = computed(() => {
  if (isNaN(decimalValue.value)) return '-'
  return decimalValue.value.toString(2)
})

const octal = computed(() => {
  if (isNaN(decimalValue.value)) return '-'
  return decimalValue.value.toString(8)
})

const decimal = computed(() => {
  if (isNaN(decimalValue.value)) return '-'
  return decimalValue.value.toString(10)
})

const hex = computed(() => {
  if (isNaN(decimalValue.value)) return '-'
  return decimalValue.value.toString(16).toUpperCase()
})

const autoConvert = () => {
  // 自动转换，无需手动操作
}

const clear = () => {
  inputValue.value = ''
}
</script>

<style scoped>
.number-converter {
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
  border: 1px solid #d1d5db;
}

.tab-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 12px;
  background: #ffffff;
}

.tab-btn {
  padding: 0 10px;
  background: #f8f9fa;
  border: none;
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
}

.tab-btn:last-child {
  border-right: none;
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

.action-btn {
  padding: 0 10px;
  background: #f8f9fa;
  border: none;
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
}

.action-btn:hover {
  background: #e9ecef;
}

.clear-btn {
  padding: 0 10px;
  background: #f8f9fa;
  border: 1px solid #d1d5db;
  font-size: 10px;
  color: #6c757d;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 45px;
  height: 100%;
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

.results-container {
  padding: 12px;
  height: 100%;
  overflow-y: auto;
}

.result-row {
  display: flex;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f1f5f9;
}

.result-row:last-child {
  border-bottom: none;
}

.result-label {
  font-size: 11px;
  color: #6c757d;
  min-width: 60px;
  font-weight: 500;
}

.result-value {
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 11px;
  color: #212529;
  word-break: break-all;
  user-select: text;
}

@media (max-width: 1024px) {
  .content-layout {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .number-converter {
    padding: 12px;
  }
}
</style>