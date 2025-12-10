//go:build !js

package system_tray

import (
	"embed"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/atotto/clipboard"
	"github.com/energye/systray"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/windows_registry"
)

// Sample：https://github.com/energye/systray/blob/main/example/main.go

//go:embed icon.ico
var iconData embed.FS

var (
	startServerFunc         func()
	shutdownServerFunc      func()
	getURLFunc              func() string
	getConfigDirFunc        func() (string, error)
	getStoreUrlsFunc        func() []string
	toggleTailscaleFunc     func() error
	setLanguageFunc         func(string) error
	getTailscaleEnabledFunc func() bool
	// 菜单项引用，用于语言切换时更新
	menuItems struct {
		mOpenBrowser           *systray.MenuItem
		mCopyURL               *systray.MenuItem
		mTailscale             *systray.MenuItem
		mExtra                 *systray.MenuItem
		mProject               *systray.MenuItem
		mContextFolder         *systray.MenuItem
		mContextFileAssoc      *systray.MenuItem
		mCreateDesktopShortcut *systray.MenuItem
		mLanguage              *systray.MenuItem
		mLangZh                *systray.MenuItem
		mLangEn                *systray.MenuItem
		mLangJa                *systray.MenuItem
		mOpenDir               *systray.MenuItem
		mConfigDir             *systray.MenuItem
		mStoreFolders          []*systray.MenuItem
		mQuit                  *systray.MenuItem
	}
)

// SetupSystray 设置系统托盘
// startServer: 启动服务器的函数
// shutdownServer: 清理服务器的函数
// getURL: 获取服务器URL的函数
// getConfigDir: 获取配置目录的函数
// getStoreUrls: 获取书库URL列表的函数
// toggleTailscale: 切换Tailscale状态的函数
// setLanguage: 设置语言的函数
// getTailscaleEnabled: 获取Tailscale是否启用的函数
func SetupSystray(
	startServer, shutdownServer func(),
	getURL func() string,
	getConfigDir func() (string, error),
	getStoreUrls func() []string,
	toggleTailscale func() error,
	setLanguage func(string) error,
	getTailscaleEnabled func() bool,
) {
	startServerFunc = startServer
	shutdownServerFunc = shutdownServer
	getURLFunc = getURL
	getConfigDirFunc = getConfigDir
	getStoreUrlsFunc = getStoreUrls
	toggleTailscaleFunc = toggleTailscale
	setLanguageFunc = setLanguage
	getTailscaleEnabledFunc = getTailscaleEnabled

	// 在主线程运行 systray
	systray.Run(onReady, onExit)
}

// onReady 系统托盘就绪时的回调
func onReady() {
	// 从嵌入的文件系统中读取图标
	iconBytes, err := iconData.ReadFile("icon.ico")
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_read_icon_file"), err)
		// 如果读取失败，使用默认图标
		systray.SetIcon(nil)
	} else {
		systray.SetIcon(iconBytes)
	}
	// 设置托盘图标旁边的文字（占用空间太大，注释掉，以后或许可以显示用户数什么的）
	// systray.SetTitle(“Comigo”)

	OnClickTray := func(menu systray.IMenu) {
		// 清理所有菜单项
		systray.ResetMenu()
		// 重新创建所有菜单项（使用最新的语言和书库链接）
		initMenuItems()
		// 显示菜单
		if menu != nil { // menu for linux nil
			menu.ShowMenu()
		}
	}
	// 设置单击托盘图标时的回调
	systray.SetOnClick(OnClickTray)
	// 设置右击托盘图标时的回调
	systray.SetOnRClick(OnClickTray)
	// 初始化菜单项
	initMenuItems()

	// 在后台启动Comigo服务
	go func() {
		if startServerFunc != nil {
			startServerFunc()
		}
	}()
}

// initMenuItems 初始化菜单项
func initMenuItems() {
	// 设置托盘工具提示
	systray.SetTooltip(locale.GetString("systray_tooltip"))

	// 创建菜单项
	menuItems.mOpenBrowser = systray.AddMenuItem(locale.GetString("systray_open_browser"), locale.GetString("systray_open_browser_tooltip"))
	menuItems.mOpenBrowser.Click(func() {
		if getURLFunc != nil {
			url := getURLFunc()
			go tools.OpenBrowserByURL(url)
			logger.Infof(locale.GetString("log_opening_browser"), url)
		}
	})

	// 复制阅读地址
	menuItems.mCopyURL = systray.AddMenuItem(locale.GetString("systray_copy_url"), locale.GetString("systray_copy_url_tooltip"))
	menuItems.mCopyURL.Click(func() {
		if getURLFunc != nil {
			url := getURLFunc()
			if err := clipboard.WriteAll(url); err != nil {
				logger.Infof(locale.GetString("log_failed_to_copy_url"), err)
			} else {
				logger.Infof(locale.GetString("log_copied_url_to_clipboard"), url)
			}
		}
	})

	// Tailscale 开关
	if getTailscaleEnabledFunc != nil {
		tailscaleEnabled := getTailscaleEnabledFunc()
		tailscaleTitle := locale.GetString("systray_enable_tailscale")
		if tailscaleEnabled {
			tailscaleTitle = locale.GetString("systray_disable_tailscale")
		}
		menuItems.mTailscale = systray.AddMenuItem(tailscaleTitle, locale.GetString("systray_toggle_tailscale_tooltip"))
		menuItems.mTailscale.Click(func() {
			if toggleTailscaleFunc != nil {
				if err := toggleTailscaleFunc(); err != nil {
					logger.Infof(locale.GetString("log_failed_to_toggle_tailscale"), err)
				}
				// 菜单会在下次点击托盘图标时自动更新，这里不需要手动更新
			}
		})
	}

	// 语言切换子菜单
	menuItems.mLanguage = systray.AddMenuItem(locale.GetString("systray_language"), locale.GetString("systray_language_tooltip"))
	menuItems.mLangZh = menuItems.mLanguage.AddSubMenuItem(locale.GetString("systray_language_zh"), locale.GetString("systray_language_zh_tooltip"))
	menuItems.mLangEn = menuItems.mLanguage.AddSubMenuItem(locale.GetString("systray_language_en"), locale.GetString("systray_language_en_tooltip"))
	menuItems.mLangJa = menuItems.mLanguage.AddSubMenuItem(locale.GetString("systray_language_ja"), locale.GetString("systray_language_ja_tooltip"))

	menuItems.mLangZh.Click(func() {
		if setLanguageFunc != nil {
			if err := setLanguageFunc("zh-CN"); err != nil {
				logger.Infof(locale.GetString("log_failed_to_set_language"), err)
			} else {
				logger.Info(locale.GetString("log_language_changed_to_chinese"))
				// 菜单会在下次点击托盘图标时自动更新，这里不需要手动更新
			}
		}
	})
	menuItems.mLangEn.Click(func() {
		if setLanguageFunc != nil {
			if err := setLanguageFunc("en-US"); err != nil {
				logger.Infof(locale.GetString("log_failed_to_set_language"), err)
			} else {
				logger.Info(locale.GetString("log_language_changed_to_english"))
				// 菜单会在下次点击托盘图标时自动更新，这里不需要手动更新
			}
		}
	})
	menuItems.mLangJa.Click(func() {
		if setLanguageFunc != nil {
			if err := setLanguageFunc("ja-JP"); err != nil {
				logger.Infof(locale.GetString("log_failed_to_set_language"), err)
			} else {
				logger.Info(locale.GetString("log_language_changed_to_japanese"))
				// 菜单会在下次点击托盘图标时自动更新，这里不需要手动更新
			}
		}
	})

	// 打开目录子菜单
	menuItems.mOpenDir = systray.AddMenuItem(locale.GetString("systray_open_directory"), locale.GetString("systray_open_directory_tooltip"))

	// 配置文件目录
	if getConfigDirFunc != nil {
		configDir, err := getConfigDirFunc()
		if err == nil && configDir != "" {
			menuItems.mConfigDir = menuItems.mOpenDir.AddSubMenuItem(locale.GetString("systray_config_directory"), locale.GetString("systray_config_directory_tooltip"))
			menuItems.mConfigDir.Click(func() {
				if getConfigDirFunc != nil {
					configDir, err := getConfigDirFunc()
					if err != nil {
						logger.Infof(locale.GetString("log_failed_to_get_config_dir"), err)
						return
					}
					openDirectory(configDir)
				}
			})
		}
	}

	// 书库文件夹
	menuItems.mStoreFolders = make([]*systray.MenuItem, 0)
	if getStoreUrlsFunc != nil {
		storeUrls := getStoreUrlsFunc()
		for _, storeUrl := range storeUrls {
			if storeUrl == "" {
				continue
			}
			menuTitle := fmt.Sprintf("%s", storeUrl)
			mStore := menuItems.mOpenDir.AddSubMenuItem(menuTitle, storeUrl)
			menuItems.mStoreFolders = append(menuItems.mStoreFolders, mStore)
			// 避免闭包问题，在循环内部创建函数并立即调用
			func(url string) {
				mStore.Click(func() {
					openDirectory(url)
				})
			}(storeUrl)
		}
	}

	// 顶层“其他/Extra/その他”菜单（所有平台都显示）
	menuItems.mExtra = systray.AddMenuItem(locale.GetString("systray_extra"), locale.GetString("systray_extra_tooltip"))
	// Windows 右键菜单注册/清理（仅在 Windows 下显示）
	if runtime.GOOS == "windows" {
		// 子菜单：文件夹右键菜单注册/清理
		folderTitle := locale.GetString("register_folder_context_menu")
		if windows_registry.HasComigoFolderContextMenu() {
			folderTitle = locale.GetString("unregister_folder_context_menu")
		}
		menuItems.mContextFolder = menuItems.mExtra.AddSubMenuItem(folderTitle, folderTitle)
		menuItems.mContextFolder.Click(func() {
			if windows_registry.HasComigoFolderContextMenu() {
				if err := windows_registry.RemoveComigoFromFolderContextMenu(); err != nil {
					logger.Infof(locale.GetString("log_failed_to_clear_folder_context_menu"), err)
				} else {
					logger.Infof("%s", locale.GetString("unregister_folder_context_menu"))
				}
			} else {
				if err := windows_registry.AddComigoToFolderContextMenu(); err != nil {
					logger.Infof(locale.GetString("log_failed_to_register_folder_context_menu"), err)
				} else {
					logger.Infof("%s", locale.GetString("register_folder_context_menu"))
				}
			}
			// 文本更新依赖下次点击托盘图标时重新构建菜单
		})

		// 子菜单：在桌面创建快捷方式
		menuItems.mCreateDesktopShortcut = menuItems.mExtra.AddSubMenuItem(locale.GetString("create_desktop_shortcut"), locale.GetString("create_desktop_shortcut"))
		menuItems.mCreateDesktopShortcut.Click(func() {
			if err := windows_registry.CreateDesktopShortcut(); err != nil {
				logger.Infof(locale.GetString("log_failed_to_create_desktop_shortcut"), err)
			} else {
				logger.Infof("%s", locale.GetString("create_desktop_shortcut"))
			}
		})

		// 子菜单：文件类型关联注册/清理（候选打开方式）
		fileAssocTitle := locale.GetString("register_file_association")
		if windows_registry.HasComigoArchiveAssociation(nil) {
			fileAssocTitle = locale.GetString("unregister_file_association")
		}
		menuItems.mContextFileAssoc = menuItems.mExtra.AddSubMenuItem(fileAssocTitle, fileAssocTitle)
		menuItems.mContextFileAssoc.Click(func() {
			if windows_registry.HasComigoArchiveAssociation(nil) {
				if err := windows_registry.UnregisterComigoAsDefaultArchiveHandler(nil); err != nil {
					logger.Infof(locale.GetString("log_failed_to_unregister_archive_handler"), err)
				} else {
					logger.Infof("%s", locale.GetString("unregister_file_association"))
				}
			} else {
				if err := windows_registry.RegisterComigoAsDefaultArchiveHandler(nil); err != nil {
					logger.Infof(locale.GetString("log_failed_to_register_archive_handler"), err)
				} else {
					logger.Infof("%s", locale.GetString("register_file_association"))
				}
			}
			// 文本更新依赖下次点击托盘图标时重新构建菜单
		})
	}
	// Comigo 项目地址子菜单（所有平台都显示）
	menuItems.mProject = menuItems.mExtra.AddSubMenuItem(locale.GetString("systray_project"), locale.GetString("systray_project_tooltip"))
	menuItems.mProject.Click(func() {
		go tools.OpenBrowserByURL("https://github.com/yumenaka/comigo")
		logger.Infof(locale.GetString("log_opening_comigo_project_page"))
	})

	// 退出
	menuItems.mQuit = systray.AddMenuItem(locale.GetString("systray_quit"), locale.GetString("systray_quit_tooltip"))
	menuItems.mQuit.Click(func() {
		logger.Info(locale.GetString("log_requesting_quit_from_systray"))
		systray.Quit()
	})
}

// onExit 系统托盘退出时的回调
func onExit() {
	// 执行清理逻辑
	if shutdownServerFunc != nil {
		shutdownServerFunc()
	}
}

// openDirectory 打开指定目录
func openDirectory(path string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default: // linux
		cmd = exec.Command("xdg-open", path)
	}
	if err := cmd.Start(); err != nil {
		logger.Infof(locale.GetString("log_failed_to_open_directory"), err)
	}
}
