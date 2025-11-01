//go:build !js

//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest
package main

import (
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/cmd/tui"
	"github.com/yumenaka/comigo/routers"
	"github.com/yumenaka/comigo/tools/logger"
)

// 运行 Comigo 服务器
func main() {
	// Create an instance of the app structure
	app := NewApp()
	// 初始化命令行flag与args，环境变量与配置文件
	cmd.Execute()
	// 启动网页服务器（不阻塞）
	routers.StartWebServer()
	// 启动或停止 Tailscale 服务（如启用）
	routers.StartTailscale()
	// 扫描书库（命令行指定）
	cmd.ScanStore(cmd.Args)
	// 在命令行显示QRCode
	cmd.ShowQRCode()
	//// 退出时清理临时文件的处理函数
	//cmd.SetShutdownHandler()
	//// RunTui()

	// 获取网页服务器（echo）
	echo := routers.GetWebServer()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Comigo",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:  nil,
			Handler: http.HandlerFunc(echo.ServeHTTP),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}

// tui实验
func RunTui() {
	// 判断是否在终端中运行
	if term.IsTerminal(os.Stdout.Fd()) {
		// 1. 初始化自定义的日志缓冲区
		logBuffer := tui.NewLogBuffer()
		// 将标准日志的输出重定向到 logBuffer
		logger.SetOutput(logBuffer)

		// 2. 创建 Bubble Tea 程序的模型
		model := tui.InitialModel(logBuffer)
		// 创建一个bubbletea的应用对象
		program := tea.NewProgram(model)

		// Comigo 服务器的初始化(初始化 Comigo 命令行flag与args，环境变量与配置文件)
		cmd.Execute()
		// 启动网页服务器（不阻塞）
		routers.StartWebServer()
		// 扫描书库（命令行指定）
		cmd.ScanStore(cmd.Args)

		// 3. 调用 Bubble Tea 对象的Start()方法开始执行，运行 TUI 程序
		if _, err := program.Run(); err != nil {
			logger.Errorf("Error running tui interface: %v", err)
		}
	} else {
		// 初始化命令行flag与args，环境变量与配置文件
		cmd.Execute()
		// 启动网页服务器（不阻塞）
		routers.StartWebServer()
		// 扫描书库（命令行指定）
		cmd.ScanStore(cmd.Args)
		// 在命令行显示QRCode
		cmd.ShowQRCode()
		// 退出时清理临时文件的处理函数
		cmd.SetShutdownHandler()
	}
}
