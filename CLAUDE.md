# CLAUDE.md

此文件为在此代码库中工作的 Claude Code (claude.ai/code) 提供指导。

## 项目概述

这是一个基于 Wails v2 + Vue 3 + TypeScript 构建的桌面开发者工具应用。它提供了一系列常用的开发者实用工具，如 JSON/XML 编辑器、时间转换器、URL 工具、Base64 转换、二维码生成等。

## 架构

**后端 (Go):**
- `main.go` - 入口点，包含 Wails 应用配置
- `app.go` - 主应用结构体和基本方法
- `json_processor.go` - JSON 处理工具
- `xml_processor.go` - XML 处理工具
- 使用 Wails v2 作为桌面应用框架

**前端 (Vue 3):**
- `frontend/src/App.vue` - 主 Vue 应用组件
- `frontend/src/views/` - 各个工具视图 (JsonEditor, TimeConverter 等)
- `frontend/src/components/` - 可复用的 Vue 组件
- `frontend/src/stores/` - Pinia 状态管理
- `frontend/src/router/` - Vue Router 配置
- 使用 Element Plus 作为 UI 组件库
- 使用 Monaco Editor 进行代码编辑
- 使用 Vite 作为构建工具

## 常用命令

### 开发
```bash
# 启动开发服务器（热重载）
wails dev

# 安装前端依赖
cd frontend && npm install

# 安装 Go 依赖
go mod tidy
```

### 构建
```bash
# 为当前平台构建
wails build

# 为 macOS 通用版构建
make build-mac
# 或
wails build -platform=darwin/universal

# 为 Windows 构建
make build-win
# 或
wails build -platform=windows/amd64

# 带调试/开发工具构建
make build-mac-dev
make build-win-dev
```

### 前端开发
```bash
cd frontend

# 开发服务器
npm run dev

# 仅构建前端
npm run build

# 类型检查
vue-tsc --noEmit
```

### macOS 分发
```bash
# 创建 DMG 包
make dmg

# 清理构建产物
make clean
```

## 项目结构

- 后端 Go 文件位于根目录
- 前端 Vue 应用位于 `frontend/` 目录
- Wails 配置文件 `wails.json`
- 构建输出至 `build/bin/`
- 文档和图片位于 `docs/`

## 关键依赖

**Go:**
- `github.com/wailsapp/wails/v2` - 桌面应用框架
- `github.com/gopherjs/gopherjs` - JS 互操作工具

**前端:**
- Vue 3 + TypeScript
- Element Plus (UI 组件)
- Monaco Editor (代码编辑)
- Pinia (状态管理)
- 各种特定工具的实用库

## 开发注意事项

- 应用使用嵌入式资源 (`//go:embed all:frontend/dist`)
- 前端通过 Wails 绑定与 Go 后端通信
- 多个工具视图组织为独立的 Vue 组件
- 工具状态和偏好设置通过 Pinia 存储管理
- 应用支持拖拽重排序工具卡片
- 所有工具都支持实时处理和预览