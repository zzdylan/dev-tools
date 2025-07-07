<template>
  <div class="time-converter">
    <!-- 顶部：操作按钮 -->
    <div class="top-header">
      <div class="tab-actions">
        <button class="current-time-btn" @click="getCurrentTime" title="当前时间">当前时间</button>
      </div>
      <div class="tab-nav">
      </div>
    </div>

    <!-- 底部：内容区域 -->
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
          type="datetime-local"
          v-model="datetime"
          step="1"
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
  datetime.value = dayjs(milliseconds).format('YYYY-MM-DDTHH:mm:ss')
}

const convertToDateTime = () => {
  if (!timestamp.value.trim()) {
    datetime.value = ''
    return
  }

  const ts = parseInt(timestamp.value)
  if (isNaN(ts)) return

  const milliseconds = ts.toString().length === 10 ? ts * 1000 : ts
  datetime.value = dayjs(milliseconds).format('YYYY-MM-DDTHH:mm:ss')
}

const convertToTimestamp = () => {
  if (!datetime.value) {
    timestamp.value = ''
    return
  }

  const ts = Math.floor(new Date(datetime.value).getTime() / 1000)
  timestamp.value = ts.toString()
}

const clearTimestamp = () => {
  timestamp.value = ''
  datetime.value = ''
}

const clearDateTime = () => {
  datetime.value = ''
  timestamp.value = ''
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
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #f8fafc;
}

.top-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #ffffff;
  height: 35px;
  flex-shrink: 0;
}

.tab-nav {
  display: flex;
  height: 100%;
}

.tab-actions {
  display: flex;
  height: 100%;
}

.current-time-btn {
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
  min-width: 60px;
  height: 100%;
}

.current-time-btn:hover {
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
