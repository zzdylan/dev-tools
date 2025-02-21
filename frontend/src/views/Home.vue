<template>
  <div class="home">
    <div class="tools-grid">
      <router-link to="/json-editor" class="tool-card">
        <div class="tool-icon">{ }</div>
        <h3>JSON 编辑器</h3>
        <p>JSON 格式化、验证、转换工具</p>
      </router-link>

      <router-link to="/xml-editor" class="tool-card">
        <div class="tool-icon">📄</div>
        <h3>XML 编辑器</h3>
        <p>XML 格式化、验证、转换工具</p>
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts"></script>

<style scoped>
.home {
  padding: 20px;
}

.tools-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
}

.tool-card {
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 20px;
  text-decoration: none;
  color: inherit;
  transition: all 0.2s;
}

.tool-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.tool-icon {
  font-size: 24px;
  margin-bottom: 12px;
}

.tool-card h3 {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: #111827;
}

.tool-card p {
  margin: 0;
  font-size: 14px;
  color: #6b7280;
}
</style>
