# macOS DMG 相关变量
APP_NAME = dev-tools
DMG_NAME = $(APP_NAME)
APP_DIR = build/bin/$(APP_NAME).app
DMG_DIR = build/bin
DMG_FILE = $(DMG_DIR)/$(DMG_NAME).dmg

# 默认目标
.PHONY: all
all: build

# 构建应用
.PHONY: build
build:
	wails build -platform=darwin/universal

# 创建 DMG
.PHONY: dmg
dmg: build
	@echo "Creating DMG..."
	@test -d "$(APP_DIR)" || (echo "Error: $(APP_DIR) not found" && exit 1)
	@rm -f "$(DMG_FILE)"
	@create-dmg \
		--volname "$(DMG_NAME)" \
		--window-pos 200 120 \
		--window-size 800 400 \
		--icon-size 100 \
		--icon "$(APP_NAME).app" 200 190 \
		--hide-extension "$(APP_NAME).app" \
		--app-drop-link 600 185 \
		"$(DMG_FILE)" \
		"$(APP_DIR)"
	@echo "DMG created at $(DMG_FILE)"

# 清理构建文件
.PHONY: clean
clean:
	rm -rf build/bin/*

build-win:
	wails build -platform=windows/amd64

build-win-dev:
	wails build -platform=windows/amd64 -debug -devtools

build-mac:
	wails build -platform=darwin/universal

build-mac-dev:
	wails build -platform=darwin/universal -debug -devtools
