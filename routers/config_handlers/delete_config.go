package config_handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/logger"
	"net/http"
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
		logger.Infof("error: Failed save to %s directory", in)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed save to" + in + " directory"})
		return
	}
	err := config.DeleteConfigIn(in)
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
