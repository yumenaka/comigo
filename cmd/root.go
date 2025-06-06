package cmd

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/util/logger"
)

var Args []string

// RootCmd 没有任何子命令的情况下时的基本命令
var RootCmd = &cobra.Command{
	Use:     locale.GetString("comigo_use"),
	Short:   locale.GetString("short_description"),
	Example: locale.GetString("comigo_example"),
	Version: config.GetVersion(),
	Long:    locale.GetString("long_description"),
	// // Run 函数按以下顺序执行：PersistentPreRun() PreRun() Run() PostRun() PersistentPostRun()
	// // 所有函数都能拿到相同的参数，即命令名称后面添加的参数。
	// // 仅当设置了 Run 函数时，才会执行 PreRun 和 PostRun 函数。
	// 因为参数设置已完成，实际运行的命令习惯写在这里
	Run: func(cmd *cobra.Command, args []string) {
		Args = args
		// 通过“可执行文件名”设置部分默认参数,目前不生效
		config.SetByExecutableFilename()
		// 设置临时文件夹
		config.AutoSetCachePath()
	},
}

// LoadConfigFile 读取顺序：RAM（代码当中设定的默认值）+命令行参数  -> HomeDirectory -> ProgramDirectory -> WorkingDirectory
func LoadConfigFile() {
	// 在非js环境下
	if runtime.GOOS == "js" {
		return
	}
	home, err := os.UserHomeDir()
	if err != nil {
		logger.Infof("%s", err)
	}
	// 在HomeDir搜索配置
	homeConfigDir := path.Join(home, ".config/comigo")
	runtimeViper.AddConfigPath(homeConfigDir)
	// 在ProgramDirectory(二进制程序所在文件夹）的配置
	ProgramDirectory, err := os.Executable()
	if err != nil {
		logger.Infof("Failed to get ProgramDirectory:", err)
		return
	}
	// 将ProgramDirectory转换为绝对路径
	absPath, err := filepath.Abs(ProgramDirectory)
	if err != nil {
		logger.Infof("Failed to get absolute path:%s", err)
		return
	}
	logger.Infof("ProgramDirectory:%s", absPath)
	runtimeViper.AddConfigPath(absPath)

	// WorkingDirectory：当前执行目录
	WorkingDirectory, err := os.Getwd()
	if err != nil {
		logger.Infof("Failed to get WorkingDirectory:%s", err)
	}
	runtimeViper.AddConfigPath(WorkingDirectory)
	runtimeViper.SetConfigType("toml")
	runtimeViper.SetConfigName("config.toml")

	// 用户命令行指定的目录或文件
	if config.GetConfigPath() != "" {
		// SetConfigFile 显式定义配置文件的路径、名称和扩展名。 Viper 将使用它并且不检查任何配置路径。
		runtimeViper.SetConfigFile(config.GetConfigPath())
	}

	// 读取设定文件
	if err := runtimeViper.ReadInConfig(); err != nil {
		if config.GetConfigPath() == "" {
			logger.Infof("%s", err)
		}
	} else {
		// 获取当前使用的配置文件路径
		// https://github.com/spf13/viper/issues/89
		tempConfigPath := runtimeViper.ConfigFileUsed()
		logger.Infof(locale.GetString("found_config_file")+"%s", tempConfigPath)
	}
	// 把设定文件的内容，解析到构造体里面。
	if err := runtimeViper.Unmarshal(config.GetCfg()); err != nil {
		logger.Infof("%s", err)
		os.Exit(1)
	}
	// //监听文件修改
	// runtimeViper.WatchConfig()
	// //文件修改时，执行重载设置、服务重启的函数
	// runtimeViper.OnConfigChange(handlerConfigReload)
}

// Execute 将所有子命令添加到根命令并适当设置标志。
// 由 main.main() 调用的。 rootCmd 只需要执行一次。
func Execute() {
	// 初始化命令行参数。不能放在初始化配置文件之后。
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
