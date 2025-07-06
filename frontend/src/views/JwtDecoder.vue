<template>
    <div class="jwt-decoder">
        <div class="decoder-container">
            <!-- JWT输入区域 -->
            <div class="input-section">
                <div class="section-header">
                    <h3>JWT Token</h3>
                    <button class="btn-decode" @click="decodeJwt">解析</button>
                </div>
                <div class="input-wrapper">
                    <textarea 
                        v-model="jwtToken" 
                        placeholder="粘贴JWT Token..."
                        class="jwt-input"
                    ></textarea>
                </div>
            </div>

            <!-- 解析结果区域 -->
            <div class="output-section">
                <div class="section-header">
                    <h3>解析结果</h3>
                    <button 
                        v-if="jwtData" 
                        class="btn-copy" 
                        @click="copyToClipboard(JSON.stringify(jwtData, null, 2))"
                    >
                        复制结果
                    </button>
                </div>
                <div class="output-wrapper">
                    <div v-if="error" class="error-message">
                        <div class="error-icon">⚠️</div>
                        <div class="error-text">{{ error }}</div>
                    </div>

                    <div v-else-if="jwtData" class="result-content">
                        <pre class="json-result">{{ JSON.stringify(jwtData, null, 2) }}</pre>
                    </div>

                    <div v-else class="placeholder">
                        <div class="placeholder-icon">🔍</div>
                        <div class="placeholder-text">输入JWT Token开始解析</div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useClipboard } from '@vueuse/core'
import { ElMessage } from 'element-plus'
import { useToolsStore } from '../stores/tools'
import { storeToRefs } from 'pinia'

const { copy } = useClipboard()
const store = useToolsStore()
const { jwtDecoder } = storeToRefs(store)

const jwtToken = computed({
    get: () => jwtDecoder.value.token,
    set: (val) => {
        jwtDecoder.value.token = val
    }
})

const jwtData = ref<any>(null)
const error = ref('')

// 解析JWT Token
const decodeJwt = () => {
    error.value = ''
    jwtData.value = null

    if (!jwtToken.value.trim()) {
        error.value = '请输入JWT Token'
        return
    }

    try {
        const parts = jwtToken.value.trim().split('.')

        if (parts.length !== 3) {
            error.value = '无效的JWT格式，JWT应该包含3个部分（header.payload.signature）'
            return
        }

        // 解析Header
        const header = JSON.parse(atob(parts[0].replace(/-/g, '+').replace(/_/g, '/')))

        // 解析Payload
        const payload = JSON.parse(atob(parts[1].replace(/-/g, '+').replace(/_/g, '/')))

        // 保存原始签名
        const signature = parts[2]

        jwtData.value = {
            header,
            payload,
            signature
        }

        ElMessage.success('JWT解析成功')
    } catch (err) {
        error.value = 'JWT解析失败：' + (err instanceof Error ? err.message : '未知错误')
    }
}

// 复制到剪贴板
const copyToClipboard = async (text: string) => {
    try {
        await copy(text)
        ElMessage.success('已复制到剪贴板')
    } catch (err) {
        ElMessage.error('复制失败')
    }
}

// 自动解析（当输入变化时）
const autoDecodeIfValid = () => {
    if (jwtToken.value.trim() && jwtToken.value.split('.').length === 3) {
        decodeJwt()
    }
}

// 监听输入变化
watch(() => jwtToken.value, () => {
    if (jwtToken.value.trim()) {
        // 延迟自动解析，避免频繁解析
        setTimeout(autoDecodeIfValid, 500)
    } else {
        jwtData.value = null
        error.value = ''
    }
})
</script>

<style scoped>
.jwt-decoder {
    padding: 24px;
    max-width: 1200px;
    margin: 0 auto;
    background: #f8fafc;
    min-height: 100vh;
}

.decoder-container {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 24px;
    height: calc(100vh - 60px);
}

.input-section,
.output-section {
    background: white;
    border-radius: 12px;
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
    padding: 20px 24px;
    border-bottom: 1px solid #e2e8f0;
    background: #f8fafc;
}

.section-header h3 {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
    color: #1e293b;
}

.btn-decode {
    padding: 8px 16px;
    background: #3b82f6;
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
}

.btn-decode:hover {
    background: #2563eb;
}

.btn-copy {
    padding: 6px 12px;
    background: #f1f5f9;
    color: #475569;
    border: 1px solid #e2e8f0;
    border-radius: 6px;
    font-size: 13px;
    cursor: pointer;
    transition: all 0.2s;
}

.btn-copy:hover {
    background: #e2e8f0;
}

.input-wrapper,
.output-wrapper {
    flex: 1;
    display: flex;
    flex-direction: column;
}

.jwt-input {
    flex: 1;
    padding: 20px;
    border: none;
    outline: none;
    font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
    font-size: 14px;
    line-height: 1.5;
    resize: none;
    background: transparent;
    color: #1e293b;
}

.jwt-input::placeholder {
    color: #94a3b8;
}

.output-wrapper {
    padding: 20px;
}

.error-message {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px;
    background: #fef2f2;
    border: 1px solid #fecaca;
    border-radius: 8px;
    color: #dc2626;
}

.error-icon {
    font-size: 20px;
}

.error-text {
    font-size: 14px;
    font-weight: 500;
}

.result-content {
    flex: 1;
}

.json-result {
    background: #f8fafc;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    padding: 16px;
    font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
    font-size: 14px;
    line-height: 1.5;
    overflow-x: auto;
    white-space: pre-wrap;
    color: #1e293b;
    margin: 0;
}

.placeholder {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    gap: 12px;
    color: #94a3b8;
}

.placeholder-icon {
    font-size: 48px;
    opacity: 0.5;
}

.placeholder-text {
    font-size: 16px;
    font-weight: 500;
}

@media (max-width: 1024px) {
    .decoder-container {
        grid-template-columns: 1fr;
        gap: 20px;
    }
    
    .jwt-decoder {
        padding: 16px;
    }
}
</style>