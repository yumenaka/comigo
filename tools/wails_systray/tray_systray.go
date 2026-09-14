//go:build wails && !js && !bindings && !darwin

package wails_systray

import (
	"github.com/energye/systray"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
)

// startPlatform 在 Windows/Linux 上接入 energye/systray。
func startPlatform(t *Tray) func() {
	start, end := systray.RunWithExternalLoop(func() {
		systray.SetIcon(trayIcon)
		systray.SetTooltip(locale.GetString("systray_tooltip"))
		mShow := systray.AddMenuItem(locale.GetString("wails_systray_show"), locale.GetString("wails_systray_show_tooltip"))
		mShow.Click(t.showWindow)
		mCopyURL := systray.AddMenuItem(locale.GetString("systray_copy_url"), locale.GetString("systray_copy_url_tooltip"))
		mCopyURL.Click(copyReaderURL)
		mOpenDir := systray.AddMenuItem(locale.GetString("systray_open_directory"), locale.GetString("systray_open_directory_tooltip"))
		storeURLs := config.GetCfg().StoreUrls
		for _, storeURL := range storeURLs {
			mStore := mOpenDir.AddSubMenuItem(storeURL, storeURL)
			mStore.Click(func() { openStoreDirectory(storeURL) })
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
		mQuit.Click(t.Quit)
		// 左右键统一打开菜单；Linux 库拿不到菜单时保留原有恢复行为。
		showMenu := func(menu systray.IMenu) {
			if menu == nil {
				t.showWindow()
				return
			}
			_ = menu.ShowMenu()
		}
		systray.SetOnClick(showMenu)
		systray.SetOnRClick(showMenu)
	}, nil)
	start()
	return end
}

func setPlatformWindowVisible(bool) {
}

func quitPlatformFallback() {
	systray.Quit()
}
