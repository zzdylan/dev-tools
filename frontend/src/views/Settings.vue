<template>
  <div class="settings">
    <h2>设置</h2>
    <div class="settings-content">
      <!-- 设置选项示例 -->
      <div class="setting-item">
        <h3>常规设置</h3>
        <div class="setting-option">
          <label>
            <input type="checkbox" v-model="settings.autoStart" /> 开机自启
          </label>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Settings {
  autoStart: boolean
}

const settings = ref<Settings>({
  autoStart: false,
})
</script>

<style scoped>
.settings {
  max-width: 800px;
  margin: 0 auto;
}

.setting-item {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  margin-top: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.setting-option {
  margin-top: 16px;
}
</style>
