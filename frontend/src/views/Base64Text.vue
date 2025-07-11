<template>
  <div class="base64-text">
    <!-- 顶部：标签页按钮 -->
    <div class="top-header">
      <div class="tab-nav">
        <button 
          class="tab-btn"
          :class="{ active: activeTab === 'encode' }"
          @click="activeTab = 'encode'"
        >
          编码
        </button>
        <button 
          class="tab-btn"
          :class="{ active: activeTab === 'decode' }"
          @click="activeTab = 'decode'"
        >
          解码
        </button>
      </div>
      <div class="tab-actions">
        <button class="clear-btn" @click="clearAll" title="清空">× 清空</button>
      </div>
    </div>

    <!-- 底部：内容区域 -->
    <div class="content-layout">
      <!-- 左侧：输入文本区域 -->
      <div class="input-panel">
        <textarea
          v-model="inputText"
          :placeholder="getInputPlaceholder()"
          class="text-area"
          autocomplete="off"
          spellcheck="false"
          @input="autoConvert"
        ></textarea>
      </div>

      <!-- 右侧：输出文本区域 -->
      <div class="output-panel">
        <textarea
          v-model="outputText"
          :placeholder="getOutputPlaceholder()"
          class="text-area"
          readonly
          autocomplete="off"
          spellcheck="false"
        ></textarea>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useToolsStore } from '../stores/tools'
import { storeToRefs } from 'pinia'

const store = useToolsStore()
const { base64Text } = storeToRefs(store)

const activeTab = computed({
  get: () => base64Text.value.activeTab || 'encode',
  set: (val) => (base64Text.value.activeTab = val),
})

const inputText = computed({
  get: () => base64Text.value.inputText || '',
  set: (val) => (base64Text.value.inputText = val),
})

const outputText = computed({
  get: () => base64Text.value.outputText || '',
  set: (val) => (base64Text.value.outputText = val),
})

// 自动转换
const autoConvert = () => {
  if (!inputText.value.trim()) {
    outputText.value = ''
    return
  }
  
  if (activeTab.value === 'encode') {
    try {
      outputText.value = btoa(inputText.value)
    } catch (error) {
      outputText.value = '编码失败'
    }
  } else {
    try {
      outputText.value = atob(inputText.value)
    } catch (error) {
      outputText.value = '解码失败：无效的 Base64 字符串'
    }
  }
}

// 切换模式时重新转换
const switchMode = () => {
  autoConvert()
}

// 清空所有内容
const clearAll = () => {
  inputText.value = ''
  outputText.value = ''
}

// 获取输入占位符
const getInputPlaceholder = () => {
  return activeTab.value === 'encode' ? '输入需要编码的文本...' : '输入需要解码的Base64字符串...'
}

// 获取输出占位符
const getOutputPlaceholder = () => {
  return activeTab.value === 'encode' ? 'Base64编码结果...' : '解码后的文本...'
}


// 监听标签页切换
watch(activeTab, switchMode)


</script>

<style scoped>
.base64-text {
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
  background: #ffffff;
  height: 28px;
  flex-shrink: 0;
}

.tab-nav {
  display: flex;
  align-items: stretch;
  background: #f8f9fa;
}

.tab-btn {
  padding: 0 10px;
  margin: 0;
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
  box-sizing: border-box;
}

.tab-btn:first-child {
  border-left: 1px solid #d1d5db;
}

.tab-btn + .tab-btn {
  border-left: 1px solid #d1d5db;
}

.tab-btn:last-child {
  border-right: 1px solid #d1d5db;
}

.tab-btn:hover {
  background: #e9ecef;
}

.tab-btn.active {
  background: #ffffff;
  color: #212529;
  font-weight: 500;
}

.tab-actions {
  display: flex;
  align-items: center;
  background: #f8f9fa;
}

.clear-btn {
  padding: 0 10px;
  margin: 0;
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
  
  .base64-text {
    padding: 12px;
  }
}
</style>
