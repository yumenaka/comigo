//go:build wails && !js

package wails_systray

import (
	"errors"
	"runtime"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/yumenaka/comigo/assets"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/windows_registry"
)

// Tray 由 Wails 管理 desktop 原生托盘与阅读窗口的生命周期。
type Tray struct {
	app      *application.App
	window   application.Window
	icon     *application.SystemTray
	menu     *application.Menu
	quitting atomic.Bool
}

// New 为 desktop 创建托盘，与阅读窗口共用 Wails 原生事件循环。
func New(app *application.App, window application.Window) *Tray {
	t := &Tray{app: app, window: window}
	if !traySupported() {
		return t
	}
	t.icon = app.SystemTray.New()
	t.icon.SetIcon(assets.GetImageData("favicon.png"))
	t.menu = application.NewMenu()
	t.refreshMenu()
	if window != nil {
		window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
			if t.HandleBeforeClose() {
				event.Cancel()
			}
		})
	}
	return t
}

// Available 表示当前桌面实际提供了托盘宿主。
func (t *Tray) Available() bool { return t.icon != nil }

// HandleBeforeClose 有托盘才隐藏窗口；无托盘的 Linux 桌面仍能正常关闭应用。
func (t *Tray) HandleBeforeClose() bool {
	if t.icon == nil || t.quitting.Load() {
		return false
	}
	t.window.Hide()
	setDockVisible(false)
	return true
}

// Show 恢复阅读窗口。
func (t *Tray) Show() {
	setDockVisible(true)
	t.window.Show()
	t.window.UnMinimise()
	t.window.Focus()
}

// Quit 不经过窗口关闭拦截，统一交由 Wails 执行退出清理。
func (t *Tray) Quit() { t.quitting.Store(true); t.app.Quit() }

// refreshMenu 在语言或配置变化后重建菜单，所有操作直接复用现有 Go 服务。
func (t *Tray) refreshMenu() {
	if t.icon == nil {
		return
	}
	t.menu.Clear()
	add := func(menu *application.Menu, key string, action func()) *application.MenuItem {
		return menu.Add(locale.GetString(key)).OnClick(func(*application.Context) { action() })
	}
	if t.window != nil {
		add(t.menu, "wails_systray_show", t.Show)
	}
	add(t.menu, "systray_open_browser", func() { tools.OpenBrowserByURL(config.GetLocalBrowserURL()) })
	add(t.menu, "systray_copy_url", func() {
		if !t.app.Clipboard.SetText(config.GetQrcodeURL()) {
			logger.Info(locale.GetString("log_failed_to_copy_url"))
		}
	})
	languages := t.menu.AddSubmenu(locale.GetString("systray_language"))
	for _, lang := range []struct{ key, code string }{{"systray_language_zh", "zh-CN"}, {"systray_language_en", "en-US"}, {"systray_language_ja", "ja-JP"}} {
		add(languages, lang.key, func() {
			if err := t.setLanguage(lang.code); err != nil {
				logger.Info(err)
				return
			}
		})
	}
	dirs := t.menu.AddSubmenu(locale.GetString("systray_open_directory"))
	add(dirs, "systray_config_directory", func() {
		path, err := config.GetConfigDir()
		if err == nil {
			err = tools.OpenURL(path)
		}
		if err != nil {
			logger.Info(err)
		}
	})
	for _, path := range config.GetCfg().StoreUrls {
		dirs.Add(tools.LibraryLabel(path)).OnClick(func(*application.Context) {
			if err := tools.OpenURL(path); err != nil {
				logger.Info(err)
			}
		})
	}
	extra := t.menu.AddSubmenu(locale.GetString("systray_extra"))
	if runtime.GOOS == "windows" {
		folderKey := "register_folder_context_menu"
		if windows_registry.HasComigoFolderContextMenu() {
			folderKey = "unregister_folder_context_menu"
		}
		add(extra, folderKey, func() {
			var err error
			if windows_registry.HasComigoFolderContextMenu() {
				err = windows_registry.RemoveComigoFromFolderContextMenu()
			} else {
				err = windows_registry.AddComigoToFolderContextMenu()
			}
			if err != nil {
				logger.Info(err)
			}
			t.refreshMenu()
		})
		add(extra, "create_desktop_shortcut", func() {
			if err := windows_registry.CreateDesktopShortcut(); err != nil {
				logger.Info(err)
			}
		})
		fileKey := "register_file_association"
		if windows_registry.HasComigoArchiveAssociation(nil) {
			fileKey = "unregister_file_association"
		}
		add(extra, fileKey, func() {
			var err error
			if windows_registry.HasComigoArchiveAssociation(nil) {
				err = windows_registry.UnregisterComigoAsDefaultArchiveHandler(nil)
			} else {
				err = windows_registry.RegisterComigoAsDefaultArchiveHandler(nil)
			}
			if err != nil {
				logger.Info(err)
			}
			t.refreshMenu()
		})
	}
	add(extra, "systray_project", func() { tools.OpenBrowserByURL("https://github.com/yumenaka/comigo") })
	extra.Add("Comigo " + config.GetVersion()).SetEnabled(false)
	add(t.menu, "systray_quit", t.Quit)
	t.icon.SetTooltip(locale.GetString("systray_tooltip"))
	t.icon.SetMenu(t.menu)
}

// setLanguage 先保存配置再更新界面，面板与原生菜单使用同一条修改路径。
func (t *Tray) setLanguage(language string) error {
	_, err := config.ChangeConfig(true, func(next *config.Config) error {
		if next.ReadOnlyMode {
			return errors.New(locale.GetString("err_config_locked"))
		}
		next.Language = language
		return nil
	})
	if err != nil {
		return err
	}
	locale.InitLanguageFromConfig(language)
	t.refreshMenu()
	return nil
}
