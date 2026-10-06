//go:build wails && !js && !bindings && !darwin

package wails_systray

import (
	"fyne.io/systray"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
)

// startPlatform 在 Windows/Linux 上接入 原生托盘菜单。
func startPlatform(t *Tray) func() {
	start, end := systray.RunWithExternalLoop(func() {
		systray.SetIcon(trayIcon)
		systray.SetTooltip(locale.GetString("systray_tooltip"))
		mShow := systray.AddMenuItem(locale.GetString("wails_systray_show"), locale.GetString("wails_systray_show_tooltip"))
		onMenuClick(mShow, t.showWindow)
		mCopyURL := systray.AddMenuItem(locale.GetString("systray_copy_url"), locale.GetString("systray_copy_url_tooltip"))
		onMenuClick(mCopyURL, copyReaderURL)
		mOpenDir := systray.AddMenuItem(locale.GetString("systray_open_directory"), locale.GetString("systray_open_directory_tooltip"))
		storeURLs := config.GetCfg().StoreUrls
		for _, storeURL := range storeURLs {
			mStore := mOpenDir.AddSubMenuItem(storeURL, storeURL)
			onMenuClick(mStore, func() { openStoreDirectory(storeURL) })
		}
		if len(storeURLs) == 0 {
			mOpenDir.Disable()
		}
		systray.AddSeparator()
		// “其他”保持倒数第二，版本作为禁用子项；退出始终在最底部。
		mExtra := systray.AddMenuItem(locale.GetString("systray_extra"), locale.GetString("systray_extra_tooltip"))
		mVersion := mExtra.AddSubMenuItem("Comigo "+config.GetVersion(), "")
		mVersion.Disable()
		mQuit := systray.AddMenuItem(locale.GetString("systray_quit"), locale.GetString("systray_quit_tooltip"))
		onMenuClick(mQuit, t.Quit)
		// 不设置图标激活回调，由宿主左右键打开同一份菜单。
	}, nil)
	start()
	return end
}

func setPlatformWindowVisible(bool) {
}

func quitPlatformFallback() {
	systray.Quit()
}

// onMenuClick 在菜单移除时随点击通道关闭而结束监听。
func onMenuClick(item *systray.MenuItem, action func()) {
	go func() {
		for range item.ClickedCh {
			action()
		}
	}()
}
