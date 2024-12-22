package main

import (
	"fmt"
	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/htmx/router"
	"golang.org/x/term"
	"log/slog"
	"os"
)

func main() {
	// 检测标准输出是否为终端
	if term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Println("程序正在终端中执行")
		// 参考这里的例子，优化操作：
		// https://github.com/charmbracelet/bubbletea/tree/main/examples
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
