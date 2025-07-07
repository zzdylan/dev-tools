<template>
  <div class="unicode-converter">
    <!-- 顶部：标签导航 -->
    <div class="top-header">
      <div class="tab-nav">
        <button class="tab-btn" :class="{ active: currentMode === 'encode' }" @click="setMode('encode')">编码</button>
        <button class="tab-btn" :class="{ active: currentMode === 'decode' }" @click="setMode('decode')">解码</button>
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
          v-model="rawText"
          placeholder="输入需要编码/解码的Unicode文本..."
          class="text-area"
        ></textarea>
      </div>

      <!-- 输出区域 -->
      <div class="output-panel">
        <textarea
          v-model="result"
          readonly
          placeholder="转换结果将显示在这里..."
          class="text-area"
        ></textarea>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useClipboard } from '@vueuse/core'
import { useToolsStore } from '../stores/tools'
import { storeToRefs } from 'pinia'

const { copy: copyToClipboard } = useClipboard()
const store = useToolsStore()
const { unicodeConverter } = storeToRefs(store)

const rawText = computed({
  get: () => unicodeConverter.value.rawText,
  set: (val) => (unicodeConverter.value.rawText = val),
})

const result = computed({
  get: () => unicodeConverter.value.result,
  set: (val) => (unicodeConverter.value.result = val),
})

const currentMode = ref('encode')

const setMode = (mode: string) => {
  currentMode.value = mode
  if (rawText.value) {
    if (mode === 'encode') {
      encode()
    } else {
      decode()
    }
  }
}

const encode = () => {
  try {
    if (!rawText.value.trim()) {
      ElMessage.warning('请输入需要编码的文本')
      return
    }
    result.value = rawText.value
      .split('')
      .map((char) => {
        const code = char.charCodeAt(0)
        return code > 127 ? `\\u${code.toString(16).padStart(4, '0')}` : char
      })
      .join('')
  } catch (error) {
    ElMessage.error('编码失败')
  }
}

const decode = () => {
  try {
    if (!rawText.value.trim()) {
      ElMessage.warning('请输入需要解码的文本')
      return
    }
    result.value = rawText.value.replace(/\\u([0-9a-fA-F]{4})/g, (_, hex) =>
      String.fromCharCode(parseInt(hex, 16))
    )
  } catch (error) {
    ElMessage.error('解码失败：无效的 Unicode 字符串')
  }
}

// 清空所有
const clearAll = () => {
  rawText.value = ''
  result.value = ''
}

// 监听输入变化，自动转换
const autoConvert = () => {
  if (!rawText.value.trim()) {
    result.value = ''
    return
  }
  
  if (currentMode.value === 'encode') {
    encode()
  } else {
    decode()
  }
}

// 监听输入文本变化
watch(rawText, autoConvert)
</script>

<style scoped>
.unicode-converter {
  height: 100%;
  background: #ffffff;
  padding: 16px;
}

.top-header {
  display: flex;
  justify-content: space-between;
  align-items: stretch;
  padding: 0;
  background: #ffffff;
  height: 28px;
  margin-bottom: 16px;
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
  height: calc(100% - 44px);
  align-items: stretch;
}

.input-panel,
.output-panel {
  border: 1px solid #d1d5db;
  background: #ffffff;
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


@media (max-width: 1024px) {
  .content-layout {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .unicode-converter {
    padding: 12px;
  }
}
</style>
