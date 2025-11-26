package window

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Controller 窗口控制器
type Controller struct {
	ctx context.Context
}

// NewController 创建窗口控制器
func NewController(ctx context.Context) *Controller {
	return &Controller{ctx: ctx}
}

// Minimize 最小化窗口
func (c *Controller) Minimize() {
	runtime.WindowMinimise(c.ctx)
}

// Close 关闭窗口
func (c *Controller) Close() {
	runtime.Quit(c.ctx)
}

// ToggleFullscreen 切换全屏状态
func (c *Controller) ToggleFullscreen() {
	if runtime.WindowIsFullscreen(c.ctx) {
		runtime.WindowUnfullscreen(c.ctx)
	} else {
		runtime.WindowFullscreen(c.ctx)
	}
}

// IsFullscreen 获取全屏状态
func (c *Controller) IsFullscreen() bool {
	return runtime.WindowIsFullscreen(c.ctx)
}

// ResetSize 重置窗口大小到默认值
func (c *Controller) ResetSize() error {
	// 设置窗口大小为默认值
	runtime.WindowSetSize(c.ctx, 1000, 700)

	// 居中窗口
	runtime.WindowCenter(c.ctx)

	// 如果窗口是最大化的，先取消最大化
	if runtime.WindowIsMaximised(c.ctx) {
		runtime.WindowUnmaximise(c.ctx)
	}

	return nil
}

// GetSize 获取窗口大小
func (c *Controller) GetSize() (width, height int) {
	return runtime.WindowGetSize(c.ctx)
}

// GetPosition 获取窗口位置
func (c *Controller) GetPosition() (x, y int) {
	return runtime.WindowGetPosition(c.ctx)
}

// IsMaximised 获取最大化状态
func (c *Controller) IsMaximised() bool {
	return runtime.WindowIsMaximised(c.ctx)
}
