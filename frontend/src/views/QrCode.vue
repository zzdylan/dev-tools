<template>
  <div class="qrcode-tool">
    <div class="tool-section">
      <!-- 顶部导航栏 -->
      <div class="top-header">
        <div class="tab-nav">
          <button class="action-btn" @click="loadExample">示例</button>
        </div>
        <div class="tab-actions">
          <button v-if="text" class="download-btn" @click="downloadQR">保存图片</button>
          <button class="clear-btn" @click="clear" title="清空">× 清空</button>
        </div>
      </div>

      <!-- 内容区域 -->
      <div class="content-layout">
        <!-- 左侧：文本输入 -->
        <div class="input-panel">
          <textarea
            v-model="text"
            placeholder="输入要生成二维码的文本内容..."
            class="text-area"
            autocomplete="off"
            autocorrect="off"
            autocapitalize="off"
            spellcheck="false"
          ></textarea>
        </div>

        <!-- 右侧：二维码显示 -->
        <div class="output-panel">
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
import { saveImage } from '../utils/fileUtils'

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
}

const loadExample = async () => {
  text.value = 'https://www.baidu.com'
  await nextTick()
  generateQRCode()
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

        try {
          // 获取canvas的base64数据
          const dataURL = canvas.toDataURL('image/png')
          const base64Data = dataURL.split(',')[1]
          
          // 使用通用保存工具
          const savedPath = await saveImage(base64Data, 'qrcode.png', '保存二维码图片')
          ElMessage.success(`二维码已保存到: ${savedPath}`)
        } catch (saveError) {
          const errorMsg = saveError?.toString() || ''
          if (errorMsg.includes('用户取消保存')) {
            // 用户取消，不显示错误提示
          } else {
            ElMessage.error(`保存失败: ${saveError}`)
          }
        }
      }
    )
  } catch (error) {
    console.error('下载失败:', error)
    ElMessage.error('下载失败')
  }
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
  height: 100%;
}

.tool-section {
  background: white;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  height: 100%;
}

.top-header {
  display: flex;
  justify-content: space-between;
  align-items: stretch;
  padding: 0;
  background: #ffffff;
  height: 28px;
  margin-bottom: 8px;
}

.tab-nav {
  display: flex;
  align-items: stretch;
  border: 1px solid #d1d5db;
}

.action-btn {
  padding: 0 10px;
  background: #f8f9fa;
  border: none;
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

.action-btn:hover {
  background: #e9ecef;
}

.tab-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 12px;
  background: #ffffff;
}

.clear-btn,
.download-btn {
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

.clear-btn:hover,
.download-btn:hover {
  background: #e9ecef;
}


.content-layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
  height: calc(100% - 40px);
  align-items: stretch;
}

.input-panel,
.output-panel {
  border: 1px solid #d1d5db;
  background: #ffffff;
  display: flex;
  flex-direction: column;
}

.text-area {
  flex: 1;
  padding: 12px;
  border: none;
  outline: none;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 11px;
  line-height: 1.4;
  resize: none;
  background: transparent;
  color: #212529;
}

.text-area::placeholder {
  color: #9ca3af;
}

.qrcode-display {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 16px;
  transition: all 0.3s ease;
  position: relative;
}

.qrcode-display.drag-over {
  background: #f0f9ff;
  border: 2px dashed #3b82f6;
}

.qrcode-canvas-wrapper {
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  cursor: pointer;
}

.qrcode-canvas {
  width: 200px;
  height: 200px;
  border: 1px solid #d1d5db;
  cursor: pointer;
  transition: all 0.2s;
}

.qrcode-canvas:hover {
  opacity: 0.8;
}



.qrcode-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: #9ca3af;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
  padding: 40px;
}

.qrcode-placeholder:hover {
  background: #f8f9fa;
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
  color: #6b7280;
  margin-top: 6px;
}

.hidden {
  display: none;
}

@media (max-width: 768px) {
  .content-layout {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .qrcode-canvas {
    width: 250px;
    height: 250px;
  }
}
</style>
