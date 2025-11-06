<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { Quit, WindowMinimise, WindowToggleMaximise, EventsOn, WindowIsMaximised } from '../../wailsjs/runtime/runtime'
import WindowMin from './icons/WindowMin.vue'
import WindowMax from './icons/WindowMax.vue'
import WindowRestore from './icons/WindowRestore.vue'
import WindowClose from './icons/WindowClose.vue'

const props = defineProps({
  size: {
    type: Number,
    default: 38,
  },
  maximised: {
    type: Boolean,
  },
})

const emit = defineEmits(['update:maximised'])

const localMaximised = ref(props.maximised)

const buttonSize = computed(() => {
  return props.size + 'px'
})

const handleMinimise = async () => {
  WindowMinimise()
}

const handleMaximise = () => {
  WindowToggleMaximise()
}

const handleClose = () => {
  Quit()
}

// 监听窗口状态变化
EventsOn('window_changed', (info: any) => {
  const { maximised } = info
  localMaximised.value = maximised
  emit('update:maximised', maximised)
})

onMounted(async () => {
  const maximised = await WindowIsMaximised()
  localMaximised.value = maximised
  emit('update:maximised', maximised)
})

watch(() => props.maximised, (val) => {
  localMaximised.value = val
})
</script>

<template>
  <div class="control-wrapper">
    <div class="btn-wrapper" @click="handleMinimise" title="最小化">
      <window-min />
    </div>
    <div v-if="localMaximised" class="btn-wrapper" @click="handleMaximise" title="还原">
      <window-restore />
    </div>
    <div v-else class="btn-wrapper" @click="handleMaximise" title="最大化">
      <window-max />
    </div>
    <div class="btn-wrapper close-btn" @click="handleClose" title="关闭">
      <window-close />
    </div>
  </div>
</template>

<style scoped>
.control-wrapper {
  display: flex;
  align-items: stretch;
  height: 100%;
}

.btn-wrapper {
  width: v-bind('buttonSize');
  height: v-bind('buttonSize');
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  --wails-draggable: none;
  color: #6b7280;
  transition: all 0.2s ease;
}

.btn-wrapper:hover {
  background-color: rgba(0, 0, 0, 0.05);
}

.btn-wrapper:active {
  background-color: rgba(0, 0, 0, 0.1);
}

.btn-wrapper.close-btn:hover {
  background-color: #ef4444;
  color: white;
}

.btn-wrapper.close-btn:active {
  background-color: #dc2626;
  color: white;
}
</style>
