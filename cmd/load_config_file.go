package cmd

import (
	"os"
	"runtime"
	"strings"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools/logger"
)

// LoadConfigFile 读取顺序：RAM 默认值 + 命令行参数 -> HomeDirectory -> ProgramDirectory -> WorkingDirectory。
func LoadConfigFile() {
	// 在非js环境下
	if runtime.GOOS == "js" {
		return
	}

	if temporaryReaderMode(RootCmd.Flags().Args(), runtimeViper.GetString("ConfigFile"), config.GetCfg().TemporaryReaderMode) {
		config.GetCfg().TemporaryReaderMode = true
		if err := runtimeViper.Unmarshal(config.GetCfg()); err != nil {
			logger.Infof("%s", err)
			os.Exit(1)
		}
		config.GetCfg().TemporaryReaderMode = true
		locale.InitLanguageFromConfig(config.GetCfg().Language)
		return
	}

	runtimeViper.SetConfigType("toml")

	// 用户命令行指定的目录或文件
	usedConfigFile := ""
	if config.GetCfg().ConfigFile != "" {
		// SetConfigFile 显式定义配置文件的路径、名称和扩展名。 Viper 将使用它并且不检查任何配置路径。
		runtimeViper.SetConfigFile(config.GetCfg().ConfigFile)
		usedConfigFile = config.GetCfg().ConfigFile
	} else if _, configFile := config.FindConfigFile(); configFile != "" {
		runtimeViper.SetConfigFile(configFile)
		usedConfigFile = configFile
	}

	// 读取设定文件
	if usedConfigFile != "" {
		if err := runtimeViper.ReadInConfig(); err != nil {
			if config.GetCfg().ConfigFile == "" {
				logger.Infof("%s", err)
			}
		} else {
			usedConfigFile = runtimeViper.ConfigFileUsed()
			logger.Infof(locale.GetString("found_config_file")+"%s", usedConfigFile)
		}
	}
	if usedConfigFile != "" {
		// 获取当前使用的配置文件路径
		// https://github.com/spf13/viper/issues/89
		config.GetCfg().ConfigFile = usedConfigFile
	}
	// 把设定文件的内容，解析到构造体里面。
	if err := runtimeViper.Unmarshal(config.GetCfg()); err != nil {
		logger.Infof("%s", err)
		os.Exit(1)
	}
	if usedConfigFile != "" {
		config.GetCfg().ConfigFile = usedConfigFile
	}
	// 根据配置重新初始化语言设置
	locale.InitLanguageFromConfig(config.GetCfg().Language)
	// //监听文件修改
	// runtimeViper.WatchConfig()
	// //文件修改时，执行重载设置、服务重启的函数
	// runtimeViper.OnConfigChange(handlerConfigReload)
}

// temporaryReaderMode 判断是否为纯文件参数临时阅读模式。
// ponytail: 只识别已存在的本地普通文件；目录、远程 URL 和不存在路径继续走普通书库模式。
func temporaryReaderMode(args []string, configFile string, force bool) bool {
	if force {
		return true
	}
	if strings.TrimSpace(configFile) != "" || len(args) == 0 {
		return false
	}
	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil || !info.Mode().IsRegular() {
			return false
		}
	}
	return true
}
