package handlers

import (
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/go-homedir"
	"github.com/yumenaka/comi/logger"
	"github.com/yumenaka/comi/util"
)

const (
	HomeDirectory    = "HomeDirectory"
	WorkingDirectory = "WorkingDirectory"
	ProgramDirectory = "ProgramDirectory"
)

func DeleteConfig(c *gin.Context) {
	in := c.Param("in")
	validDirs := []string{WorkingDirectory, HomeDirectory, ProgramDirectory}

	if !contains(validDirs, in) {
		logger.Info("error: Failed save to" + in + " directory")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed save to" + in + " directory"})
		return
	}
	err := deleteConfigIn(in)
	if err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Failed to save config"})
		return
	}
	GetConfigStatus(c)
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
	case ProgramDirectory:
		executablePath, err := os.Executable()
		if err != nil {
			return err
		}
		filePath = path.Join(path.Dir(executablePath), "config.toml")
	}
	if err != nil {
		return err
	}
	return util.DeleteFileIfExist(filePath)
}
