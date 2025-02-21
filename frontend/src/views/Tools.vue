<template>
  <div class="tools">
    <h2>工具库</h2>
    <div class="tools-grid">
      <!-- 工具卡片示例 -->
      <div class="tool-card">
        <h3>开发工具1</h3>
        <p>工具描述...</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Tool {
  id: number
  name: string
  description: string
}

// 这里可以添加 TypeScript 代码
</script>

<style scoped>
.tools {
  max-width: 1200px;
  margin: 0 auto;
}

.tools-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 24px;
  margin-top: 24px;
}

.tool-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}
</style>
