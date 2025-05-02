package main

import "github.com/yumenaka/comigo/cmd"

func main() {
	cmd.Execute()
	// // 解析命令行参数
	// cmd.InitFlags()
	// // 初始化配置文件
	// cobra.OnInitialize(cmd.LoadConfigFile)
	// // 通过“可执行文件名”设置部分默认参数,目前不生效
	// config.SetByExecutableFilename()
	// // 设置临时文件夹
	// config.AutoSetCachePath()
	// // 扫描命令行指定的书库与文件
	// cmd.ScanStore(os.Args)
	// // 启动网页服务器（不阻塞）
	// routers.StartWebServer()
	// // 在命令行显示QRCode
	// cmd.ShowQRCode()
	// // 退出时清理临时文件
	// cmd.SetShutdownHandler()
}
