<template>
  <div class="unicode-converter">
    <div class="converter-container">
      <!-- 输入区域 -->
      <div class="input-section">
        <div class="section-header">
          <h3>原始文本</h3>
          <button v-if="rawText" class="btn-clear" @click="clearRawText">清空</button>
        </div>
        <div class="input-wrapper">
          <textarea
            v-model="rawText"
            placeholder="输入需要编码/解码的Unicode文本..."
            class="text-input"
          ></textarea>
        </div>
      </div>

      <!-- 转换按钮 -->
      <div class="controls-section">
        <button class="btn-encode" @click="encode">编码</button>
        <button class="btn-decode" @click="decode">解码</button>
      </div>

      <!-- 输出区域 -->
      <div class="output-section">
        <div class="section-header">
          <h3>转换结果</h3>
          <button v-if="result" class="btn-clear" @click="clearResult">清空</button>
        </div>
        <div class="output-wrapper">
          <textarea
            v-model="result"
            readonly
            placeholder="转换结果将显示在这里..."
            class="text-output"
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
const { unicodeConverter } = storeToRefs(store)

const rawText = computed({
  get: () => unicodeConverter.value.rawText,
  set: (val) => (unicodeConverter.value.rawText = val),
})

const result = computed({
  get: () => unicodeConverter.value.result,
  set: (val) => (unicodeConverter.value.result = val),
})

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

// 清空原始文本
const clearRawText = () => {
  rawText.value = ''
}

// 清空结果
const clearResult = () => {
  result.value = ''
}
</script>

<style scoped>
.unicode-converter {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
  background: #f8fafc;
  min-height: 100vh;
}

.converter-container {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  gap: 24px;
  height: calc(100vh - 60px);
  align-items: stretch;
}

.input-section,
.output-section {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  border: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
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

.input-wrapper,
.output-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.text-input,
.text-output {
  flex: 1;
  padding: 20px;
  border: none;
  outline: none;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
  resize: none;
  background: transparent;
  color: #1e293b;
}

.text-input::placeholder,
.text-output::placeholder {
  color: #94a3b8;
}

.text-output {
  background: #f8fafc;
}

.controls-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
}

.btn-encode {
  padding: 12px 24px;
  background: #10b981;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  min-width: 80px;
}

.btn-encode:hover {
  background: #059669;
  transform: translateY(-1px);
}

.btn-decode {
  padding: 12px 24px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  min-width: 80px;
}

.btn-decode:hover {
  background: #2563eb;
  transform: translateY(-1px);
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

@media (max-width: 1024px) {
  .converter-container {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .unicode-converter {
    padding: 16px;
  }
}
</style>
