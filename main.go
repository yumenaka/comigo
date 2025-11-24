//go:build !js

//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest
package main

import (
	"os"

	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/routers"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
)

// 运行 Comigo 服务器
func main() {
	// 检查是否只是查看版本或帮助信息
	for _, arg := range os.Args {
		if arg == "-v" || arg == "--version" || arg == "-h" || arg == "--help" {
			// 初始化命令行flag与args，环境变量与配置文件
			cmd.Execute()
			// 打印信息后直接退出
			return
		}
	}

	// 初始化命令行flag与args，环境变量与配置文件
	// 需要在单实例检查之前执行，以便获取 cmd.Args
	cmd.Execute()

	// 处理新参数的回调函数（当已有实例运行时，新实例会调用此函数）
	handleNewArgs := func(args []string) error {
		if len(args) == 0 {
			return nil
		}
		logger.Infof("Received new args from another instance: %v", args)
		// 添加新的扫描路径
		cmd.CreateStoreUrls(args)
		// 扫描新添加的书库
		cmd.ScanStore()
		// 保存书籍元数据
		cmd.SaveMetadata()
		return nil
	}

	// 确保单实例模式运行
	isFirstInstance, err := tools.EnsureSingleInstance(cmd.Args, handleNewArgs)
	if err != nil {
		logger.Infof("Single instance check failed: %v", err)
		// 如果单实例检查失败，仍然继续运行（向后兼容）
	} else if !isFirstInstance {
		// 已有实例运行，参数已发送，直接退出
		logger.Infof("Args sent to existing instance, exiting...")
		return
	}

	// 第一个实例，正常启动
	// 注册退出时清理单实例资源
	defer tools.CleanupSingleInstance()

	// 启动网页服务器（不阻塞）
	routers.StartWebServer()
	// 启动或停止 Tailscale 服务（如启用）
	routers.StartTailscale()
	// 生成书库URL
	cmd.CreateStoreUrls(cmd.Args)
	// 加载书籍元数据
	cmd.LoadMetadata()
	// 扫描书库（命令行指定）
	cmd.ScanStore()
	// 保存书籍元数据
	cmd.SaveMetadata()
	// 在命令行显示QRCode
	cmd.ShowQRCode()
	// 退出时清理临时文件的处理函数
	cmd.SetShutdownHandler()
}

// // tui实验
// func RunTui() {
// 	// 判断是否在终端中运行
// 	if term.IsTerminal(os.Stdout.Fd()) {
// 		// 1. 初始化自定义的日志缓冲区
// 		logBuffer := tui.NewLogBuffer()
// 		// 将标准日志的输出重定向到 logBuffer
// 		logger.SetOutput(logBuffer)

// 		// 2. 创建 Bubble Tea 程序的模型
// 		model := tui.InitialModel(logBuffer)
// 		// 创建一个bubbletea的应用对象
// 		program := tea.NewProgram(model)

// 		// Comigo 服务器的初始化(初始化 Comigo 命令行flag与args，环境变量与配置文件)
// 		cmd.Execute()
// 		// 启动网页服务器（不阻塞）
// 		routers.StartWebServer()
// 		// 扫描书库（命令行指定）
// 		cmd.ScanStore(cmd.Args)

// 		// 3. 调用 Bubble Tea 对象的Start()方法开始执行，运行 TUI 程序
// 		if _, err := program.Run(); err != nil {
// 			logger.Errorf("Error running tui interface: %v", err)
// 		}
// 	} else {
// 		// 初始化命令行flag与args，环境变量与配置文件
// 		cmd.Execute()
// 		// 启动网页服务器（不阻塞）
// 		routers.StartWebServer()
// 		// 扫描书库（命令行指定）
// 		cmd.ScanStore(cmd.Args)
// 		// 在命令行显示QRCode
// 		cmd.ShowQRCode()
// 		// 退出时清理临时文件的处理函数
// 		cmd.SetShutdownHandler()
// 	}
// }
