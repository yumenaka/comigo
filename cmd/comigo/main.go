//go:generate goversioninfo -icon=../../icon.ico -manifest=goversioninfo.exe.manifest versioninfo.json
package main

import (
	"os"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/system_tray"
)

// 运行 Comigo 服务器
func main() {
	config.UseTrayConfigProfile()
	// 检查是否只是查看版本或帮助信息
	for _, arg := range os.Args {
		if arg == "-v" || arg == "--version" || arg == "-h" || arg == "--help" ||
			arg == "-u" || arg == "--upgrade" {
			// 初始化命令行flag与args，环境变量与配置文件
			cmd.Execute()
			// 打印信息后直接退出（含自升级流程内的 os.Exit）
			return
		}
	}
	// 初始化命令行flag与args，环境变量与配置文件
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
	var releaseSingleInstance func()
	if config.GetCfg().EnableSingleInstance {
		releaseSingleInstance = tools.CleanupSingleInstance
	}
	// 设置系统托盘并启动服务器
	exitCode := system_tray.SetupSystray(
		startServer,
		shutdownServer,
		getServerURL,
		getBrowserURL,
		config.GetConfigDir,
		getStoreUrls,
		toggleTailscale,
		setLanguage,
		getTailscaleEnabled,
		releaseSingleInstance,
	)
	if exitCode >= 0 {
		os.Exit(exitCode)
	}
}

// startServer 启动服务器
func startServer() {
	// 启动网页服务器（不阻塞）
	if err := routers.StartWebServer(); err != nil {
		logger.Infof("%v", err)
		return
	}
	// 启动或停止 Tailscale 服务（如启用）
	routers.StartTailscale()
	// 加载用户插件，与 CLI 和桌面入口保持一致。
	cmd.LoadUserPlugins()
	// 分析命令行参数，生成书库URL
	cmd.AddStoreUrls(cmd.Args)
	// 加载书籍元数据（包括书签）
	cmd.LoadMetadata()
	// 扫描书库
	cmd.ScanStore()
	// 保存书籍元数据（包括书签）
	cmd.SaveMetadata()
	// 启动自动扫描（如果配置了间隔）
	config.StartOrStopAutoRescan()
	// 在命令行显示 QRCode
	cmd.ShowQRCode()
	// 判断是否需要打开浏览器
	config.OpenBrowserIfNeeded()
}

// getServerURL 获取服务器URL
func getServerURL() string {
	return config.GetQrcodeURL()
}

// getBrowserURL 获取托盘打开浏览器使用的本机 URL，避免局域网地址或自定义 Host 在本机不可访问。
func getBrowserURL() string {
	return config.GetLocalBrowserURL()
}

// getStoreUrls 获取书库URL列表
func getStoreUrls() []string {
	return config.GetCfg().StoreUrls
}

// toggleTailscale 切换Tailscale状态
func toggleTailscale() error {
	cfg := config.GetCfg()
	cfg.EnableTailscale = !cfg.EnableTailscale

	if cfg.EnableTailscale {
		routers.StartTailscale()
	} else {
		routers.StopTailscale()
	}

	// 保存配置
	return config.SaveConfig(config.DefaultConfigLocation())
}

// setLanguage 设置语言
func setLanguage(lang string) error {
	return locale.SetLanguage(lang)
}

// getTailscaleEnabled 获取Tailscale是否启用
func getTailscaleEnabled() bool {
	return config.GetCfg().EnableTailscale
}

// shutdownServer 清理服务器资源
func shutdownServer() {
	// 清理单实例资源
	if config.GetCfg().EnableSingleInstance {
		tools.CleanupSingleInstance()
	}

	cmd.Shutdown()
}
