<template>
    <div class="jwt-decoder">

        <div class="input-section">
            <div class="section-title">JWT Token</div>
            <div class="textarea-wrapper">
                <textarea v-model="jwtToken" placeholder="请输入JWT Token..." class="jwt-input" rows="6"></textarea>
                <div class="tools-group">
                    <button class="tool-btn" @click="decodeJwt">
                        <span class="tool-icon">🔍</span>
                        解析
                    </button>
                </div>
            </div>
        </div>

        <div class="output-section">
            <div class="result-content">
                <div v-if="error" class="error-message">
                    {{ error }}
                </div>

                <div v-else-if="jwtData" class="result-panel">
                    <div class="section-header">
                        <span class="section-title">解析结果</span>
                        <button class="copy-btn" @click="copyToClipboard(JSON.stringify(jwtData, null, 2))">
                            <span class="tool-icon">📋</span>
                            复制
                        </button>
                    </div>
                    <pre class="json-result">{{ JSON.stringify(jwtData, null, 2) }}</pre>
                </div>

                <div v-else class="placeholder">
                    请输入JWT Token进行解析
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
    padding: 20px;
    background: #f5f5f5;
    min-height: 100vh;
}

.title {
    font-size: 24px;
    font-weight: bold;
    margin-bottom: 20px;
    color: #333;
}

.input-section {
    background: white;
    border-radius: 8px;
    padding: 20px;
    margin-bottom: 20px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.section-title {
    font-size: 16px;
    font-weight: bold;
    margin-bottom: 10px;
    color: #333;
}

.textarea-wrapper {
    position: relative;
}

.jwt-input {
    width: 100%;
    padding: 12px;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-family: 'Consolas', 'Monaco', monospace;
    font-size: 14px;
    line-height: 1.4;
    resize: vertical;
    outline: none;
    transition: border-color 0.3s;
}

.jwt-input:focus {
    border-color: #409eff;
}

.tools-group {
    display: flex;
    gap: 8px;
    margin-top: 12px;
}

.tool-btn {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 8px 12px;
    background: #409eff;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    transition: background-color 0.3s;
}

.tool-btn:hover {
    background: #66b1ff;
}

.tool-icon {
    font-size: 14px;
}

.output-section {
    background: white;
    border-radius: 8px;
    padding: 20px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.result-content {
    min-height: 200px;
}

.error-message {
    color: #f56c6c;
    background: #fef0f0;
    padding: 12px;
    border-radius: 4px;
    border: 1px solid #fbc4c4;
}

.result-section {
    animation: fadeIn 0.3s ease-in-out;
}

.section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
}

.copy-btn {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 6px 10px;
    background: #f0f0f0;
    color: #666;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 12px;
    transition: background-color 0.3s;
}

.copy-btn:hover {
    background: #e0e0e0;
}

.json-result {
    background: #f8f9fa;
    border: 1px solid #e9ecef;
    border-radius: 4px;
    padding: 16px;
    font-family: 'Consolas', 'Monaco', monospace;
    font-size: 14px;
    line-height: 1.5;
    overflow-x: auto;
    white-space: pre-wrap;
    color: #333;
}

.placeholder {
    text-align: center;
    color: #999;
    font-size: 16px;
    padding: 40px;
}

@keyframes fadeIn {
    from {
        opacity: 0;
        transform: translateY(10px);
    }

    to {
        opacity: 1;
        transform: translateY(0);
    }
}
</style>