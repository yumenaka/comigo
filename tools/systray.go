//go:build !js

package tools

import (
	"embed"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/atotto/clipboard"
	"github.com/energye/systray"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/logger"
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
		mOpenBrowser  *systray.MenuItem
		mCopyURL      *systray.MenuItem
		mTailscale    *systray.MenuItem
		mLanguage     *systray.MenuItem
		mLangZh       *systray.MenuItem
		mLangEn       *systray.MenuItem
		mLangJa       *systray.MenuItem
		mOpenDir      *systray.MenuItem
		mConfigDir    *systray.MenuItem
		mRefreshDir   *systray.MenuItem
		mStoreFolders []*systray.MenuItem
		mQuit         *systray.MenuItem
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
		logger.Infof("Failed to read icon file: %v, using default icon", err)
		// 如果读取失败，使用默认图标
		systray.SetIcon(nil)
	} else {
		systray.SetIcon(iconBytes)
	}
	// 设置托盘图标旁边的文字（占用空间太大，注释掉，以后或许可以显示用户数什么的）
	// systray.SetTitle(“Comigo”)

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
	// 单击托盘图标时的回调
	systray.SetOnClick(func(menu systray.IMenu) {
		if menu != nil { // menu for linux nil
			menu.ShowMenu()
		}
	})

	// 创建菜单项
	menuItems.mOpenBrowser = systray.AddMenuItem(locale.GetString("systray_open_browser"), locale.GetString("systray_open_browser_tooltip"))
	menuItems.mOpenBrowser.Click(func() {
		if getURLFunc != nil {
			url := getURLFunc()
			go OpenBrowser(url)
			logger.Infof("Opening browser: %s", url)
		}
	})

	// 复制阅读地址
	menuItems.mCopyURL = systray.AddMenuItem(locale.GetString("systray_copy_url"), locale.GetString("systray_copy_url_tooltip"))
	menuItems.mCopyURL.Click(func() {
		if getURLFunc != nil {
			url := getURLFunc()
			if err := clipboard.WriteAll(url); err != nil {
				logger.Infof("Failed to copy URL to clipboard: %v", err)
			} else {
				logger.Infof("Copied URL to clipboard: %s", url)
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
					logger.Infof("Failed to toggle Tailscale: %v", err)
				} else {
					// 更新菜单标题
					updateMenuTitles()
				}
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
				logger.Infof("Failed to set language: %v", err)
			} else {
				logger.Info("Language changed to Chinese")
				updateMenuTitles()
			}
		}
	})
	menuItems.mLangEn.Click(func() {
		if setLanguageFunc != nil {
			if err := setLanguageFunc("en-US"); err != nil {
				logger.Infof("Failed to set language: %v", err)
			} else {
				logger.Info("Language changed to English")
				updateMenuTitles()
			}
		}
	})
	menuItems.mLangJa.Click(func() {
		if setLanguageFunc != nil {
			if err := setLanguageFunc("ja-JP"); err != nil {
				logger.Infof("Failed to set language: %v", err)
			} else {
				logger.Info("Language changed to Japanese")
				updateMenuTitles()
			}
		}
	})

	// 打开目录子菜单
	menuItems.mOpenDir = systray.AddMenuItem(locale.GetString("systray_open_directory"), locale.GetString("systray_open_directory_tooltip"))

	// 初始化目录子菜单（配置目录和书库文件夹）
	refreshOpenDirSubMenu()

	// 刷新按钮（配置目录和书库文件夹后）
	menuItems.mRefreshDir = menuItems.mOpenDir.AddSubMenuItem(locale.GetString("systray_refresh_directories"), locale.GetString("systray_refresh_directories_tooltip"))
	menuItems.mRefreshDir.Click(func() {
		logger.Info("Refreshing directory menu...")
		refreshOpenDirSubMenu()
	})

	// 退出
	menuItems.mQuit = systray.AddMenuItem(locale.GetString("systray_quit"), locale.GetString("systray_quit_tooltip"))
	menuItems.mQuit.Click(func() {
		logger.Info("Requesting quit from system tray")
		systray.Quit()
	})
}

// updateMenuTitles 更新所有菜单项的标题
func updateMenuTitles() {
	// 更新托盘工具提示
	systray.SetTooltip(locale.GetString("systray_tooltip"))

	// 更新主菜单项
	if menuItems.mOpenBrowser != nil {
		menuItems.mOpenBrowser.SetTitle(locale.GetString("systray_open_browser"))
		menuItems.mOpenBrowser.SetTooltip(locale.GetString("systray_open_browser"))
	}
	if menuItems.mCopyURL != nil {
		menuItems.mCopyURL.SetTitle(locale.GetString("systray_copy_url"))
		menuItems.mCopyURL.SetTooltip(locale.GetString("systray_copy_url_tooltip"))
	}

	// 更新 Tailscale 菜单项
	if menuItems.mTailscale != nil && getTailscaleEnabledFunc != nil {
		tailscaleEnabled := getTailscaleEnabledFunc()
		if tailscaleEnabled {
			menuItems.mTailscale.SetTitle(locale.GetString("systray_disable_tailscale"))
			menuItems.mTailscale.SetTooltip(locale.GetString("systray_disable_tailscale"))
		} else {
			menuItems.mTailscale.SetTitle(locale.GetString("systray_enable_tailscale"))
			menuItems.mTailscale.SetTooltip(locale.GetString("systray_enable_tailscale"))
		}
	}

	// 更新语言切换菜单
	if menuItems.mLanguage != nil {
		menuItems.mLanguage.SetTitle(locale.GetString("systray_language"))
		menuItems.mLanguage.SetTooltip(locale.GetString("systray_language"))
	}
	if menuItems.mLangZh != nil {
		menuItems.mLangZh.SetTitle(locale.GetString("systray_language_zh"))
		menuItems.mLangZh.SetTooltip(locale.GetString("systray_language_zh"))
	}
	if menuItems.mLangEn != nil {
		menuItems.mLangEn.SetTitle(locale.GetString("systray_language_en"))
		menuItems.mLangEn.SetTooltip(locale.GetString("systray_language_en"))
	}
	if menuItems.mLangJa != nil {
		menuItems.mLangJa.SetTitle(locale.GetString("systray_language_ja"))
		menuItems.mLangJa.SetTooltip(locale.GetString("systray_language_ja"))
	}

	// 更新打开目录菜单
	if menuItems.mOpenDir != nil {
		menuItems.mOpenDir.SetTitle(locale.GetString("systray_open_directory"))
		menuItems.mOpenDir.SetTooltip(locale.GetString("systray_open_directory_tooltip"))
	}
	if menuItems.mConfigDir != nil {
		menuItems.mConfigDir.SetTitle(locale.GetString("systray_config_directory"))
		menuItems.mConfigDir.SetTooltip(locale.GetString("systray_config_directory_tooltip"))
	}
	if menuItems.mRefreshDir != nil {
		menuItems.mRefreshDir.SetTitle(locale.GetString("systray_refresh_directories"))
		menuItems.mRefreshDir.SetTooltip(locale.GetString("systray_refresh_directories_tooltip"))
	}
	// 书库文件夹与设置文件夹没有需要翻译的字段，只在点击刷新按钮后更新
	// 更新退出菜单
	if menuItems.mQuit != nil {
		menuItems.mQuit.SetTitle(locale.GetString("systray_quit"))
		menuItems.mQuit.SetTooltip(locale.GetString("systray_quit_tooltip"))
	}
}

// onExit 系统托盘退出时的回调
func onExit() {
	// 执行清理逻辑
	if shutdownServerFunc != nil {
		shutdownServerFunc()
	}
}

// refreshOpenDirSubMenu 刷新"打开目录"子菜单
// 为了确保刷新按钮始终在最后，我们需要先隐藏刷新按钮，刷新目录，然后重新创建刷新按钮
func refreshOpenDirSubMenu() {
	// 保存刷新按钮的引用（如果存在）
	refreshDirExists := menuItems.mRefreshDir != nil
	if refreshDirExists {
		menuItems.mRefreshDir.Hide()
		menuItems.mRefreshDir = nil
	}

	// 隐藏并清理旧的配置目录菜单项
	if menuItems.mConfigDir != nil {
		menuItems.mConfigDir.Hide()
		menuItems.mConfigDir = nil
	}

	// 隐藏并清理旧的书库文件夹菜单项
	for _, mStore := range menuItems.mStoreFolders {
		if mStore != nil {
			mStore.Hide()
		}
	}
	menuItems.mStoreFolders = make([]*systray.MenuItem, 0)

	// 重新创建配置文件目录菜单项
	if getConfigDirFunc != nil {
		configDir, err := getConfigDirFunc()
		if err == nil && configDir != "" {
			menuItems.mConfigDir = menuItems.mOpenDir.AddSubMenuItem(locale.GetString("systray_config_directory"), locale.GetString("systray_config_directory_tooltip"))
			menuItems.mConfigDir.Click(func() {
				if getConfigDirFunc != nil {
					configDir, err := getConfigDirFunc()
					if err != nil {
						logger.Infof("Failed to get config dir: %v", err)
						return
					}
					openDirectory(configDir)
				}
			})
		} else {
			logger.Infof("Config directory not available: %v", err)
		}
	}

	// 重新创建书库文件夹菜单项
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

	// 重新创建刷新按钮（放在最后）
	if refreshDirExists {
		menuItems.mRefreshDir = menuItems.mOpenDir.AddSubMenuItem(locale.GetString("systray_refresh_directories"), locale.GetString("systray_refresh_directories_tooltip"))
		menuItems.mRefreshDir.Click(func() {
			logger.Info("Refreshing directory menu...")
			refreshOpenDirSubMenu()
		})
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
		logger.Infof("Failed to open directory: %v", err)
	}
}
