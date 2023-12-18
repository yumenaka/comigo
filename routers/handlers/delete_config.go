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

const (
	HomeDirectory          = "HomeDirectory"
	WorkingDirectory       = "WorkingDirectory"
	ProgramDirectoryectory = "ProgramDirectory"
)

func DeleteConfig(c *gin.Context) {
	in := c.Param("in")
	validDirs := []string{WorkingDirectory, HomeDirectory, ProgramDirectoryectory}

	if !contains(validDirs, in) {
		logger.Info("error: Failed save to" + in + " directory")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed save to" + in + " directory"})
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

// contains 函数来检查切片是否包含特定字符串
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func deleteConfigIn(Directory string) (err error) {
	logger.Info("Delete Config in " + Directory)
	var filePath string

	switch Directory {
	case HomeDirectory:
		home, err := homedir.Dir()
		if err == nil {
			filePath = path.Join(home, ".config/comigo/config.toml")
		}
	case WorkingDirectory:
		filePath = "config.toml"
	case ProgramDirectoryectory:
		filePath, err = getProgramDirectoryectoryConfigFilePath()
	}
	if err != nil {
		return err
	}
	return deleteFileIfExist(filePath)
}

func getProgramDirectoryectoryConfigFilePath() (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return path.Join(path.Dir(executablePath), "config.toml"), nil
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
