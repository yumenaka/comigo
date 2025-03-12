package main

import (
	"log/slog"
	"os"

	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/htmx/router"
	"github.com/yumenaka/comigo/routers"
)

func main() {
	// 运行 Comigo 服务器
	runComigoServer()
}

//func runTui() {
//	if term.IsTerminal(os.Stdout.Fd()) {
//		// 1. 初始化自定义的日志缓冲区
//		logBuffer := tui.NewLogBuffer()
//		// 将标准日志的输出重定向到 logBuffer
//		logger.SetOutput(logBuffer)
//
//		// 2. 创建 Bubble Tea 程序
//		m := tui.InitialModel(logBuffer)
//		p := tea.NewProgram(m)
//
//		// 3. 运行 TUI 程序
//		if _, err := p.Run(); err != nil {
//			logger.Errorf("Error running tui interface: %v", err)
//		}
//	}
//}

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
