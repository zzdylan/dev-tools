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

      <router-link to="/time-converter" class="tool-card">
        <div class="tool-icon">⏰</div>
        <h3>时间戳转换</h3>
        <p>时间戳与日期时间互转工具</p>
      </router-link>

      <router-link to="/url-converter" class="tool-card">
        <div class="tool-icon">🔗</div>
        <h3>URL 编解码</h3>
        <p>URL 编码解码转换工具</p>
      </router-link>

      <router-link to="/url-parser" class="tool-card">
        <div class="tool-icon">🔍</div>
        <h3>URL 解析</h3>
        <p>URL 参数解析分析工具</p>
      </router-link>

      <router-link to="/qrcode" class="tool-card">
        <div class="tool-icon">📱</div>
        <h3>二维码工具</h3>
        <p>二维码生成与解析工具</p>
      </router-link>

      <router-link to="/base64-image" class="tool-card">
        <div class="tool-icon">🖼️</div>
        <h3>Base64 图像</h3>
        <p>图片与 Base64 编解码转换</p>
      </router-link>

      <router-link to="/base64-text" class="tool-card">
        <div class="tool-icon">📝</div>
        <h3>Base64 文本</h3>
        <p>文本 Base64 编解码转换</p>
      </router-link>

      <router-link to="/number-converter" class="tool-card">
        <div class="tool-icon">🔢</div>
        <h3>进制转换</h3>
        <p>二进制、八进制、十进制、十六进制转换</p>
      </router-link>

      <router-link to="/text-diff" class="tool-card">
        <div class="tool-icon">📋</div>
        <h3>文本对比</h3>
        <p>文本差异对比工具</p>
      </router-link>

      <router-link to="/curl-converter" class="tool-card">
        <div class="tool-icon">🔄</div>
        <h3>cURL 转换</h3>
        <p>将 cURL 命令转换为各种编程语言代码</p>
      </router-link>

      <router-link to="/unicode-converter" class="tool-card">
        <div class="tool-icon">🔤</div>
        <h3>Unicode 转换</h3>
        <p>Unicode 编码解码转换工具</p>
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
