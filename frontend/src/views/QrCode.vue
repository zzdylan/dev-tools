<template>
  <div class="qrcode-tool">
    <div class="tool-container">
      <!-- 左侧：文本输入 -->
      <div class="input-section">
        <div class="section-header">
          <h3>文本内容</h3>
          <button v-if="text" class="btn-clear" @click="clear">清空</button>
        </div>
        <textarea
          v-model="text"
          placeholder="输入要生成二维码的文本内容..."
          class="text-input"
        ></textarea>
      </div>

      <!-- 右侧：二维码显示 -->
      <div class="qrcode-section">
        <div class="section-header">
          <h3>二维码</h3>
          <button v-if="text" class="btn-download" @click="downloadQR">
            保存图片
          </button>
        </div>
        <div 
          class="qrcode-display" 
          @drop.prevent="handleDrop" 
          @dragover.prevent="handleDragOver"
          @dragenter.prevent="handleDragEnter"
          @dragleave.prevent="handleDragLeave"
          :class="{ 'drag-over': isDragOver }"
        >
          <input
            type="file"
            ref="fileInput"
            accept="image/*"
            class="hidden"
            @change="handleFileSelect"
          />
          <div v-if="text" class="qrcode-canvas-wrapper" @click="triggerFileInput">
            <canvas
              ref="canvasRef"
              width="300"
              height="300"
              class="qrcode-canvas"
            ></canvas>
            <div class="qrcode-overlay">
              <div class="overlay-content">
                <div class="overlay-icon">📤</div>
                <div class="overlay-text">点击上传图片解析二维码</div>
              </div>
            </div>
          </div>
          <div v-else class="qrcode-placeholder" @click="triggerFileInput">
            <div class="placeholder-icon">📱</div>
            <div class="placeholder-text">输入文本生成二维码</div>
            <div class="placeholder-hint">或点击/拖拽图片解析二维码</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, computed, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { useClipboard } from '@vueuse/core'
import jsQR from 'jsqr'
import QRCode from 'qrcode'
import { useToolsStore } from '../stores/tools'
import { storeToRefs } from 'pinia'

const { copy } = useClipboard()
const store = useToolsStore()
const { qrCode } = storeToRefs(store)
const text = computed({
  get: () => qrCode.value.text,
  set: (val) => (qrCode.value.text = val),
})
const fileInput = ref<HTMLInputElement | null>(null)
const canvasRef = ref<HTMLCanvasElement | null>(null)
const isDragOver = ref(false)

// 生成二维码的函数
const generateQRCode = async () => {
  const canvas = canvasRef.value
  console.log('generateQRCode called, text:', text.value, 'canvas:', canvas)
  if (text.value && canvas) {
    try {
      await QRCode.toCanvas(canvas, text.value, {
        width: 300,
        margin: 1,
        color: {
          dark: '#000000',
          light: '#FFFFFF',
        },
      })
      console.log('二维码生成成功')
    } catch (error) {
      console.error('生成二维码失败:', error)
      ElMessage.error('生成二维码失败')
    }
  } else {
    console.log('无法生成二维码 - text:', text.value, 'canvas:', canvas)
  }
}

// 监听文本变化，更新二维码
watch(text, generateQRCode)

// 组件挂载后生成二维码
onMounted(generateQRCode)

const clear = () => {
  text.value = ''
  // ElMessage.success('已清空')
}

const downloadQR = async () => {
  try {
    // 直接创建一个新的二维码图片
    const canvas = document.createElement('canvas')
    const ctx = canvas.getContext('2d')
    if (!ctx) {
      throw new Error('无法创建 canvas 上下文')
      return
    }

    // 设置 canvas 大小
    canvas.width = 300
    canvas.height = 300

    // 绘制白色背景
    ctx.fillStyle = '#FFFFFF'
    ctx.fillRect(0, 0, canvas.width, canvas.height)

    // 创建一个新的 QRCode 实例
    QRCode.toCanvas(
      canvas,
      text.value,
      {
        width: 300,
        margin: 1,
        color: {
          dark: '#000000',
          light: '#FFFFFF',
        },
      },
      async (error) => {
        if (error) {
          ElMessage.error('生成二维码失败')
          return
        }

        // 转换为 PNG 并下载
        canvas.toBlob(async (blob) => {
          if (blob) {
            // 尝试使用现代的文件保存API
            if ('showSaveFilePicker' in window) {
              try {
                const fileHandle = await (window as any).showSaveFilePicker({
                  suggestedName: 'qrcode.png',
                  types: [{
                    description: 'PNG图片',
                    accept: { 'image/png': ['.png'] }
                  }]
                })
                const writable = await fileHandle.createWritable()
                await writable.write(blob)
                await writable.close()
                ElMessage.success('二维码已保存')
              } catch (err) {
                // 用户取消了保存
                if ((err as Error).name !== 'AbortError') {
                  // 降级到传统下载方式
                  fallbackDownload(blob)
                }
              }
            } else {
              // 降级到传统下载方式
              fallbackDownload(blob)
            }
          }
        }, 'image/png')
      }
    )
  } catch (error) {
    console.error('下载失败:', error)
    ElMessage.error('下载失败')
  }
}

// 传统下载方式
const fallbackDownload = (blob: Blob) => {
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = 'qrcode.png'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
  ElMessage.success('二维码已下载到默认下载目录')
}

// 触发文件选择
const triggerFileInput = () => {
  fileInput.value?.click()
}

// 处理文件选择
const handleFileSelect = (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (file) {
    decodeQRCode(file)
  }
}

// 处理拖拽进入
const handleDragEnter = (event: DragEvent) => {
  event.preventDefault()
  isDragOver.value = true
}

// 处理拖拽离开
const handleDragLeave = (event: DragEvent) => {
  event.preventDefault()
  isDragOver.value = false
}

// 处理拖拽悬停
const handleDragOver = (event: DragEvent) => {
  event.preventDefault()
  isDragOver.value = true
}

// 处理拖拽放下
const handleDrop = (event: DragEvent) => {
  event.preventDefault()
  isDragOver.value = false
  
  const file = event.dataTransfer?.files?.[0]
  if (file && file.type.startsWith('image/')) {
    decodeQRCode(file)
  } else {
    ElMessage.error('请上传图片文件')
  }
}

// 解析二维码图片
const decodeQRCode = (file: File) => {
  if (!file) {
    ElMessage.error('未选择文件')
    return
  }

  if (!file.type.startsWith('image/')) {
    ElMessage.error('请选择图片文件')
    return
  }


  const reader = new FileReader()
  reader.onload = (e) => {
    const img = new Image()
    img.onload = async () => {
      try {
        const canvas = document.createElement('canvas')
        const ctx = canvas.getContext('2d')
        if (!ctx) {
          ElMessage.error('无法创建画布')
          return
        }

        canvas.width = img.width
        canvas.height = img.height
        ctx.drawImage(img, 0, 0)

        const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height)
        console.log('图片尺寸:', canvas.width, 'x', canvas.height)
        console.log('图片数据:', imageData.data.length, '字节')
        
        const code = jsQR(imageData.data, imageData.width, imageData.height)
        console.log('解析结果:', code)

        if (code) {
          text.value = code.data
          console.log('解析出的文本:', code.data)
          // 等待DOM更新后再生成二维码
          await nextTick()
          await generateQRCode()
        } else {
          ElMessage.error('未能识别二维码，请确保图片清晰且包含二维码')
          console.log('未能识别二维码')
        }
      } catch (error) {
        console.error('解析二维码时出错:', error)
        ElMessage.error('解析二维码时出错')
      }
    }
    img.onerror = () => {
      ElMessage.error('图片加载失败')
    }
    img.src = e.target?.result as string
  }
  reader.onerror = () => {
    ElMessage.error('文件读取失败')
  }
  reader.readAsDataURL(file)
}

const copyText = async () => {
  try {
    await copy(text.value)
    ElMessage.success('文本已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}
</script>

<style scoped>
.qrcode-tool {
  padding: 16px;
  max-width: 1200px;
  margin: 0 auto;
  background: #f8fafc;
  min-height: 100vh;
}

.tool-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  height: calc(100vh - 60px);
}

.input-section,
.qrcode-section {
  background: white;
  border-radius: 8px;
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
  padding: 12px 16px;
  border-bottom: 1px solid #e2e8f0;
  background: #f8fafc;
}

.section-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #1e293b;
}

.btn-clear {
  padding: 4px 8px;
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fecaca;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-clear:hover {
  background: #fee2e2;
}

.btn-download {
  padding: 4px 8px;
  background: #f1f5f9;
  color: #475569;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-download:hover {
  background: #e2e8f0;
}

.text-input {
  flex: 1;
  padding: 14px;
  border: none;
  outline: none;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.4;
  resize: none;
  background: transparent;
  color: #1e293b;
}

.text-input::placeholder {
  color: #94a3b8;
}

.qrcode-display {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 14px;
  transition: all 0.3s ease;
}

.qrcode-display.drag-over {
  background: #f0f9ff;
  border: 2px dashed #3b82f6;
  border-radius: 6px;
}

.qrcode-canvas-wrapper {
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
}

.qrcode-canvas {
  width: 250px;
  height: 250px;
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.2s;
}

.qrcode-canvas:hover {
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.qrcode-hint {
  margin-top: 8px;
  font-size: 11px;
  color: #64748b;
  text-align: center;
}

.qrcode-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: #94a3b8;
  text-align: center;
}

.placeholder-icon {
  font-size: 36px;
  opacity: 0.6;
}

.placeholder-text {
  font-size: 14px;
  font-weight: 500;
}

.placeholder-hint {
  font-size: 12px;
  color: #64748b;
  margin-top: 6px;
}

.hidden {
  display: none;
}

.qrcode-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  opacity: 0;
  transition: opacity 0.2s;
}

.qrcode-canvas-wrapper:hover .qrcode-overlay {
  opacity: 1;
}

.overlay-content {
  text-align: center;
  color: white;
}

.overlay-icon {
  font-size: 24px;
  margin-bottom: 6px;
}

.overlay-text {
  font-size: 12px;
  font-weight: 500;
}

.qrcode-placeholder {
  cursor: pointer;
  transition: all 0.2s;
}

.qrcode-placeholder:hover {
  background: #f1f5f9;
  border-radius: 6px;
}

@media (max-width: 1024px) {
  .tool-container {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .qrcode-tool {
    padding: 16px;
  }
  
  .qrcode-canvas {
    width: 250px;
    height: 250px;
  }
}
</style>
