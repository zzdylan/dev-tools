<template>
  <div class="color-converter">
    <!-- 顶部：功能按钮 -->
    <div class="top-header">
      <div class="tab-nav">
        <button class="action-btn" @click="loadExample">示例</button>
      </div>
      <div class="tab-actions">
        <button class="clear-btn" @click="clearAll" title="清空">× 清空</button>
      </div>
    </div>

    <!-- 底部：内容区域 -->
    <div class="content-layout">
      <!-- 左侧：输入区域 -->
      <div class="input-panel">
        <div class="input-section">
          <div class="input-label">十六进制颜色</div>
          <input
            v-model="hexColor"
            placeholder="#FF5733"
            class="color-input"
            @input="convertFromHex"
          />
        </div>
        <div class="input-section">
          <div class="input-label">RGB颜色</div>
          <div class="rgb-inputs">
            <input
              v-model.number="rgbColor.r"
              type="number"
              min="0"
              max="255"
              placeholder="R"
              class="rgb-input"
              @input="convertFromRgb"
            />
            <input
              v-model.number="rgbColor.g"
              type="number"
              min="0"
              max="255"
              placeholder="G"
              class="rgb-input"
              @input="convertFromRgb"
            />
            <input
              v-model.number="rgbColor.b"
              type="number"
              min="0"
              max="255"
              placeholder="B"
              class="rgb-input"
              @input="convertFromRgb"
            />
          </div>
        </div>
      </div>

      <!-- 右侧：颜色预览 -->
      <div class="output-panel">
        <div class="color-preview" :style="{ backgroundColor: previewColor }">
          <div class="preview-text" :style="{ color: getContrastColor(previewColor) }">
            颜色预览
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useToolsStore } from '../stores/tools'
import { storeToRefs } from 'pinia'

const store = useToolsStore()
const { colorConverter } = storeToRefs(store)

interface RGBColor {
  r: number
  g: number
  b: number
}

interface HSLColor {
  h: number
  s: number
  l: number
}

const hexColor = computed({
  get: () => colorConverter.value.hexColor || '#FF5733',
  set: (val) => (colorConverter.value.hexColor = val),
})

const rgbColor = computed({
  get: () => colorConverter.value.rgbColor || { r: 255, g: 87, b: 51 },
  set: (val) => (colorConverter.value.rgbColor = val),
})


// 颜色预览
const previewColor = computed(() => {
  const { r, g, b } = rgbColor.value
  return `rgb(${r || 0}, ${g || 0}, ${b || 0})`
})

// 输出格式化
const outputHex = computed(() => {
  const { r, g, b } = rgbColor.value
  const toHex = (n: number) => Math.max(0, Math.min(255, n || 0)).toString(16).padStart(2, '0').toUpperCase()
  return `#${toHex(r)}${toHex(g)}${toHex(b)}`
})

const outputRgb = computed(() => {
  const { r, g, b } = rgbColor.value
  return `rgb(${r || 0}, ${g || 0}, ${b || 0})`
})



// 从十六进制转换
const convertFromHex = () => {
  const hex = hexColor.value.replace('#', '')
  if (hex.length === 6 || hex.length === 3) {
    const fullHex = hex.length === 3 ? hex.split('').map((c: string) => c + c).join('') : hex
    const r = parseInt(fullHex.substring(0, 2), 16)
    const g = parseInt(fullHex.substring(2, 4), 16)
    const b = parseInt(fullHex.substring(4, 6), 16)
    
    if (!isNaN(r) && !isNaN(g) && !isNaN(b)) {
      rgbColor.value = { r, g, b }
    }
  }
}

// 从RGB转换
const convertFromRgb = () => {
  const { r, g, b } = rgbColor.value
  if (r >= 0 && r <= 255 && g >= 0 && g <= 255 && b >= 0 && b <= 255) {
    hexColor.value = `#${r.toString(16).padStart(2, '0')}${g.toString(16).padStart(2, '0')}${b.toString(16).padStart(2, '0')}`.toUpperCase()
  }
}


// 获取对比色
const getContrastColor = (bgColor: string) => {
  const { r, g, b } = rgbColor.value
  const brightness = (r * 299 + g * 587 + b * 114) / 1000
  return brightness > 128 ? '#000000' : '#ffffff'
}

// 加载示例
const loadExample = () => {
  hexColor.value = '#FF5733'
  convertFromHex()
}

// 清空所有内容
const clearAll = () => {
  hexColor.value = '#000000'
  rgbColor.value = { r: 0, g: 0, b: 0 }
}

// 初始化时转换一次
convertFromHex()
</script>

<style scoped>
.color-converter {
  height: 100%;
  width: 100%;
  margin: 0 !important;
  padding: 0 !important;
  background: white;
  display: flex;
  flex-direction: column;
}

.top-header {
  display: flex;
  justify-content: space-between;
  align-items: stretch;
  padding: 0;
  margin: 0;
  background: #f8f9fa;
  height: 28px;
  flex-shrink: 0;
}

.tab-nav {
  display: flex;
  align-items: stretch;
  background: #f8f9fa;
}

.action-btn {
  padding: 0 10px;
  margin: 0;
  background: #f8f9fa;
  border-top: none;
  border-bottom: none;
  border-left: 1px solid #d1d5db;
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
  box-sizing: border-box;
}

.action-btn:hover {
  background: #e9ecef;
}

.tab-actions {
  display: flex;
  align-items: stretch;
  background: #f8f9fa;
}

.clear-btn {
  padding: 0 10px;
  margin: 0;
  background: #f8f9fa;
  border-top: none;
  border-bottom: none;
  border-left: 1px solid #d1d5db;
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
  box-sizing: border-box;
}

.clear-btn:hover {
  background: #e9ecef;
}

.content-layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  flex: 1;
  align-items: stretch;
}

.input-panel,
.output-panel {
  border: 1px solid #d1d5db;
  background: #ffffff;
  padding: 16px;
}

.input-section {
  margin-bottom: 16px;
}

.input-section:last-child {
  margin-bottom: 0;
}

.input-label {
  font-size: 12px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 6px;
}

.color-input {
  width: 100%;
  padding: 6px 10px;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 11px;
  color: #212529;
  background: #ffffff;
  outline: none;
  box-sizing: border-box;
}

.color-input:focus {
  border-color: #3b82f6;
  box-shadow: 0 0 0 1px #3b82f6;
}

.color-input::placeholder {
  color: #9ca3af;
}

.rgb-inputs {
  display: flex;
  gap: 8px;
}

.rgb-input {
  flex: 1;
  padding: 6px 8px;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 11px;
  color: #212529;
  background: #ffffff;
  outline: none;
  text-align: center;
  box-sizing: border-box;
}

.rgb-input:focus {
  border-color: #3b82f6;
  box-shadow: 0 0 0 1px #3b82f6;
}

.rgb-input::placeholder {
  color: #9ca3af;
}

.color-preview {
  height: 120px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.preview-text {
  font-size: 14px;
  font-weight: 500;
  text-shadow: 0 0 4px rgba(0, 0, 0, 0.3);
}


@media (max-width: 1024px) {
  .content-layout {
    grid-template-columns: 1fr;
    grid-template-rows: 200px 150px;
    gap: 12px;
  }
  
  .input-panel {
    padding: 12px;
    overflow-y: auto;
  }
  
  .input-section {
    margin-bottom: 12px;
  }
  
  .output-panel {
    padding: 12px;
  }
}
</style>