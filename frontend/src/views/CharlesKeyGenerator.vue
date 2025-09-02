<template>
  <div class="charles-generator">
    <!-- 顶部：功能按钮 -->
    <div class="top-header">
      <div class="tab-nav">
        <button class="action-btn" @click="loadExample">
          示例
        </button>
      </div>
      <div class="tab-actions">
        <button class="clear-btn" @click="clearAll" title="清空">× 清空</button>
      </div>
    </div>

    <!-- 底部：内容区域 -->
    <div class="content-layout">
      <!-- 左侧：输入用户名区域 -->
      <div class="input-panel">
        <textarea
          v-model="userName"
          placeholder="输入用户名..."
          class="text-area"
          autocomplete="off"
          spellcheck="false"
          @input="generateKey"
        ></textarea>
      </div>

      <!-- 右侧：输出激活码区域 -->
      <div class="output-panel">
        <textarea
          v-model="licenseKey"
          placeholder="生成的Charles激活码..."
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
import { computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useToolsStore } from '../stores/tools'
import { storeToRefs } from 'pinia'
import { GenerateCharlesKey } from '../../wailsjs/go/main/App'

const store = useToolsStore()
const { charlesGenerator } = storeToRefs(store)

const userName = computed({
  get: () => charlesGenerator.value.userName || '',
  set: (val) => (charlesGenerator.value.userName = val),
})

const licenseKey = computed({
  get: () => charlesGenerator.value.licenseKey || '',
  set: (val) => (charlesGenerator.value.licenseKey = val),
})

// 生成激活码
const generateKey = async () => {
  if (!userName.value.trim()) {
    licenseKey.value = ''
    return
  }
  
  try {
    const key = await GenerateCharlesKey(userName.value.trim())
    licenseKey.value = key
  } catch (error) {
    licenseKey.value = '生成失败'
    ElMessage.error('生成激活码失败: ' + error)
  }
}

// 加载示例
const loadExample = () => {
  userName.value = 'charles'
  generateKey()
}

// 清空所有内容
const clearAll = () => {
  userName.value = ''
  licenseKey.value = ''
}

// 组件挂载时，如果有内容自动生成
onMounted(() => {
  if (userName.value.trim()) {
    generateKey()
  }
})
</script>

<style scoped>
.charles-generator {
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
  grid-template-columns: 1fr;
  gap: 1px;
  align-items: stretch;
}

.input-panel,
.output-panel {
  border: 1px solid #d1d5db;
  background: #ffffff;
  height: 120px;
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

</style>