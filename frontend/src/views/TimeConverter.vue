<template>
  <div class="time-converter">
    <div class="converter-section">
      <div class="input-group">
        <div class="label">时间戳</div>
        <div class="input-with-buttons">
          <input
            type="text"
            v-model="timestamp"
            placeholder="输入时间戳"
            @input="handleTimestampInput"
          />
          <button class="tool-btn" @click="getCurrentTimestamp">
            当前时间戳
          </button>
        </div>
      </div>

      <div class="input-group">
        <div class="label">日期时间</div>
        <div class="input-with-buttons">
          <input
            type="datetime-local"
            v-model="datetime"
            step="1"
            @input="handleDateTimeInput"
          />
          <button class="tool-btn" @click="getCurrentDateTime">当前时间</button>
        </div>
      </div>
    </div>

    <div class="format-section">
      <h3>格式化结果</h3>
      <div class="format-results">
        <div class="format-item">
          <span class="format-label">UTC：</span>
          <span class="format-value">{{ utcTime }}</span>
        </div>
        <div class="format-item">
          <span class="format-label">本地时间：</span>
          <span class="format-value">{{ localTime }}</span>
        </div>
        <div class="format-item">
          <span class="format-label">相对时间：</span>
          <span class="format-value">{{ relativeTime }}</span>
        </div>
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

dayjs.extend(utc)
dayjs.extend(relativeTimePlugin)
dayjs.locale('zh-cn')

const timestamp = ref('')
const datetime = ref('')

const handleTimestampInput = () => {
  if (!timestamp.value) {
    datetime.value = ''
    return
  }

  const ts = parseInt(timestamp.value)
  if (isNaN(ts)) return

  const milliseconds = ts.toString().length === 10 ? ts * 1000 : ts
  datetime.value = dayjs(milliseconds).format('YYYY-MM-DDTHH:mm:ss')
}

const handleDateTimeInput = () => {
  if (!datetime.value) {
    timestamp.value = ''
    return
  }

  const ts = Math.floor(new Date(datetime.value).getTime() / 1000)
  timestamp.value = ts.toString()
}

const getCurrentTimestamp = () => {
  const ts = Math.floor(Date.now() / 1000)
  timestamp.value = ts.toString()
  handleTimestampInput()
}

const getCurrentDateTime = () => {
  datetime.value = dayjs().format('YYYY-MM-DDTHH:mm:ss')
  handleDateTimeInput()
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
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.converter-section {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.input-group {
  margin-bottom: 20px;
}

.input-group:last-child {
  margin-bottom: 0;
}

.label {
  font-size: 14px;
  color: #374151;
  margin-bottom: 8px;
}

.input-with-buttons {
  display: flex;
  gap: 12px;
}

input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
}

input:focus {
  outline: none;
  border-color: #60a5fa;
  box-shadow: 0 0 0 2px rgba(96, 165, 250, 0.2);
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

.format-section {
  margin-top: 24px;
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.format-section h3 {
  margin: 0 0 16px 0;
  font-size: 16px;
  color: #111827;
}

.format-results {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.format-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.format-label {
  font-size: 14px;
  color: #6b7280;
  min-width: 80px;
}

.format-value {
  font-size: 14px;
  color: #111827;
  font-family: monospace;
}
</style>
