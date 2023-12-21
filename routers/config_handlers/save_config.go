package config_handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/logger"
	"net/http"
)

// SaveConfigHandler 保存服务器配置到文件
func SaveConfigHandler(c *gin.Context) {
	SaveTo := c.Param("to")
	if !(SaveTo == "WorkingDirectory" || SaveTo == "HomeDirectory" || SaveTo == "ProgramDirectory") {
		logger.Info("error: Failed save to " + SaveTo + " directory")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed save to " + SaveTo + " directory"})
		return
	}
	err := config.SaveConfig(SaveTo)
	if err != nil {
		logger.Info(err.Error())
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Failed to save config"})
		return
	}
	GetConfigStatus(c)
}
