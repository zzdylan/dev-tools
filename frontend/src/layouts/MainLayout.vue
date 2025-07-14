<template>
  <div
    class="layout"
    :class="{ 'sidebar-collapsed': !showSidebar, 'mac-layout': isMac }"
  >
    <!-- 头部 -->
    <header class="top-header">
      <!-- Mac风格的窗口控制按钮 - 左边 -->
      <div v-if="isMac" class="mac-window-controls">
        <button
          class="mac-btn mac-close-btn"
          @click="closeWindow"
          title="关闭"
        ></button>
        <button
          class="mac-btn mac-minimize-btn"
          @click="minimizeWindow"
          title="最小化"
        ></button>
        <button
          class="mac-btn mac-fullscreen-btn"
          @click="toggleFullscreen"
          title="全屏"
        ></button>
      </div>

      <div class="logo-section">
        <!-- 移除DevTools文字 -->
      </div>
      <button
        class="icon-btn collapse-btn"
        @click="toggleSidebar"
        :title="`${showSidebar ? '隐藏' : '显示'}菜单栏`"
      >
        <MenuOutline class="collapse-icon" />
      </button>
      <div class="header-title">
        <h1 class="title">{{ currentMenuTitle }}</h1>
      </div>

      <!-- 设置按钮 -->
      <div class="header-actions">
        <button
          class="icon-btn settings-btn"
          @click="toggleSettings"
          title="设置"
        >
          <SettingsOutline class="settings-icon" />
        </button>
      </div>

      <!-- Windows风格的窗口控制按钮 - 右边 -->
      <div v-if="!isMac" class="header-right">
        <button
          class="window-btn minimize-btn"
          @click="minimizeWindow"
          title="最小化"
        >
          <span class="window-icon">−</span>
        </button>
        <button class="window-btn close-btn" @click="closeWindow" title="关闭">
          <span class="window-icon">×</span>
        </button>
      </div>
    </header>

    <div class="main-container">
      <!-- 左侧菜单 -->
      <aside class="sidebar" :class="{ 'sidebar-hidden': !showSidebar }">
        <div class="search-box">
          <span class="search-icon">🔍</span>
          <input
            type="text"
            v-model="searchQuery"
            placeholder="Search..."
            class="search-input"
          />
        </div>
        <nav class="side-nav">
          <router-link
            v-for="item in filteredMenuItems"
            :key="item.path"
            :to="item.path"
            class="nav-item"
            :exact="item.path === '/'"
            active-class=""
            exact-active-class=""
            :class="{ 'router-link-active': isActiveRoute(item.path) }"
          >
            <span class="nav-icon">{{ item.icon }}</span>
            {{ item.title }}
          </router-link>
        </nav>
      </aside>

      <!-- 主要内容区 -->
      <main class="content">
        <router-view></router-view>
      </main>
    </div>
  </div>

  <!-- 设置对话框 -->
  <div v-if="showSettings" class="settings-overlay" @click="closeSettings">
    <div class="settings-modal" @click.stop>
      <div class="settings-header">
        <div class="settings-controls">
          <button class="settings-control close" @click="closeSettings"></button>
          <button class="settings-control minimize"></button>
          <button class="settings-control maximize"></button>
        </div>
        <h2 class="settings-title">设置</h2>
      </div>
      
      <div class="settings-content">
        <div class="settings-panel single-panel">
          <div class="settings-section">
            <h3 class="section-title">数据管理</h3>
            <div class="setting-item">
              <div class="setting-description">
                <label class="setting-label">清空缓存</label>
                <p class="setting-hint">清空所有本地存储的数据，包括工具配置、历史记录等</p>
              </div>
              <button class="clear-cache-btn" @click="clearCache">清空缓存</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import { MenuOutline, SettingsOutline } from "@vicons/ionicons5";
import { useToolsStore } from "../stores/tools";
import { storeToRefs } from "pinia";
import { ElMessage, ElMessageBox } from 'element-plus';
import {
  MinimizeWindow,
  CloseWindow,
  ToggleFullscreen,
} from "../../wailsjs/go/main/App";

const route = useRoute();
const currentMenuTitle = ref("");
const showSidebar = ref(true);

// 检测是否是Mac平台
const isMac = ref(false);
onMounted(() => {
  // 检测用户代理或使用Wails API检测平台
  isMac.value = navigator.platform.toUpperCase().indexOf("MAC") >= 0;
});

// 菜单标题映射
const menuTitles: Record<string, string> = {
  "/": "全部功能列表",
  "/json-editor": "JSON 编辑器",
  "/xml-editor": "XML 编辑器",
  "/time-converter": "时间戳转换",
  // '/settings': '设置',
  "/url-converter": "URL 编解码",
  "/url-parser": "URL 解析",
  "/qrcode": "二维码工具",
  "/base64-image": "Base64 图像",
  "/base64-text": "Base64 文本",
  "/number-converter": "进制转换",
  "/text-diff": "文本对比",
  "/curl-converter": "cURL 转换",
  "/unicode-converter": "Unicode 转换",
  "/json-to-go": "JSON转Go",
  "/jwt-decoder": "JWT 解析",
};

// 监听路由变化更新标题
watch(
  () => route.path,
  (path) => {
    currentMenuTitle.value = menuTitles[path] || "";
  },
  { immediate: true }
);

const toggleSidebar = () => {
  showSidebar.value = !showSidebar.value;
};

// 设置相关
const showSettings = ref(false);

const toggleSettings = () => {
  showSettings.value = !showSettings.value;
};

const closeSettings = () => {
  showSettings.value = false;
};

// 清空缓存功能
const clearCache = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要清空所有缓存数据吗？\n\n这将清空所有工具的配置、历史记录等数据，操作不可撤销。',
      '清空缓存确认',
      {
        confirmButtonText: '确定清空',
        cancelButtonText: '取消',
        type: 'warning',
        dangerouslyUseHTMLString: false
      }
    );
    
    // 清空所有localStorage数据
    localStorage.clear();
    
    // 清空所有sessionStorage数据
    sessionStorage.clear();
    
    // 重置store状态
    store.$reset();
    
    ElMessage.success('缓存已清空，建议重启应用以确保完全生效');
    
    // 关闭设置对话框
    closeSettings();
    
  } catch (err) {
    // 用户取消操作时不显示错误
    if (err === 'cancel') {
      return;
    }
    ElMessage.error('清空缓存失败');
  }
};

const store = useToolsStore();
const { menuConfig } = storeToRefs(store);

// 在组件挂载时初始化菜单
onMounted(() => {
  store.initializeMenu();
});

const menuItems = computed(() => {
  const home = { path: "/", icon: "🏠", title: "全部功能列表" };
  const visibleItems = menuConfig.value.items
    .filter((item) => item.visible)
    .sort((a, b) => {
      // 确保排序稳定性
      const orderDiff = a.order - b.order;
      return orderDiff !== 0 ? orderDiff : a.id.localeCompare(b.id);
    });
  return [home, ...visibleItems];
});

const searchQuery = ref("");

// 过滤后的菜单项
const filteredMenuItems = computed(() => {
  const query = searchQuery.value.toLowerCase().trim();
  if (!query) return menuItems.value;
  return menuItems.value.filter((item) =>
    item.title.toLowerCase().includes(query)
  );
});

// 检查路由是否激活，处理重定向情况
const isActiveRoute = (path: string) => {
  const currentPath = route.path;

  // 首页只有完全匹配才激活
  if (path === "/") {
    return currentPath === "/";
  }

  // 处理重定向路由的特殊情况
  if (path === "/json-editor" && currentPath.startsWith("/json-editor")) {
    return true;
  }
  if (path === "/xml-editor" && currentPath.startsWith("/xml-editor")) {
    return true;
  }

  // 普通路由匹配
  return currentPath === path;
};

// 窗口控制方法
const minimizeWindow = () => {
  MinimizeWindow();
};

const closeWindow = () => {
  CloseWindow();
};

const toggleFullscreen = () => {
  // TODO: Implement when wails bindings are updated
  // console.log("Toggle fullscreen");
  ToggleFullscreen();
};
</script>

<style scoped>
.layout {
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.top-header {
  height: 48px;
  background: #ffffff;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  align-items: center;
  padding: 0 16px;
  --wails-draggable: drag;
}

.logo-section {
  display: flex;
  align-items: center;
  gap: 8px;
}

.logo {
  text-decoration: none;
  color: #000;
}

.logo h1 {
  font-size: 1.2rem;
  font-weight: 600;
  margin: 0;
}

.header-title {
  flex: 1;
  text-align: center;
}

.header-title .title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #24292f;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 4px;
}

.icon-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;
}

.icon-btn:hover {
  background: #f3f4f6;
}

.window-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.window-btn:hover {
  background: #f3f4f6;
}

.minimize-btn:hover {
  background: #e5e7eb;
}

.close-btn:hover {
  background: #fee2e2;
  color: #dc2626;
}

.window-icon {
  font-size: 16px;
  font-weight: 500;
  color: #6b7280;
  line-height: 1;
}

.close-btn:hover .window-icon {
  color: #dc2626;
}

.settings-icon {
  font-size: 1.2rem;
}

.main-container {
  flex: 1;
  display: flex;
  background: #ffffff;
  min-height: 0;
  overflow: hidden;
}

.sidebar {
  width: 220px;
  min-width: 220px;
  flex-shrink: 0;
  background: #f9fafb;
  padding: 16px 0;
  border-right: 1px solid #e5e7eb;
}

.sidebar-hidden {
  width: 0 !important;
  min-width: 0 !important;
  padding: 0 !important;
  margin: 0 !important;
  overflow: hidden;
  border-right: none !important;
}

.search-box {
  padding: 0 16px;
  margin-bottom: 16px;
  position: relative;
}

.search-icon {
  position: absolute;
  left: 24px;
  top: 50%;
  transform: translateY(-50%);
  color: #6b7280;
  font-size: 14px;
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 8px 12px;
  padding-left: 36px;
  border: none;
  border-radius: 4px;
  font-size: 13px;
  outline: none;
  background: #edf0f4;
  color: #24292f;
}

.search-input::placeholder {
  color: #6b7280;
}

.search-input:focus {
  background: #e5e7eb;
}

/* 当没有搜索结果时的样式 */
.no-results {
  padding: 16px;
  text-align: center;
  color: #6b7280;
  font-size: 0.875rem;
}

.side-nav {
  display: flex;
  flex-direction: column;
}

.nav-item {
  padding: 8px 16px;
  color: #374151;
  text-decoration: none;
  font-size: 0.875rem;
  display: flex;
  align-items: center;
  gap: 8px;
}

.nav-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.nav-item:hover {
  background: #f3f4f6;
}

.nav-item.router-link-active {
  background: #e5e7eb;
  color: #000;
  font-weight: 500;
}

.content {
  flex: 1;
  min-width: 0;
  min-height: 0;
  padding: 0;
  background: #ffffff;
  overflow: auto;
}

.collapse-btn {
  padding: 4px 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #e5e7eb;
  border-radius: 4px;
}

.collapse-btn:hover {
  background: #f3f4f6;
}

.collapse-icon {
  width: 20px;
  height: 20px;
  color: #666;
}

/* 当侧边栏隐藏时，旋转按钮图标 */
.layout:not(.sidebar-collapsed) .collapse-icon {
  transform: rotate(0deg);
}

.sidebar-collapsed .collapse-icon {
  transform: rotate(180deg);
}

/* Mac风格窗口控制按钮 */
.mac-window-controls {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-right: 16px;
}

.mac-btn {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: none;
  cursor: pointer;
  transition: opacity 0.2s;
}

.mac-btn:hover {
  opacity: 0.8;
}

.mac-close-btn {
  background: #ff5f57;
  position: relative;
}

.mac-close-btn:hover::after {
  content: "×";
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: #4a0e0e;
  font-size: 10px;
  font-weight: bold;
  line-height: 1;
}

.mac-minimize-btn {
  background: #ffbd2e;
  position: relative;
}

.mac-minimize-btn:hover::after {
  content: "−";
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: #7a5200;
  font-size: 10px;
  font-weight: bold;
  line-height: 1;
}

.mac-fullscreen-btn {
  background: #28ca42;
  position: relative;
}

.mac-fullscreen-btn:hover::after {
  content: "⤢";
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: #0f3e17;
  font-size: 8px;
  font-weight: bold;
  line-height: 1;
}

/* Mac环境下调整collapse按钮位置 */
.layout.mac-layout .collapse-btn {
  margin-left: 150px; /* Mac环境下减少margin，因为左边有窗口控制按钮 */
}

/* Windows/Linux环境下的collapse按钮位置 */
.layout:not(.mac-layout) .collapse-btn {
  margin-left: 210px;
}

/* 头部操作按钮 */
.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-right: 16px;
}

.settings-btn {
  padding: 6px;
  border: 1px solid #e5e7eb;
  border-radius: 4px;
}

.settings-btn:hover {
  background: #f3f4f6;
}

.settings-icon {
  width: 18px;
  height: 18px;
  color: #666;
}

/* 设置对话框样式 */
.settings-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.settings-modal {
  background: white;
  border-radius: 8px;
  width: 600px;
  max-width: 90vw;
  max-height: 80vh;
  overflow: hidden;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
}

.settings-header {
  position: relative;
  padding: 16px 20px;
  border-bottom: 1px solid #e5e7eb;
  background: #f8f9fa;
}

.settings-controls {
  position: absolute;
  top: 16px;
  left: 20px;
  display: flex;
  gap: 8px;
}

.settings-control {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: none;
  cursor: pointer;
  transition: opacity 0.2s;
}

.settings-control:hover {
  opacity: 0.8;
}

.settings-control.close {
  background: #ff5f57;
}

.settings-control.minimize {
  background: #ffbd2e;
}

.settings-control.maximize {
  background: #28ca42;
}

.settings-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #24292f;
  text-align: center;
}

.settings-content {
  display: flex;
  min-height: 200px;
}

.settings-panel {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

.settings-panel.single-panel {
  width: 100%;
}

.settings-section {
  margin-bottom: 24px;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  margin: 0 0 12px 0;
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

.setting-description {
  flex: 1;
}

.setting-label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 4px;
  display: block;
}

.setting-hint {
  font-size: 12px;
  color: #6b7280;
  margin: 0;
  line-height: 1.4;
}

.clear-cache-btn {
  padding: 8px 16px;
  background: #dc2626;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
  margin-left: 16px;
}

.clear-cache-btn:hover {
  background: #b91c1c;
}
</style>
