<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { EventsOn, WindowIsFullscreen, WindowIsMaximised } from '../wailsjs/runtime/runtime'

const maximised = ref(false)
const hideRadius = ref(false)

// 检测平台
const isMacOS = () => {
  return navigator.userAgent.toUpperCase().indexOf('MAC') >= 0
}

const isWindows = () => {
  return navigator.userAgent.toUpperCase().indexOf('WIN') >= 0
}

// 计算容器样式
const wrapperStyle = computed(() => {
  if (isWindows()) {
    return {}
  }
  return hideRadius.value
    ? {}
    : {
        border: '1px solid #e5e7eb',
        borderRadius: '10px',
      }
})

const onToggleFullscreen = (fullscreen: boolean) => {
  hideRadius.value = fullscreen
}

const onToggleMaximize = (isMaximised: boolean) => {
  if (isMaximised) {
    maximised.value = true
    if (!isMacOS()) {
      hideRadius.value = true
    }
  } else {
    maximised.value = false
    if (!isMacOS()) {
      hideRadius.value = false
    }
  }
}

EventsOn('window_changed', (info: any) => {
  const { fullscreen, maximised } = info
  onToggleFullscreen(fullscreen === true)
  onToggleMaximize(maximised)
})

onMounted(async () => {
  const fullscreen = await WindowIsFullscreen()
  onToggleFullscreen(fullscreen === true)
  const maximisedState = await WindowIsMaximised()
  onToggleMaximize(maximisedState)
})
</script>

<template>
  <div id="app-wrapper" :style="wrapperStyle">
    <router-view></router-view>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto,
    'Helvetica Neue', Arial, sans-serif;
  line-height: 1.4;
  color: #374151;
  font-size: 13px;
  background-color: transparent;
}

#app-wrapper {
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  box-sizing: border-box;
  background-color: #ffffff;
}
</style>
