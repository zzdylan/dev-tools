<template>
  <div class="base64-image">
    <div class="tool-section">
      <div class="tool-layout">
        <!-- 左侧：图片上传和预览 -->
        <div class="left-panel">
          <div class="section-header">
            <h3>图片</h3>
            <button v-if="imageUrl" class="btn-clear" @click="clearImage">清空</button>
          </div>
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
            <img v-else :src="imageUrl" class="preview-image" @click="triggerFileInput" />
          </div>
        </div>

        <!-- 中间：转换按钮 -->
        <div class="controls-section">
          <button class="btn-encode" @click="encodeImageToBase64">编码</button>
          <button class="btn-decode" @click="decodeBase64ToImage">解码</button>
        </div>

        <!-- 右侧：Base64 结果 -->
        <div class="right-panel">
          <div class="section-header">
            <h3>Base64 结果</h3>
            <button v-if="base64Result" class="btn-clear" @click="clearBase64">清空</button>
          </div>
          <textarea
            v-model="base64Result"
            placeholder="Base64 编码结果将显示在这里..."
            class="result-area"
          ></textarea>
        </div>
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
    // 不自动转换，只显示图片
  }
  reader.readAsDataURL(file)
}

// 手动编码图片为Base64
const encodeImageToBase64 = () => {
  if (!imageUrl.value) {
    ElMessage.warning('请先上传图片')
    return
  }
  
  try {
    base64Result.value = imageUrl.value
    ElMessage.success('编码成功')
  } catch (error) {
    ElMessage.error('编码失败')
  }
}

// 手动解码Base64为图片
const decodeBase64ToImage = () => {
  if (!base64Result.value || !base64Result.value.trim()) {
    ElMessage.warning('请输入Base64字符串')
    return
  }
  
  try {
    // 清理Base64字符串，移除换行符和空格
    const cleanBase64 = base64Result.value.replace(/\s+/g, '')
    
    // 如果已经是完整的data URL格式
    if (cleanBase64.startsWith('data:image/')) {
      imageUrl.value = cleanBase64
      ElMessage.success('解码成功')
    } 
    // 如果只是Base64字符串，尝试添加data URL前缀
    else if (cleanBase64.length > 100) {
      // 默认尝试jpeg格式
      const dataUrl = `data:image/jpeg;base64,${cleanBase64}`
      imageUrl.value = dataUrl
      ElMessage.success('解码成功')
    } else {
      ElMessage.error('无效的Base64字符串')
    }
  } catch (error) {
    ElMessage.error('解码失败')
  }
}

// 清空图片
const clearImage = () => {
  imageUrl.value = ''
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

// 清空 Base64
const clearBase64 = () => {
  base64Result.value = ''
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
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  gap: 24px;
  height: 100%;
  align-items: stretch;
}

.left-panel,
.right-panel {
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

.upload-area {
  flex: 1;
  margin: 20px;
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
  padding: 20px;
  border: none;
  outline: none;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
  resize: none;
  background: transparent;
  color: #1e293b;
  white-space: normal;
  word-wrap: break-word;
  word-break: normal;
}

.result-area::placeholder {
  color: #94a3b8;
}

.preview-image {
  cursor: pointer;
}

.preview-image:hover {
  opacity: 0.8;
}

.hidden {
  display: none;
}
</style>
