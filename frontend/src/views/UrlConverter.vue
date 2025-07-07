<template>
  <div class="url-converter">
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
    <div class="content-container">
      <!-- 输入区域 -->
      <div class="input-section">
        <div class="input-wrapper">
          <textarea
            v-model="rawText"
            placeholder="输入需要编码/解码的URL文本..."
            class="text-input"
            autocomplete="off"
            spellcheck="false"
          ></textarea>
        </div>
      </div>

      <!-- 输出区域 -->
      <div class="output-section">
        <div class="output-wrapper">
          <textarea
            v-model="result"
            readonly
            placeholder="转换结果将显示在这里..."
            class="text-output"
            autocomplete="off"
            spellcheck="false"
          ></textarea>
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
const { urlConverter } = storeToRefs(store)
const rawText = computed({
  get: () => urlConverter.value.rawText,
  set: (val) => (urlConverter.value.rawText = val),
})
const result = computed({
  get: () => urlConverter.value.result,
  set: (val) => (urlConverter.value.result = val),
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
    result.value = encodeURIComponent(rawText.value)
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
    result.value = decodeURIComponent(rawText.value)
  } catch (error) {
    ElMessage.error('解码失败：无效的编码字符串')
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
import { watch } from 'vue'
watch(rawText, autoConvert)
</script>

<style scoped>
.url-converter {
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
  height: 100%;
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

.tab-btn.active {
  background: #ffffff;
  color: #212529;
  font-weight: 500;
}

.tab-btn:hover {
  background: #e9ecef;
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

.content-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  height: calc(100% - 44px);
  align-items: stretch;
}

.input-section,
.output-section {
  border: 1px solid #d1d5db;
  background: #ffffff;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.output-section {
  border-left: none;
}


.input-wrapper,
.output-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.text-input,
.text-output {
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

.text-input:focus,
.text-output:focus {
  outline: none;
  box-shadow: none;
}

.text-input::placeholder,
.text-output::placeholder {
  color: #9ca3af;
  font-size: 11px;
}

.text-output {
  background: transparent;
}

@media (max-width: 1024px) {
  .content-container {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .input-section,
  .output-section {
    border: 1px solid #d1d5db;
  }
}
</style>
