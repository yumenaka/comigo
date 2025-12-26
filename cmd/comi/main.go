package main

import (
	"os"

	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/config"
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
	// 初始化命令行flag与args: 读取环境变量 -> 根据可执行文件名设置部分默认参数 -> 读取配置文件
	cmd.Execute()
	// 启动网页服务器（不阻塞）
	routers.StartWebServer()
	// 启动或停止 Tailscale 服务（如启用）
	routers.StartTailscale()
	// 分析命令行参数，生成书库URL
	cmd.AddStoreUrls(cmd.Args)
	// 如果没有指定扫描路径，就把当前工作目录作为扫描路径
	cmd.SetCwdAsScanPathIfNeed()
	// 设置上传路径
	cmd.SetUploadPath(cmd.Args)
	// 加载书籍元数据（包括书签）
	cmd.LoadMetadata()
	// 扫描书库
	cmd.ScanStore()
	// 保存书籍元数据（包括书签）
	cmd.SaveMetadata()
	// 启动自动扫描（如果配置了间隔）
	config.StartOrStopAutoRescan()
	// 在命令行显示QRCode
	cmd.ShowQRCode()
	// 退出时清理临时文件的处理函数
	cmd.SetShutdownHandler()
}
