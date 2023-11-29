package handlers

import (
	"errors"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/go-homedir"
	"github.com/yumenaka/comi/logger"
)

// DeleteConfig 删除服务器配置文件
func DeleteConfig(c *gin.Context) {
	in := c.Param("in")
	if !(in == "ExecutionDirectory" || in == "HomeDirectory" || in == "ProgramDirectory") {
		logger.Info("error: Failed save to " + in + " directory")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed save to " + in + " directory"})
		return
	}
	err := deleteConfigIn(in)
	if err != nil {
		logger.Info(err.Error())
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Failed to save config"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Config yaml save successfully"})
}

func deleteConfigIn(Directory string) error {
	logger.Info("Delete Config in " + Directory)
	// Home 目录
	if Directory == "HomeDirectory" {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		err = deleteFileIfExist(path.Join(home, ".config/comigo/config.toml"))
		if err != nil {
			return err
		}
	}
	//当前执行目录
	if Directory == "ExecutionDirectory" {
		err := deleteFileIfExist("config.toml")
		if err != nil {
			return err
		}
	}
	// 可执行程序自身的文件路径
	if Directory == "ProgramDirectory" {
		executablePath, err := os.Executable()
		if err != nil {
			return err
		}
		err = deleteFileIfExist(path.Join(executablePath, "config.toml"))
		if err != nil {
			return err
		}
	}
	return nil
}

// 删除文件
func deleteFileIfExist(filePath string) error {
	// 使用os.Stat检查文件是否存在
	if _, err := os.Stat(filePath); err == nil {
		// 文件存在，尝试删除
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	} else if os.IsNotExist(err) {
		return errors.New("File does not exist:" + filePath)
	} else {
		return err
	}
	return nil
}
