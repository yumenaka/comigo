//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest
package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/system_tray"
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
	// 如果启用了单实例模式，进行单实例检查
	if config.GetCfg().EnableSingleInstance {
		// 处理新参数的回调函数（当已有实例运行时，新实例会调用此函数）
		handleNewArgs := func(args []string) error {
			if len(args) == 0 {
				return nil
			}
			logger.Infof(locale.GetString("log_received_new_args_from_instance"), args)
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

	// 设置系统托盘并启动服务器
	system_tray.SetupSystray(
		startServer,
		shutdownServer,
		getServerURL,
		getConfigDir,
		getStoreUrls,
		toggleTailscale,
		setLanguage,
		getTailscaleEnabled,
	)
}

// startServer 启动服务器
func startServer() {
	// 启动网页服务器（不阻塞）
	routers.StartWebServer()
	// 启动或停止 Tailscale 服务（如启用）
	routers.StartTailscale()
	// 分析命令行参数，生成书库URL
	cmd.CreateStoreUrls(cmd.Args)
	// 加载书籍元数据（包括书签）
	cmd.LoadMetadata()
	// 扫描书库
	cmd.ScanStore()
	// 保存书籍元数据（包括书签）
	cmd.SaveMetadata()
	// 在命令行显示QRCode
	cmd.ShowQRCode()
}

// getServerURL 获取服务器URL
func getServerURL() string {
	return config.GetQrcodeURL()
}

// getConfigDir 获取配置目录
func getConfigDir() (string, error) {
	return config.GetConfigDir()
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

	// 清理临时文件
	if config.GetCfg().ClearCacheExit {
		logger.Infof("\r"+locale.GetString("start_clear_file")+" CacheDir:%s ", config.GetCfg().CacheDir)
		allBooks, err := model.IStore.ListBooks()
		if err != nil {
			logger.Infof(locale.GetString("log_error_listing_books"), err)
		}
		for _, book := range allBooks {
			// 清理某一本书的缓存
			cachePath := filepath.Join(config.GetCfg().CacheDir, book.BookID)
			err := os.RemoveAll(cachePath)
			if err != nil {
				logger.Infof(locale.GetString("log_error_clearing_temp_files"), cachePath)
			} else if config.GetCfg().Debug {
				logger.Infof(locale.GetString("log_cleared_temp_files"), cachePath)
			}
		}
		logger.Infof("%s", locale.GetString("clear_temp_file_completed"))
	}

	// 关闭服务器
	if config.Server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := config.Server.Shutdown(ctx); err != nil {
			log.Fatal("Comigo Server forced to shutdown: ", err)
		}
	}
	log.Println("Comigo Server exit.")
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
