<template>
  <div class="layout" :class="{ 'sidebar-collapsed': !showSidebar }">
    <!-- 头部 -->
    <header class="top-header">
      <div class="logo-section">
        <router-link to="/" class="logo">
          <h1>DevTools</h1>
        </router-link>
        <button class="icon-btn collapse-btn" @click="toggleSidebar">
          <MenuOutline class="collapse-icon" />
        </button>
      </div>
      <div class="header-title">
        <h1 class="title">{{ currentMenuTitle }}</h1>
      </div>
      <div class="header-right">
        <button class="icon-btn">
          <!-- <i class="settings-icon">⚙️</i> -->
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
            exact-active-class="router-link-active"
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
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { MenuOutline, SettingsOutline } from '@vicons/ionicons5'

const route = useRoute()
const currentMenuTitle = ref('')
const showSidebar = ref(true)

// 菜单标题映射
const menuTitles: Record<string, string> = {
  '/': '全部功能列表',
  '/json-editor': 'JSON 编辑器',
  '/xml-editor': 'XML 编辑器',
  '/time-converter': '时间戳转换',
  '/settings': '设置',
}

// 监听路由变化更新标题
watch(
  () => route.path,
  (path) => {
    currentMenuTitle.value = menuTitles[path] || ''
  },
  { immediate: true }
)

const toggleSidebar = () => {
  showSidebar.value = !showSidebar.value
}

// 定义菜单项数据结构
interface MenuItem {
  path: string
  icon: string
  title: string
}

// 菜单项数据
const menuItems: MenuItem[] = [
  { path: '/', icon: '🏠', title: '全部功能列表' },
  { path: '/json-editor', icon: '{ }', title: 'JSON 编辑器' },
  { path: '/xml-editor', icon: '📄', title: 'XML 编辑器' },
  { path: '/time-converter', icon: '⏰', title: '时间戳转换' },
  { path: '/settings', icon: '⚙️', title: '设置' },
]

const searchQuery = ref('')

// 过滤后的菜单项
const filteredMenuItems = computed(() => {
  const query = searchQuery.value.toLowerCase().trim()
  if (!query) return menuItems
  return menuItems.filter((item) => item.title.toLowerCase().includes(query))
})
</script>

<style scoped>
.layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.top-header {
  height: 48px;
  background: #ffffff;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  align-items: center;
  padding: 0 16px;
  justify-content: space-between;
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

.settings-icon {
  font-size: 1.2rem;
}

.main-container {
  flex: 1;
  display: flex;
  background: #ffffff;
}

.sidebar {
  width: 220px;
  background: #f9fafb;
  padding: 16px 0;
  border-right: 1px solid #e5e7eb;
}

.sidebar-hidden {
  width: 0;
  padding: 0;
  overflow: hidden;
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
  padding: 24px;
  background: #ffffff;
}

.collapse-btn {
  margin-left: 8px;
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

.sidebar-collapsed .collapse-icon {
  transform: rotate(180deg);
}
</style>
