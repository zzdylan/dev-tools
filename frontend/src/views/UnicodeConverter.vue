<template>
  <div class="unicode-converter">
    <div class="converter-section">
      <div class="input-group">
        <div class="label">原始文本</div>
        <div class="input-with-buttons">
          <textarea
            v-model="rawText"
            placeholder="输入需要编码/解码的文本"
            class="text-area"
          ></textarea>
          <div class="button-group">
            <button class="tool-btn" @click="encode">编码</button>
            <button class="tool-btn" @click="decode">解码</button>
            <button class="tool-btn" @click="clear">清空</button>
            <button class="tool-btn" @click="copy(rawText)">复制</button>
          </div>
        </div>
      </div>

      <div class="input-group">
        <div class="label">编码/解码结果</div>
        <div class="input-with-buttons">
          <textarea
            v-model="result"
            readonly
            placeholder="编码/解码结果将显示在这里"
            class="text-area"
          ></textarea>
          <div class="button-group">
            <button class="tool-btn" @click="copy(result)">复制结果</button>
            <button class="tool-btn" @click="swapTexts">交换内容</button>
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

const clear = () => {
  rawText.value = ''
  result.value = ''
}

const copy = async (text: string) => {
  if (!text) {
    ElMessage.warning('没有可复制的内容')
    return
  }
  await copyToClipboard(text)
  ElMessage.success('已复制到剪贴板')
}

const swapTexts = () => {
  const temp = rawText.value
  rawText.value = result.value
  result.value = temp
}
</script>

<style scoped>
.unicode-converter {
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
  margin-bottom: 20px;
}

.input-group:last-child {
  margin-bottom: 0;
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
  min-height: 120px;
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
</style>
