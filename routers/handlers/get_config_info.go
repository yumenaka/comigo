package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/go-homedir"
	"github.com/yumenaka/comi/logger"
	"net/http"
	"os"
	"path"
)

type ConfigInfo struct {
	RAMConfigEnable          bool `json:"ram_config_enable"`
	HomeDirectoryEnable      bool `json:"home_directory_enable"`
	ExecutionDirectoryEnable bool `json:"execution_directory_enable"`
	ProgramDirectoryEnable   bool `json:"program_directory_enable"`
}

// GetConfig 获取json格式的当前配置
func GetConfigInfo(c *gin.Context) {
	ConfigInfo := ConfigInfo{}
	logger.Info("Check Config Path" + Directory)
	// Home 目录
	home, err := homedir.Dir()
	if err != nil {
		logger.Info("error: Failed find home directory")
	}
	if fileExists(path.Join(home, ".config/comigo/config.toml")) {
		ConfigInfo.HomeDirectoryEnable = true
	}
	//当前执行目录
	if fileExists("config.toml") {
		ConfigInfo.ExecutionDirectoryEnable = true
	}
	// 可执行程序自身的文件路径
	executablePath, err := os.Executable()
	if err != nil {
		return
	}
	if fileExists(path.Join(executablePath, "config.toml")) {
		ConfigInfo.ProgramDirectoryEnable = true
	}
	c.JSON(http.StatusOK, ConfigInfo)
}

// fileExists 检查指定路径的文件是否存在
func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !os.IsNotExist(err)
}
