<template>
  <div class="number-converter">
    <div class="converter-section">
      <div class="input-group">
        <div class="label">输入数值</div>
        <div class="input-with-buttons">
          <div class="base-select">
            <select v-model="inputBase" class="select-control">
              <option value="2">二进制</option>
              <option value="8">八进制</option>
              <option value="10">十进制</option>
              <option value="16">十六进制</option>
            </select>
          </div>
          <input
            v-model="inputValue"
            type="text"
            placeholder="输入需要转换的数值"
            class="input-area"
            @input="handleInput"
          />
          <div class="button-group">
            <button class="tool-btn" @click="clear">清空</button>
          </div>
        </div>
      </div>

      <div class="results-grid">
        <div class="result-item">
          <div class="result-label">二进制 (Binary)</div>
          <div class="result-value">{{ binary }}</div>
          <button class="copy-btn" @click="copy(binary)">复制</button>
        </div>

        <div class="result-item">
          <div class="result-label">八进制 (Octal)</div>
          <div class="result-value">{{ octal }}</div>
          <button class="copy-btn" @click="copy(octal)">复制</button>
        </div>

        <div class="result-item">
          <div class="result-label">十进制 (Decimal)</div>
          <div class="result-value">{{ decimal }}</div>
          <button class="copy-btn" @click="copy(decimal)">复制</button>
        </div>

        <div class="result-item">
          <div class="result-label">十六进制 (Hex)</div>
          <div class="result-value">{{ hex }}</div>
          <button class="copy-btn" @click="copy(hex)">复制</button>
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

// 检测输入的是什么进制的数
const detectBase = (value: string): number => {
  value = value.toLowerCase().trim()
  if (value.startsWith('0x')) return 16
  if (value.startsWith('0b')) return 2
  if (value.startsWith('0o')) return 8
  return 10
}

// 将输入转换为十进制数
const parseInput = (value: string): number => {
  value = value.toLowerCase().trim()
  return parseInt(value, parseInt(inputBase.value))
}

const handleInput = () => {
  // inputBase.value = detectBase(inputValue.value)
}

const decimalValue = computed(() => {
  if (!inputValue.value.trim()) return NaN
  try {
    return parseInput(inputValue.value)
  } catch {
    return NaN
  }
})

const binary = computed(() => {
  if (isNaN(decimalValue.value)) return '-'
  return '0b' + decimalValue.value.toString(2).padStart(8, '0')
})

const octal = computed(() => {
  if (isNaN(decimalValue.value)) return '-'
  return '0o' + decimalValue.value.toString(8).padStart(3, '0')
})

const decimal = computed(() => {
  if (isNaN(decimalValue.value)) return '-'
  return decimalValue.value.toString(10)
})

const hex = computed(() => {
  if (isNaN(decimalValue.value)) return '-'
  return '0x' + decimalValue.value.toString(16).toUpperCase().padStart(2, '0')
})

const clear = () => {
  inputValue.value = ''
  ElMessage.success('已清空')
}

const copy = async (text: string) => {
  if (!text || text === '-') {
    ElMessage.warning('没有可复制的内容')
    return
  }
  await copyToClipboard(text)
  ElMessage.success('已复制到剪贴板')
}
</script>

<style scoped>
.number-converter {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.converter-section {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.input-group {
  margin-bottom: 24px;
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

.base-select {
  min-width: 100px;
}

.select-control {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  background: white;
  cursor: pointer;
}

.select-control:focus {
  outline: none;
  border-color: #60a5fa;
  box-shadow: 0 0 0 2px rgba(96, 165, 250, 0.2);
}

.input-area {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  font-family: monospace;
}

.input-area:focus {
  outline: none;
  border-color: #60a5fa;
  box-shadow: 0 0 0 2px rgba(96, 165, 250, 0.2);
}

.button-group {
  display: flex;
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

.results-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.result-item {
  background: #f9fafb;
  padding: 16px;
  border-radius: 6px;
  border: 1px solid #e5e7eb;
  position: relative;
}

.result-label {
  font-size: 12px;
  color: #6b7280;
  margin-bottom: 8px;
}

.result-value {
  font-family: monospace;
  font-size: 14px;
  color: #111827;
  word-break: break-all;
}

.copy-btn {
  position: absolute;
  top: 12px;
  right: 12px;
  padding: 4px 8px;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  background: white;
  color: #374151;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.copy-btn:hover {
  background: #f3f4f6;
  border-color: #9ca3af;
}
</style>
