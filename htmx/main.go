package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/htmx/router"
	"github.com/yumenaka/comigo/htmx/tui"
	"github.com/yumenaka/comigo/util/logger"
	"log/slog"
	"os"
)

func main() {
	// 检测标准输出是否为终端
	if term.IsTerminal(os.Stdout.Fd()) {
		fmt.Println("程序正在终端中执行")
		// 1. 初始化自定义的日志缓冲区
		logBuffer := tui.NewLogBuffer()
		// 将标准日志的输出重定向到 logBuffer
		logger.SetOutput(logBuffer)

		// 2. 创建 Bubble Tea 程序
		m := tui.InitialModel(logBuffer)
		p := tea.NewProgram(m)

		// 3. 运行 TUI 程序
		go func() {
			if _, err := p.Run(); err != nil {
				logger.Errorf("Error running tui interface: %v", err)
			}
		}()

	} else {
		fmt.Println("程序不在终端中执行")
	}
	// Initialize flags.
	cmd.InitFlags()
	// Run Comigo server.
	if err := router.RunServer(); err != nil {
		slog.Error("Failed to start server!", "details", err.Error())
		os.Exit(1)
	}
}
