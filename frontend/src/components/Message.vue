<template>
  <Transition name="fade">
    <div v-if="visible" class="message" :class="type">
      {{ content }}
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const visible = ref(false)
const content = ref('')
const type = ref('success')

const show = (message: string, messageType = 'success') => {
  content.value = message
  type.value = messageType
  visible.value = true
  setTimeout(() => {
    visible.value = false
  }, 2000)
}

defineExpose({ show })
</script>

<style scoped>
.message {
  position: fixed;
  top: 20px;
  right: 20px;
  padding: 10px 20px;
  border-radius: 4px;
  z-index: 9999;
  color: #fff;
  font-size: 14px;
}

.success {
  background-color: #10b981;
}

.error {
  background-color: #ef4444;
}

.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
