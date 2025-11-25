//go:build !js

package tools

import (
	"embed"

	"github.com/getlantern/systray"
	"github.com/yumenaka/comigo/tools/logger"
)

//go:embed icon.ico
var iconData embed.FS

var (
	startServerFunc    func()
	shutdownServerFunc func()
	getURLFunc         func() string
)

// SetupSystray 设置系统托盘
// startServer: 启动服务器的函数
// shutdownServer: 清理服务器的函数
// getURL: 获取服务器URL的函数
func SetupSystray(startServer, shutdownServer func(), getURL func() string) {
	startServerFunc = startServer
	shutdownServerFunc = shutdownServer
	getURLFunc = getURL

	// 在主线程运行 systray
	systray.Run(onReady, onExit)
}

// onReady 系统托盘就绪时的回调
func onReady() {
	// 从嵌入的文件系统中读取图标
	iconBytes, err := iconData.ReadFile("icon.ico")
	if err != nil {
		logger.Infof("Failed to read icon file: %v, using default icon", err)
		// 如果读取失败，使用默认图标
		systray.SetIcon(nil)
	} else {
		systray.SetIcon(iconBytes)
	}

	// 设置托盘工具提示
	systray.SetTooltip("Comigo 漫画阅读器")

	// 创建菜单项
	mOpenBrowser := systray.AddMenuItem("打开浏览器", "在浏览器中打开 Comigo")
	mQuit := systray.AddMenuItem("退出", "退出 Comigo")

	// 在后台启动服务器
	go func() {
		if startServerFunc != nil {
			startServerFunc()
		}
	}()

	// 处理菜单点击事件
	go func() {
		for {
			select {
			case <-mOpenBrowser.ClickedCh:
				// 打开浏览器
				if getURLFunc != nil {
					url := getURLFunc()
					go OpenBrowser(url)
					logger.Infof("Opening browser: %s", url)
				}

			case <-mQuit.ClickedCh:
				// 退出程序
				logger.Info("Quit requested from system tray")
				systray.Quit()
			}
		}
	}()
}

// onExit 系统托盘退出时的回调
func onExit() {
	// 执行清理逻辑
	if shutdownServerFunc != nil {
		shutdownServerFunc()
	}
}
