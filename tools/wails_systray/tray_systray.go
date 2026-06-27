//go:build wails && !js && !bindings && !darwin && !android

package wails_systray

import (
	"github.com/energye/systray"
	"github.com/yumenaka/comigo/assets/locale"
)

// startPlatform 在 Windows/Linux 上接入 energye/systray。
func startPlatform(t *Tray) func() {
	start, end := systray.RunWithExternalLoop(func() {
		systray.SetIcon(trayIcon)
		systray.SetTooltip(locale.GetString("systray_tooltip"))
		mShow := systray.AddMenuItem(locale.GetString("wails_systray_show"), locale.GetString("wails_systray_show_tooltip"))
		mShow.Click(t.showWindow)
		systray.AddSeparator()
		mQuit := systray.AddMenuItem(locale.GetString("systray_quit"), locale.GetString("systray_quit_tooltip"))
		mQuit.Click(t.quit)
		systray.SetOnClick(func(systray.IMenu) {
			t.showWindow()
		})
		systray.SetOnRClick(func(menu systray.IMenu) {
			if menu != nil {
				_ = menu.ShowMenu()
			}
		})
	}, nil)
	start()
	return end
}

func setPlatformWindowVisible(bool) {
}

func quitPlatformFallback() {
	systray.Quit()
}
