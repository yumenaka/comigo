//go:build !js

package tools

import (
	"embed"
	"fmt"

	"github.com/energye/systray"
	"github.com/yumenaka/comigo/tools/logger"
)

// Sample：https://github.com/energye/systray/blob/main/example/main.go

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
	// 设置托盘图标标题（占用空间太大，注释掉）
	// systray.SetTitle(“Comigo”)

	// 设置托盘工具提示
	systray.SetTooltip("Comigo 漫画阅读器")
	// 单击托盘图标时的回调
	systray.SetOnClick(func(menu systray.IMenu) {
		if menu != nil { // menu for linux nil
			menu.ShowMenu()
		}
		fmt.Println("SetOnClick")
	})
	//// 双击托盘图标时的回调
	//systray.SetOnDClick(func(menu systray.IMenu) {
	//	if menu != nil { // menu for linux nil
	//		menu.ShowMenu()
	//	}
	//	fmt.Println("SetOnDClick")
	//})

	// 创建菜单项
	mOpenBrowser := systray.AddMenuItem("打开浏览器", "在浏览器中打开 Comigo")
	//mOpenBrowser.Enable()
	mOpenBrowser.Click(func() {
		// 打开浏览器
		if getURLFunc != nil {
			url := getURLFunc()
			go OpenBrowser(url)
			logger.Infof("Opening browser: %s", url)
		}
	})
	mQuit := systray.AddMenuItem("退出", "退出 Comigo")

	//mQuit.Enable()
	mQuit.Click(func() {
		fmt.Println("Requesting quit")
		systray.Quit()
		//systray.Quit()// macos error
		//end() // macos error
		fmt.Println("Finished quitting")
	})
	// 在后台启动Comigo服务
	go func() {
		if startServerFunc != nil {
			startServerFunc()
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
