<template>
  <div class="number-converter">
    <div class="input-section">
      <div class="section-header">
        <h3>输入数值</h3>
        <button v-if="inputValue" class="btn-clear" @click="clear">清空</button>
      </div>
      <div class="input-content">
        <div class="base-selector">
          <div class="base-buttons">
            <button 
              v-for="base in bases" 
              :key="base.value"
              :class="['base-btn', { active: inputBase === base.value }]"
              @click="inputBase = base.value"
            >
              {{ base.label }}
            </button>
          </div>
        </div>
        <input
          v-model="inputValue"
          type="text"
          placeholder="输入需要转换的数值..."
          class="number-input"
        />
      </div>
    </div>

    <div class="results-section" v-if="inputValue && !isNaN(decimalValue)">
      <div class="section-header">
        <h3>转换结果</h3>
      </div>
      <div class="results-grid">
        <div class="result-card">
          <div class="result-label">二进制</div>
          <input 
            :value="binary" 
            readonly 
            class="result-input"
          />
        </div>

        <div class="result-card">
          <div class="result-label">八进制</div>
          <input 
            :value="octal" 
            readonly 
            class="result-input"
          />
        </div>

        <div class="result-card">
          <div class="result-label">十进制</div>
          <input 
            :value="decimal" 
            readonly 
            class="result-input"
          />
        </div>

        <div class="result-card">
          <div class="result-label">十六进制</div>
          <input 
            :value="hex" 
            readonly 
            class="result-input"
          />
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

// 进制选项
const bases = [
  { value: '2', label: '二进制' },
  { value: '8', label: '八进制' },
  { value: '10', label: '十进制' },
  { value: '16', label: '十六进制' }
]

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
}
</script>

<style scoped>
.number-converter {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
  background: #f8fafc;
  min-height: 100vh;
}

.input-section,
.results-section {
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

.input-content {
  padding: 24px;
  display: flex;
  gap: 16px;
  align-items: center;
}

.base-selector {
  min-width: 200px;
}

.base-buttons {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}

.base-btn {
  padding: 8px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  background: white;
  color: #64748b;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  text-align: center;
}

.base-btn:hover {
  background: #f8fafc;
  border-color: #cbd5e1;
}

.base-btn.active {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.base-btn.active:hover {
  background: #2563eb;
  border-color: #2563eb;
}

.number-input {
  flex: 1;
  padding: 12px 16px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 14px;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  transition: all 0.2s;
}

.number-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.number-input::placeholder {
  color: #94a3b8;
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

.results-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
  padding: 24px;
}

.result-card {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  padding: 20px;
  transition: all 0.2s;
  cursor: pointer;
}

.result-card:hover {
  background: #f1f5f9;
  border-color: #cbd5e1;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.result-label {
  font-size: 14px;
  font-weight: 600;
  color: #64748b;
  margin-bottom: 12px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.result-input {
  width: 100%;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 16px;
  font-weight: 500;
  color: #1e293b;
  background: white;
  padding: 12px;
  border-radius: 6px;
  border: 1px solid #e2e8f0;
  outline: none;
  cursor: pointer;
  transition: all 0.2s;
}

.result-input:hover {
  border-color: #cbd5e1;
}

.result-input:focus {
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.result-input:focus {
  user-select: all;
}

.result-input::selection {
  background: #dbeafe;
}

@media (max-width: 768px) {
  .input-content {
    flex-direction: column;
    align-items: stretch;
  }
  
  .base-selector {
    min-width: auto;
  }
  
  .base-buttons {
    grid-template-columns: repeat(4, 1fr);
  }
  
  .results-grid {
    grid-template-columns: 1fr;
  }
  
  .number-converter {
    padding: 16px;
  }
}
</style>
