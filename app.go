package main

import (
	"fmt"
)

// App 暴露给 Wails 绑定的最小服务。
type App struct{}

// NewApp 创建 Wails 绑定服务。
func NewApp() *App {
	return &App{}
}

// Greet 保留 Wails 示例绑定，便于生成绑定时确认服务可用。
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
