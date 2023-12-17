package handlers

import (
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/go-homedir"
	"github.com/pelletier/go-toml/v2"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/logger"
)

// SaveConfig 保存服务器配置到文件
func SaveConfig(c *gin.Context) {
	SaveTo := c.Param("to")
	if !(SaveTo == "WorkingDirectory" || SaveTo == "HomeDirectory" || SaveTo == "ProgramDirectoryectory") {
		logger.Info("error: Failed save to " + SaveTo + " directory")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed save to " + SaveTo + " directory"})
		return
	}
	err := saveConfigTo(SaveTo)
	if err != nil {
		logger.Info(err.Error())
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Failed to save config"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Config yaml save successfully"})
}

func saveConfigTo(Directory string) error {
	//保存配置
	bytes, err := toml.Marshal(config.Config)
	if err != nil {
		return err
	}
	logger.Info("Config Save To " + Directory)
	// Home 目录
	if Directory == "HomeDirectory" {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		err = os.MkdirAll(path.Join(home, ".config/comigo/"), os.ModePerm)
		if err != nil {
			return err
		}
		err = os.WriteFile(path.Join(home, ".config/comigo/config.toml"), bytes, 0644)
		if err != nil {
			return err
		}
	}
	//当前执行目录
	if Directory == "WorkingDirectory" {
		err = os.WriteFile("config.toml", bytes, 0644)
		if err != nil {
			return err
		}
	}
	// 可执行程序自身的文件路径
	if Directory == "ProgramDirectoryectory" {
		executablePath, err := os.Executable()
		if err != nil {
			return err
		}
		err = os.WriteFile(path.Join(executablePath, "config.toml"), bytes, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
