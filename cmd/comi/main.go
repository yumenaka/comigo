//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest
package main

import (
	"os"

	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/routers"
)

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
	cmd.Execute()
	// 启动网页服务器（不阻塞）
	routers.StartWebServer()
	// 启动或停止 Tailscale 服务（如启用）
	routers.StartTailscale()
	// 扫描书库（命令行指定）
	cmd.ScanStore(cmd.Args)
	// 在命令行显示QRCode
	cmd.ShowQRCode()
	// 退出时清理临时文件的处理函数
	cmd.SetShutdownHandler()
}
