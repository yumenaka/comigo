//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest

package main

import (
	"os"

	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/routers"
)

// 运行 Comigo 服务器
func main() {
	// 初始化命令行flag，环境变量与配置文件
	cmd.Execute()
	// 扫描书库（命令行指定）
	cmd.ScanStore(os.Args)
	// 启动网页服务器（不阻塞）
	routers.StartWebServer()
	// 在命令行显示QRCode
	cmd.ShowQRCode()
	// 退出时清理临时文件的处理函数
	cmd.SetShutdownHandler()
}

// // tui实验（TODO）
// func RunTui() {
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
// }
