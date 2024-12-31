package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/htmx/router"
	"github.com/yumenaka/comigo/htmx/tui"
	"github.com/yumenaka/comigo/routers"
	"github.com/yumenaka/comigo/util/logger"
	"io"
	"log/slog"
	"os"
)

func main() {
	// 禁止Gin自带的控制台输出
	gin.DefaultWriter = io.Discard
	// TODO: mac 测试 TUI 界面，无法使用air热加载（无法绑定端口？）。需要使用 go run 运行。BUG原因未知。
	// go runTui()
	// 运行 Comigo 服务器
	runComigoServer()
}

func runTui() {
	if term.IsTerminal(os.Stdout.Fd()) {
		// 1. 初始化自定义的日志缓冲区
		logBuffer := tui.NewLogBuffer()
		// 将标准日志的输出重定向到 logBuffer
		logger.SetOutput(logBuffer)

		// 2. 创建 Bubble Tea 程序
		m := tui.InitialModel(logBuffer)
		p := tea.NewProgram(m)

		// 3. 运行 TUI 程序
		if _, err := p.Run(); err != nil {
			logger.Errorf("Error running tui interface: %v", err)
		}
	}
}

func runComigoServer() {
	// Initialize flags.
	cmd.InitFlags()
	routers.SetWebServerPort()
	// Run Comigo server.
	if err := router.RunServer(); err != nil {
		slog.Error("Failed to start server!", "details", err.Error())
		os.Exit(1)
	}
}
