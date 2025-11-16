package cmd

import (
	"os"
	"path"
	"runtime"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools/logger"
)

// LoadConfigFile 读取顺序：RAM（代码当中设定的默认值）+命令行参数  -> HomeDirectory -> ProgramDirectory -> WorkingDirectory
func LoadConfigFile() {
	// 在非js环境下
	if runtime.GOOS == "js" {
		return
	}

	// 使用viper读取配置文件
	// 候选目录1：HomeDirectory
	home, err := os.UserHomeDir()
	if err == nil {
		configDir := path.Join(home, ".config/comigo")
		runtimeViper.AddConfigPath(configDir)
	} else {
		logger.Infof("Failed to get HomeDirectory:%s", err)
	}
	// 候选目录2：二进制程序所在文件夹
	// 可执行文件返回启动当前进程的可执行文件的路径名。无法保证路径仍然指向正确的可执行文件。如果使用符号链接来启动进程，则结果可能是符号链接或其指向的路径，具体取决于操作系统。
	// 如果需要稳定的结果，[pathfilepath.EvalSymlinks] 可能会有所帮助。除非发生错误，否则可执行文件将返回绝对路径。主要用例是查找与可执行文件相关的资源
	ProgramDirectory, err := os.Executable()
	if err == nil {
		logger.Infof("ProgramDirectory:%s", ProgramDirectory)
		runtimeViper.AddConfigPath(ProgramDirectory)
	} else {
		logger.Infof("Failed to get ProgramDirectory:", err)
	}

	// 候选目录3：当前执行目录
	WorkingDirectory, err := os.Getwd()
	if err == nil {
		runtimeViper.AddConfigPath(WorkingDirectory)
	} else {
		logger.Infof("Failed to get WorkingDirectory:%s", err)
	}

	runtimeViper.SetConfigType("toml")
	runtimeViper.SetConfigName("config.toml")

	// 用户命令行指定的目录或文件
	if config.GetCfg().ConfigFile != "" {
		// SetConfigFile 显式定义配置文件的路径、名称和扩展名。 Viper 将使用它并且不检查任何配置路径。
		runtimeViper.SetConfigFile(config.GetCfg().ConfigFile)
	}

	// 读取设定文件
	if err := runtimeViper.ReadInConfig(); err != nil {
		if config.GetCfg().ConfigFile == "" {
			logger.Infof("%s", err)
		}
	} else {
		// 获取当前使用的配置文件路径
		// https://github.com/spf13/viper/issues/89
		tempConfigDir := runtimeViper.ConfigFileUsed()
		logger.Infof(locale.GetString("found_config_file")+"%s", tempConfigDir)
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
