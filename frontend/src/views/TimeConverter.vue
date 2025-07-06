<template>
  <div class="time-converter">
    <div class="converter-container">
      <!-- 时间戳输入区域 -->
      <div class="input-section">
        <div class="section-header">
          <h3>时间戳</h3>
          <button v-if="timestamp" class="btn-clear" @click="clearTimestamp">清空</button>
        </div>
        <div class="input-wrapper">
          <input
            type="text"
            v-model="timestamp"
            placeholder="输入时间戳 (秒或毫秒)..."
            class="time-input"
            @input="convertToDateTime"
          />
        </div>
      </div>

      <!-- 转换按钮 -->
      <div class="controls-section">
        <button class="btn-now" @click="getCurrentTime">当前时间</button>
      </div>

      <!-- 日期时间输入区域 -->
      <div class="output-section">
        <div class="section-header">
          <h3>日期时间</h3>
          <button v-if="datetime" class="btn-clear" @click="clearDateTime">清空</button>
        </div>
        <div class="output-wrapper">
          <input
            type="datetime-local"
            v-model="datetime"
            step="1"
            class="time-input"
            @input="convertToTimestamp"
          />
        </div>
      </div>
    </div>

    <!-- 格式化结果区域 -->
    <div class="results-section" v-if="timestamp">
      <div class="section-header">
        <h3>格式化结果</h3>
      </div>
      <div class="results-content">
        <div class="result-card">
          <div class="result-label">UTC 时间</div>
          <div class="result-value">{{ utcTime }}</div>
        </div>
        <div class="result-card">
          <div class="result-label">本地时间</div>
          <div class="result-value">{{ localTime }}</div>
        </div>
        <div class="result-card">
          <div class="result-label">相对时间</div>
          <div class="result-value">{{ relativeTime }}</div>
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
  padding: 16px;
  max-width: 1000px;
  margin: 0 auto;
  background: #f8fafc;
  min-height: 100vh;
}

.converter-container {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  gap: 16px;
  margin-bottom: 16px;
}

.results-section {
  margin-top: 16px;
}

.input-section,
.output-section,
.results-section {
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

.input-wrapper,
.output-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 14px;
}

.results-content {
  flex: 1;
  padding: 16px;
}

.controls-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.time-input {
  width: 100%;
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

.time-input::placeholder {
  color: #94a3b8;
}

.btn-now {
  padding: 8px 16px;
  background: #8b5cf6;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  min-width: 80px;
}

.btn-now:hover {
  background: #7c3aed;
  transform: translateY(-1px);
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

.results-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.result-card {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 12px;
  transition: all 0.2s;
}

.result-card:hover {
  background: #f1f5f9;
}

.result-label {
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 6px;
}

.result-value {
  font-size: 14px;
  font-weight: 500;
  color: #1e293b;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  word-break: break-all;
}

@media (max-width: 1024px) {
  .converter-container {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .time-converter {
    padding: 16px;
  }
}
</style>
