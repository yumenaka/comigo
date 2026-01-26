//go:generate goversioninfo -icon=../../icon.ico -manifest=goversioninfo.exe.manifest versioninfo.json
package main

import (
	"os"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
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
	// 需要在单实例检查之前执行，以便获取 cmd.Args
	cmd.Execute()
	// 如果启用了单实例模式，进行单实例检查
	if config.GetCfg().EnableSingleInstance {
		// 处理新参数的回调函数（当已有实例运行时，新实例会调用此函数）
		handleNewArgs := func(args []string) error {
			if len(args) == 0 {
				return nil
			}
			logger.Infof(locale.GetString("log_received_new_args_from_instance"), args)
			// 添加新扫描路径
			cmd.AddStoreUrls(args)
			// 扫描新添加的书库
			cmd.ScanStore()
			// 保存书籍元数据
			cmd.SaveMetadata()
			// 生成书组
			model.GenerateBookGroup()
			// 判断是否需要打开浏览器
			config.OpenBrowserIfNeeded()
			return nil
		}

		// 确保单实例模式运行
		isFirstInstance, err := tools.EnsureSingleInstance(cmd.Args, handleNewArgs)
		if err != nil {
			logger.Infof(locale.GetString("log_single_instance_check_failed"), err)
			// 如果单实例检查失败，仍然继续运行（向后兼容）
		} else if !isFirstInstance {
			// 已有实例运行，参数已发送，直接退出
			logger.Infof(locale.GetString("log_args_sent_to_existing_instance"))
			return
		}

		// 第一个实例，正常启动
		// 注册退出时清理单实例资源
		defer tools.CleanupSingleInstance()
	}
	// 启动网页服务器（不阻塞）
	routers.StartWebServer()
	// 启动或停止 Tailscale 服务（如启用）
	routers.StartTailscale()
	// 分析命令行参数，生成书库URL
	cmd.AddStoreUrls(cmd.Args)
	// 如果没有指定扫描路径，就把当前工作目录作为扫描路径
	cmd.SetCwdAsScanPathIfNeed()
	// 加载书籍元数据（包括书签）
	cmd.LoadMetadata()
	// 扫描书库
	cmd.ScanStore()
	// 保存书籍元数据（包括书签）
	cmd.SaveMetadata()
	// 生成书组
	model.GenerateBookGroup()
	// 启动自动扫描（如果配置了间隔）
	config.StartOrStopAutoRescan()
	// 在命令行显示QRCode
	cmd.ShowQRCode()
	// 判断是否需要打开浏览器
	config.OpenBrowserIfNeeded()
	// 退出时清理临时文件的处理函数
	cmd.SetShutdownHandler()
}
