<template>
  <div class="base64-image">
    <div class="tool-section">
      <div class="tool-layout">
        <!-- 左侧：图片上传和预览 -->
        <div class="left-panel">
          <div class="upload-area" @drop.prevent="handleDrop" @dragover.prevent>
            <input
              type="file"
              ref="fileInput"
              accept="image/*"
              class="hidden"
              @change="handleFileSelect"
            />
            <div
              v-if="!imageUrl"
              class="upload-content"
              @click="triggerFileInput"
            >
              <div class="upload-icon">📤</div>
              <div class="upload-text">
                点击或拖拽图片到此处
                <div class="upload-hint">支持 jpg、png、gif 格式</div>
              </div>
            </div>
            <img v-else :src="imageUrl" class="preview-image" />
          </div>
          <div class="button-group">
            <button class="tool-btn" @click="triggerFileInput">选择图片</button>
            <button class="tool-btn" @click="clear">清空</button>
          </div>
        </div>

        <!-- 右侧：Base64 结果 -->
        <div class="right-panel">
          <textarea
            v-model="base64Result"
            placeholder="Base64 编码结果将显示在这里"
            class="result-area"
          ></textarea>
          <div class="button-group">
            <button class="tool-btn" @click="decodeBase64">解码为图片</button>
            <button class="tool-btn" @click="copyResult">复制结果</button>
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

const { copy } = useClipboard()
const store = useToolsStore()
const { base64Image } = storeToRefs(store)
const imageUrl = computed({
  get: () => base64Image.value.imageUrl,
  set: (val) => (base64Image.value.imageUrl = val),
})
const base64Result = computed({
  get: () => base64Image.value.base64Result,
  set: (val) => (base64Image.value.base64Result = val),
})
const fileInput = ref<HTMLInputElement | null>(null)

const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileSelect = (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (file) {
    handleImage(file)
  }
}

const handleDrop = (event: DragEvent) => {
  const file = event.dataTransfer?.files[0]
  if (file) {
    handleImage(file)
  }
}

const handleImage = (file: File) => {
  if (!file.type.startsWith('image/')) {
    ElMessage.error('请选择图片文件')
    return
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    const result = e.target?.result as string
    imageUrl.value = result
    base64Result.value = result
  }
  reader.readAsDataURL(file)
}

const clear = () => {
  imageUrl.value = ''
  base64Result.value = ''
  if (fileInput.value) {
    fileInput.value.value = ''
  }
  ElMessage.success('已清空')
}

const copyResult = async () => {
  if (!base64Result.value) {
    ElMessage.warning('没有可复制的内容')
    return
  }
  await copy(base64Result.value)
  ElMessage.success('已复制到剪贴板')
}

const decodeBase64 = () => {
  try {
    if (!base64Result.value.trim()) {
      ElMessage.warning('请输入需要解码的 Base64 字符串')
      return
    }

    // 验证是否是合法的 Base64 图片字符串
    if (!base64Result.value.startsWith('data:image/')) {
      ElMessage.error('无效的 Base64 图片格式')
      return
    }

    imageUrl.value = base64Result.value
    ElMessage.success('解码成功')
  } catch (error) {
    console.error('解码失败:', error)
    ElMessage.error('解码失败：无效的 Base64 字符串')
  }
}
</script>

<style scoped>
.base64-image {
  height: 100%;
}

.tool-section {
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  height: 100%;
}

.tool-layout {
  display: flex;
  height: 100%;
}

.left-panel {
  flex: 1;
  padding: 20px;
  border-right: 1px solid #e5e7eb;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.right-panel {
  width: 500px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.upload-area {
  flex: 1;
  border: 2px dashed #d1d5db;
  border-radius: 8px;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  transition: all 0.2s;
  overflow: hidden;
}

.upload-area:hover {
  border-color: #60a5fa;
  background: #f8fafc;
}

.preview-image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.upload-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 40px;
}

.upload-icon {
  font-size: 48px;
}

.upload-text {
  color: #374151;
  font-size: 14px;
  text-align: center;
}

.upload-hint {
  color: #6b7280;
  font-size: 12px;
  margin-top: 4px;
}

.result-area {
  flex: 1;
  padding: 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 12px;
  font-family: monospace;
  resize: none;
  white-space: pre-wrap;
  word-break: break-all;
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
}

.tool-btn:hover {
  background: #f3f4f6;
  border-color: #9ca3af;
}

.hidden {
  display: none;
}
</style>
