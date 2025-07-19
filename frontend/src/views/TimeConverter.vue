<template>
  <div class="time-converter">
    <!-- 顶部导航栏 -->
    <div class="top-header">
      <div class="tab-nav">
        <button class="action-btn" @click="getCurrentTime" title="当前时间">当前时间</button>
      </div>
      <div class="tab-actions">
        <button class="clear-btn" @click="clearAll" title="清空">× 清空</button>
      </div>
    </div>

    <!-- 内容区域 -->
    <div class="converter-container">
      <!-- 时间戳输入 -->
      <div class="input-row">
        <label class="input-label">时间戳</label>
        <input
          type="text"
          v-model="timestamp"
          placeholder="输入时间戳 (秒或毫秒)..."
          class="time-input"
          @input="convertToDateTime"
        />
      </div>

      <!-- 日期时间输入 -->
      <div class="input-row">
        <label class="input-label">日期时间</label>
        <input
          type="text"
          v-model="datetime"
          placeholder="2025-01-15 14:30:00"
          class="time-input"
          @input="convertToTimestamp"
        />
      </div>

      <!-- UTC 时间 -->
      <div class="input-row">
        <label class="input-label">UTC 时间</label>
        <input
          type="text"
          :value="utcTime"
          readonly
          class="time-input readonly"
        />
      </div>

      <!-- 本地时间 -->
      <div class="input-row">
        <label class="input-label">本地时间</label>
        <input
          type="text"
          :value="localTime"
          readonly
          class="time-input readonly"
        />
      </div>

      <!-- 相对时间 -->
      <div class="input-row">
        <label class="input-label">相对时间</label>
        <input
          type="text"
          :value="relativeTime"
          readonly
          class="time-input readonly"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import dayjs from 'dayjs'
import utc from 'dayjs/plugin/utc'
import relativeTimePlugin from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'
import { useToolsStore } from '../stores/tools'
import { storeToRefs } from 'pinia'

dayjs.extend(utc)
dayjs.extend(relativeTimePlugin)
dayjs.locale('zh-cn')

const store = useToolsStore()
const { timeConverter } = storeToRefs(store)
const timestamp = computed({
  get: () => timeConverter.value.timestamp,
  set: (val) => (timeConverter.value.timestamp = val),
})
const datetime = computed({
  get: () => timeConverter.value.datetime,
  set: (val) => (timeConverter.value.datetime = val),
})

const getCurrentTime = () => {
  const ts = Math.floor(Date.now() / 1000)
  timestamp.value = ts.toString()
  const milliseconds = ts * 1000
  datetime.value = dayjs(milliseconds).format('YYYY-MM-DD HH:mm:ss')
}

const convertToDateTime = () => {
  if (!timestamp.value.trim()) {
    datetime.value = ''
    return
  }

  const ts = parseInt(timestamp.value)
  if (isNaN(ts)) return

  const milliseconds = ts.toString().length === 10 ? ts * 1000 : ts
  datetime.value = dayjs(milliseconds).format('YYYY-MM-DD HH:mm:ss')
}

const convertToTimestamp = () => {
  if (!datetime.value.trim()) {
    timestamp.value = ''
    return
  }

  // 尝试解析各种日期格式
  let date = dayjs(datetime.value)
  
  // 如果 dayjs 无法解析，尝试使用 Date 构造函数
  if (!date.isValid()) {
    date = dayjs(new Date(datetime.value))
  }
  
  if (!date.isValid()) {
    return
  }

  const ts = Math.floor(date.valueOf() / 1000)
  timestamp.value = ts.toString()
}

const clearAll = () => {
  timestamp.value = ''
  datetime.value = ''
}


const utcTime = computed(() => {
  if (!timestamp.value) return '-'
  const ts = parseInt(timestamp.value)
  if (isNaN(ts)) return '-'
  const milliseconds = ts.toString().length === 10 ? ts * 1000 : ts
  return dayjs(milliseconds).utc().format('YYYY-MM-DD HH:mm:ss UTC')
})

const localTime = computed(() => {
  if (!timestamp.value) return '-'
  const ts = parseInt(timestamp.value)
  if (isNaN(ts)) return '-'
  const milliseconds = ts.toString().length === 10 ? ts * 1000 : ts
  return dayjs(milliseconds).format('YYYY-MM-DD HH:mm:ss')
})

const relativeTime = computed(() => {
  if (!timestamp.value) return '-'
  const ts = parseInt(timestamp.value)
  if (isNaN(ts)) return '-'
  const milliseconds = ts.toString().length === 10 ? ts * 1000 : ts
  return dayjs(milliseconds).fromNow()
})
</script>

<style scoped>
.time-converter {
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

.tab-actions {
  display: flex;
  align-items: stretch;
  background: #f8f9fa;
}

.action-btn {
  padding: 0 10px;
  margin: 0;
  background: #f8f9fa;
  border: none;
  border-left: 1px solid #d1d5db;
  border-right: 1px solid #d1d5db;
  border-top: none;
  border-bottom: none;
  font-size: 10px;
  color: #6c757d;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 60px;
  height: 100%;
  box-sizing: border-box;
}

.action-btn:hover {
  background: #e9ecef;
}

.clear-btn {
  padding: 0 10px;
  margin: 0;
  background: #f8f9fa;
  border: none;
  border-left: 1px solid #d1d5db;
  border-right: 1px solid #d1d5db;
  border-top: none;
  border-bottom: none;
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

.converter-container {
  flex: 1;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  background: white;
}


.input-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
}

.input-label {
  min-width: 80px;
  width: 80px;
  font-size: 12px;
  color: #6c757d;
  font-weight: 500;
  flex-shrink: 0;
}


.time-input {
  flex: 1;
  padding: 6px 8px;
  border: none;
  outline: none;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.4;
  background: #f8f9fa;
  color: #1e293b;
  height: 28px;
  box-sizing: border-box;
}

.time-input.readonly {
  background: #f1f5f9;
  color: #64748b;
  cursor: default;
}

.time-input::placeholder {
  color: #94a3b8;
}


@media (max-width: 768px) {
  .input-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
  
  .input-label {
    min-width: auto;
  }
  
  .time-input {
    width: 100%;
  }
}
</style>
