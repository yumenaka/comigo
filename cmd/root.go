package cmd

import (
	"os"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/windows_registry"
)

var Args []string

// RootCmd 没有任何子命令的情况下时的基本命令
var RootCmd = &cobra.Command{
	Use:     locale.GetString("comigo_use"),
	Short:   locale.GetString("short_description"),
	Example: locale.GetString("comigo_example"),
	Version: config.GetVersion(),
	Long:    locale.GetString("long_description"),
	// Run 函数按以下顺序执行：PersistentPreRun() PreRun() Run() PostRun() PersistentPostRun()
	// 所有函数都能拿到相同的参数，即命令名称后面添加的参数。仅当设置了 Run 函数时，才会执行 PreRun 和 PostRun 函数。
	// 因为参数设置已完成，实际运行的命令习惯写在这里
	Run: func(cmd *cobra.Command, args []string) {
		Args = args
		// 通过“可执行文件名”设置部分默认参数
		SetByExecutableFilename()
		cfg := config.GetCfg()
		// 默认启用几个内置插件
		if cfg.EnablePlugin {
			cfg.EnabledPluginList = []string{"auto_flip", "auto_scroll"}
		}
		// 设置临时文件夹
		config.AutoSetCacheDir()
		// 在 Windows 上，根据命令行参数注册/卸载资源管理器文件夹右键菜单
		if runtime.GOOS == "windows" {
			// 先处理卸载，再处理注册，避免同时传入两个参数时出现冲突
			if cfg.UnregisterContextMenu {
				if err := windows_registry.RemoveComigoFromFolderContextMenu(); err != nil {
					logger.Infof(locale.GetString("log_failed_to_unregister_windows_context_menu"), err)
				} else {
					logger.Infof("%s", locale.GetString("unregister_context_menu"))
				}
			}
			if cfg.RegisterContextMenu {
				if err := windows_registry.AddComigoToFolderContextMenu(); err != nil {
					logger.Infof(locale.GetString("log_failed_to_register_windows_context_menu"), err)
				} else {
					logger.Infof("%s", locale.GetString("register_context_menu"))
				}
			}
		}
	},
}

// Execute 将所有子命令添加到根命令并适当设置标志。 由 main.main() 调用。 rootCmd 只需要执行一次。
func Execute() {
	// 初始化命令行参数。必须在config包初始化之前调用，因为config包会用到命令行参数 --config 指定的配置文件路径
	InitFlags()
	// 初始化配置文件
	cobra.OnInitialize(LoadConfigFile) // "OnInitialize"传入的函数，应该会在所有命令执行之前、包括rootCmd.Run之前执行。
	// 执行命令
	if err := RootCmd.Execute(); err != nil {
		logger.Infof("%s", err)
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}
}
