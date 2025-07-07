<template>
  <div class="base64-image">
    <div class="tool-section">
      <!-- 顶部：标签页按钮 -->
      <div class="top-header">
        <div class="tab-nav">
          <button
            class="tab-btn"
            :class="{ active: activeTab === 'original' }"
            @click="activeTab = 'original'"
          >
            原始
          </button>
          <button
            class="tab-btn"
            :class="{ active: activeTab === 'data' }"
            @click="activeTab = 'data'"
          >
            数据
          </button>
          <button
            class="tab-btn"
            :class="{ active: activeTab === 'css' }"
            @click="activeTab = 'css'"
          >
            CSS
          </button>
        </div>
        <div class="tab-actions">
          <button class="clear-btn" @click="clearBase64" title="清空">
            × 清空
          </button>
        </div>
      </div>

      <!-- 底部：内容区域 -->
      <div class="content-layout">
        <!-- 左侧：文本区域 -->
        <div class="text-panel">
          <textarea
            :value="getTabContent()"
            @input="
              updateTabContent(
                ($event.target as HTMLTextAreaElement)?.value || ''
              )
            "
            :placeholder="getTabPlaceholder()"
            class="result-area"
            autocomplete="off"
            spellcheck="false"
          ></textarea>
          <div class="status-bar">
            <span class="byte-info">{{ getByteInfo() }}</span>
            <span v-if="imageSize" class="size-info">{{ imageSize }}</span>
          </div>
        </div>

        <!-- 右侧：图片区域 -->
        <div class="image-panel">
          <div class="upload-area" @drop.prevent="handleDrop" @dragover.prevent>
            <input
              type="file"
              ref="fileInput"
              accept="image/*"
              class="hidden"
              @change="handleFileSelect"
            />
            <div
              v-if="!imageUrl"
              class="upload-content"
              @click="triggerFileInput"
            >
              <div class="upload-icon">📤</div>
              <div class="upload-text">
                点击或拖拽图片到此处
                <div class="upload-hint">支持 jpg、png、gif 格式</div>
              </div>
            </div>
            <img
              v-else
              :src="imageUrl"
              class="preview-image"
              @click="triggerFileInput"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { ElMessage } from "element-plus";
import { useClipboard } from "@vueuse/core";
import { useToolsStore } from "../stores/tools";
import { storeToRefs } from "pinia";

const { copy } = useClipboard();
const store = useToolsStore();
const { base64Image } = storeToRefs(store);
const imageUrl = computed({
  get: () => base64Image.value.imageUrl,
  set: (val) => (base64Image.value.imageUrl = val),
});
const base64Result = computed({
  get: () => base64Image.value.base64Result,
  set: (val) => (base64Image.value.base64Result = val),
});
const fileInput = ref<HTMLInputElement | null>(null);

const activeTab = computed({
  get: () => base64Image.value.activeTab || 'original',
  set: (val) => (base64Image.value.activeTab = val),
})

const imageSize = computed({
  get: () => base64Image.value.imageSize || '',
  set: (val) => (base64Image.value.imageSize = val),
})

// 标签页内容
const originalData = computed({
  get: () => base64Image.value.originalData || '',
  set: (val) => (base64Image.value.originalData = val),
})

const cssData = computed({
  get: () => base64Image.value.cssData || '',
  set: (val) => (base64Image.value.cssData = val),
})

const exampleData = computed({
  get: () => base64Image.value.exampleData || '',
  set: (val) => (base64Image.value.exampleData = val),
})

const triggerFileInput = () => {
  fileInput.value?.click();
};

const handleFileSelect = (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0];
  if (file) {
    handleImage(file);
  }
};

const handleDrop = (event: DragEvent) => {
  const file = event.dataTransfer?.files[0];
  if (file) {
    handleImage(file);
  }
};

const handleImage = (file: File) => {
  if (!file.type.startsWith("image/")) {
    ElMessage.error("请选择图片文件");
    return;
  }

  const reader = new FileReader();
  reader.onload = (e) => {
    const result = e.target?.result as string;
    imageUrl.value = result;

    // 创建临时图片获取尺寸
    const img = new Image();
    img.onload = () => {
      imageSize.value = `${img.width}x${img.height}`;
      // 自动编码
      encodeImageToBase64();
    };
    img.src = result;
  };
  reader.readAsDataURL(file);
};

// 自动编码图片为Base64
const encodeImageToBase64 = () => {
  if (!imageUrl.value) {
    return;
  }

  try {
    const dataUrl = imageUrl.value;
    const base64Data = dataUrl.split(",")[1] || "";

    // 更新各个标签页的内容
    originalData.value = base64Data;
    base64Result.value = dataUrl;

    // 生成CSS格式
    cssData.value = `background-image: url('${dataUrl}');`;
  } catch (error) {
    console.error("编码失败:", error);
  }
};

// 自动解码Base64为图片
const decodeBase64ToImage = (base64String: string) => {
  if (!base64String || !base64String.trim()) {
    return;
  }

  try {
    // 清理Base64字符串，移除换行符和空格
    const cleanBase64 = base64String.replace(/\s+/g, "");

    // 如果已经是完整的data URL格式
    if (cleanBase64.startsWith("data:image/")) {
      imageUrl.value = cleanBase64;

      // 获取图片尺寸
      const img = new Image();
      img.onload = () => {
        imageSize.value = `${img.width}x${img.height}`;
      };
      img.src = cleanBase64;
    }
    // 如果只是Base64字符串，尝试添加data URL前缀
    else if (cleanBase64.length > 100) {
      // 默认尝试jpeg格式
      const dataUrl = `data:image/jpeg;base64,${cleanBase64}`;
      imageUrl.value = dataUrl;

      // 获取图片尺寸
      const img = new Image();
      img.onload = () => {
        imageSize.value = `${img.width}x${img.height}`;
      };
      img.onerror = () => {
        // 如果jpeg失败，尝试png
        const pngDataUrl = `data:image/png;base64,${cleanBase64}`;
        imageUrl.value = pngDataUrl;
      };
      img.src = dataUrl;
    }
  } catch (error) {
    console.error("解码失败:", error);
  }
};

// 清空图片
const clearImage = () => {
  imageUrl.value = "";
  if (fileInput.value) {
    fileInput.value.value = "";
  }
};

// 清空 Base64
const clearBase64 = () => {
  base64Result.value = "";
  originalData.value = "";
  cssData.value = "";
  imageUrl.value = "";
  imageSize.value = "";
  if (fileInput.value) {
    fileInput.value.value = "";
  }
};

// 获取当前标签页内容
const getTabContent = () => {
  switch (activeTab.value) {
    case "original":
      return originalData.value;
    case "data":
      return base64Result.value;
    case "css":
      return cssData.value;
    default:
      return "";
  }
};

// 更新当前标签页内容
const updateTabContent = (value: string) => {
  switch (activeTab.value) {
    case "original":
      originalData.value = value;
      // 尝试解码Base64
      if (value && value.length > 100) {
        decodeBase64ToImage(value);
      }
      break;
    case "data":
      base64Result.value = value;
      // 尝试解码Base64
      if (value && value.length > 100) {
        decodeBase64ToImage(value);
      }
      break;
    case "css":
      cssData.value = value;
      break;
  }
};

// 获取标签页占位符
const getTabPlaceholder = () => {
  switch (activeTab.value) {
    case "original":
      return "纯Base64数据（不包含data:image前缀）...";
    case "data":
      return "data:image/...格式的完整Base64数据...";
    case "css":
      return "CSS background-image 格式...";
    default:
      return "";
  }
};

// 获取字节信息
const getByteInfo = () => {
  const content = getTabContent();
  if (!content) return "0 bytes";

  const bytes = new TextEncoder().encode(content).length;
  return `${bytes} bytes`;
};
</script>

<style scoped>
.base64-image {
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
  margin-bottom: 16px;
}

.content-layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  height: calc(100% - 48px);
  align-items: stretch;
}

.text-panel,
.image-panel {
  background: white;
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

.tab-nav {
  display: flex;
  align-items: stretch;
  border: 1px solid #d1d5db;
}

.tab-btn {
  padding: 0 10px;
  background: #f8f9fa;
  border: none;
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
}

.tab-btn:last-child {
  border-right: none;
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
  gap: 8px;
  padding: 0 12px;
  background: #ffffff;
}

.clear-btn {
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

.clear-btn:hover {
  background: #e9ecef;
}

.text-panel {
  border: none;
  box-shadow: none;
}

.text-panel .result-area {
  flex: 1;
  border: 1px solid #d1d5db;
}

.status-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 12px;
  background: #f3f4f6;
  border-top: 1px solid #d1d5db;
  font-size: 10px;
  color: #6b7280;
  height: 20px;
  min-height: 20px;
}

.byte-info,
.size-info {
  font-family: "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Consolas,
    "Courier New", monospace;
  font-size: 10px;
}

.btn-clear {
  padding: 6px 12px;
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fecaca;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-clear:hover {
  background: #fee2e2;
}

.image-panel .upload-area {
  flex: 1;
  margin: 16px;
  border: 2px dashed #d1d5db;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  transition: all 0.2s;
  overflow: hidden;
}

.upload-area:hover {
  border-color: #60a5fa;
  background: #f8fafc;
}

.preview-image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.upload-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 40px;
}

.upload-icon {
  font-size: 48px;
}

.upload-text {
  color: #374151;
  font-size: 14px;
  text-align: center;
}

.upload-hint {
  color: #6b7280;
  font-size: 12px;
  margin-top: 4px;
}

.result-area {
  flex: 1;
  padding: 8px;
  border: none;
  outline: none;
  font-family: "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Consolas,
    "Courier New", monospace;
  font-size: 10px;
  line-height: 1.3;
  resize: none;
  background: #ffffff;
  color: #111827;
  white-space: pre-wrap;
  word-wrap: break-word;
  word-break: break-all;
}

.result-area::placeholder {
  color: #9ca3af;
  font-size: 10px;
}

.result-area:focus {
  outline: none;
  box-shadow: none;
}

.preview-image {
  cursor: pointer;
}

.preview-image:hover {
  opacity: 0.8;
}

.hidden {
  display: none;
}
</style>
