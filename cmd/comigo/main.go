//go:generate goversioninfo -icon=../../icon.ico -manifest=../comi/goversioninfo.exe.manifest -internal-name=comigo-tray.exe -original-name=comigo-tray.exe ../comi/versioninfo.json
package main

import (
	"fmt"
	"os"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/system_tray"
)

// 运行 Comigo 服务器
func main() {
	config.UseTrayConfigProfile()
	// 与 CLI、Wails 共用命令解析，管理命令执行后不启动托盘或扫描书库。
	if handled, err := cmd.RunProcessCommand(os.Args[1:], os.Stdout); handled {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	// 初始化命令行参数与配置。
	cmd.Execute()
	// 托盘菜单创建前加入启动书库，首次打开目录菜单即可看到参数中的路径。
	cmd.AddStoreUrls(cmd.Args)
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
	cmd.Shutdown()
}
