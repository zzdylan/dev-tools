<template>
  <div class="qrcode-tool">
    <div class="tool-section">
      <div class="tool-layout">
        <!-- 左侧：文本输入 -->
        <div class="left-panel">
          <textarea
            v-model="text"
            placeholder="输入文本内容 / 解析结果将显示在这里"
            class="text-area"
          ></textarea>
          <div class="upload-area" @drop.prevent="handleDrop" @dragover.prevent>
            <input
              type="file"
              ref="fileInput"
              accept="image/*"
              class="hidden"
              @change="handleFileSelect"
            />
            <div class="upload-content" @click="triggerFileInput">
              <div class="upload-icon">📤</div>
              <div class="upload-text">
                点击或拖拽图片到此处解析二维码
                <div class="upload-hint">支持 jpg、png、gif 格式</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧：二维码显示和控制 -->
        <div class="right-panel">
          <div class="qrcode-container">
            <canvas
              v-show="text"
              ref="canvasRef"
              width="460"
              height="460"
              class="qrcode"
            ></canvas>
            <button v-if="text" class="download-btn" @click="downloadQR">
              <span class="btn-icon">💾</span>
              <span class="btn-text">导出</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useClipboard } from '@vueuse/core'
import jsQR from 'jsqr'
import QRCode from 'qrcode'

const { copy: copyToClipboard } = useClipboard()
const text = ref('I am the qrcode generator of devTools😁')
const fileInput = ref<HTMLInputElement | null>(null)
const canvasRef = ref<HTMLCanvasElement | null>(null)

// 生成二维码的函数
const generateQRCode = async () => {
  const canvas = canvasRef.value
  if (text.value && canvas) {
    try {
      await QRCode.toCanvas(canvas, text.value, {
        width: 460,
        margin: 0,
        color: {
          dark: '#000000',
          light: '#FFFFFF',
        },
      })
    } catch (error) {
      console.error('生成二维码失败:', error)
      ElMessage.error('生成二维码失败')
    }
  }
}

// 监听文本变化，更新二维码
watch(text, generateQRCode)

// 组件挂载后生成二维码
onMounted(generateQRCode)

const clear = () => {
  text.value = ''
  ElMessage.success('已清空')
}

const downloadQR = () => {
  try {
    // 直接创建一个新的二维码图片
    const canvas = document.createElement('canvas')
    const ctx = canvas.getContext('2d')
    if (!ctx) {
      throw new Error('无法创建 canvas 上下文')
      return
    }

    // 设置 canvas 大小
    canvas.width = 460
    canvas.height = 460

    // 绘制白色背景
    ctx.fillStyle = '#FFFFFF'
    ctx.fillRect(0, 0, canvas.width, canvas.height)

    // 创建一个新的 QRCode 实例
    QRCode.toCanvas(
      canvas,
      text.value,
      {
        width: 460,
        margin: 0,
        color: {
          dark: '#000000',
          light: '#FFFFFF',
        },
      },
      (error) => {
        if (error) {
          ElMessage.error('生成二维码失败')
          return
        }

        // 转换为 PNG 并下载
        canvas.toBlob((blob) => {
          if (blob) {
            const url = URL.createObjectURL(blob)
            const link = document.createElement('a')
            link.href = url
            link.download = 'qrcode.png'
            document.body.appendChild(link)
            link.click()
            document.body.removeChild(link)
            URL.revokeObjectURL(url)
            ElMessage.success('二维码已下载')
          }
        }, 'image/png')
      }
    )
  } catch (error) {
    console.error('下载失败:', error)
    ElMessage.error('下载失败')
  }
}

const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileSelect = async (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (file) {
    await decodeQRCode(file)
  }
}

const handleDrop = async (event: DragEvent) => {
  const file = event.dataTransfer?.files[0]
  if (file) {
    await decodeQRCode(file)
  }
}

const decodeQRCode = async (file: File) => {
  try {
    const image = new Image()
    image.src = URL.createObjectURL(file)

    await new Promise((resolve, reject) => {
      image.onload = resolve
      image.onerror = reject
    })

    const canvas = document.createElement('canvas')
    const context = canvas.getContext('2d')
    if (!context) {
      throw new Error('无法创建 canvas 上下文')
    }

    canvas.width = image.width
    canvas.height = image.height
    context.drawImage(image, 0, 0)

    const imageData = context.getImageData(0, 0, canvas.width, canvas.height)
    const code = jsQR(imageData.data, imageData.width, imageData.height)

    if (code) {
      text.value = code.data
      ElMessage.success('二维码解析成功')
    } else {
      ElMessage.error('未能识别二维码')
    }

    URL.revokeObjectURL(image.src)
  } catch (error) {
    console.error('解析失败:', error)
    ElMessage.error('二维码解析失败')
  }
}
</script>

<style scoped>
.qrcode-tool {
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

.text-area {
  width: 100%;
  flex: 1;
  padding: 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  font-family: monospace;
  resize: none;
}

.right-panel {
  width: 500px;
  display: flex;
  flex-direction: column;
}

.qrcode-container {
  display: flex;
  justify-content: center;
  align-items: center;
  flex: 1;
  background: white;
  position: relative;
}

.qrcode {
  width: 460px;
  height: 460px;
  background: #fff;
}

.qrcode-footer {
  display: flex;
  flex-direction: column;
  padding: 12px;
  border-top: 1px solid #e5e7eb;
}

.control-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.control-item label {
  font-size: 14px;
  color: #374151;
}

.select-control {
  padding: 6px;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-size: 14px;
}

.range-control {
  width: 100%;
}

.color-controls {
  display: flex;
  gap: 16px;
}

.color-input {
  display: flex;
  gap: 8px;
}

.color-text {
  width: 80px;
  padding: 4px 8px;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-size: 14px;
}

.color-picker {
  width: 30px;
  height: 30px;
  padding: 0;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  cursor: pointer;
}

.upload-area {
  border: 2px dashed #d1d5db;
  border-radius: 8px;
  padding: 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
}

.upload-area:hover {
  border-color: #60a5fa;
  background: #f8fafc;
}

.hidden {
  display: none;
}

.upload-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.upload-icon {
  font-size: 32px;
}

.upload-text {
  color: #374151;
  font-size: 14px;
}

.upload-hint {
  color: #6b7280;
  font-size: 12px;
  margin-top: 4px;
}

.download-btn {
  position: absolute;
  right: 20px;
  top: 20px;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  background: #f9fafb;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-text {
  font-size: 14px;
}

.download-btn:hover {
  background: #e5e7eb;
  border-color: #9ca3af;
}
</style>
