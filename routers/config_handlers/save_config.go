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
	//如果不是三个目录之一，就不能保存
	if !(SaveTo == "WorkingDirectory" || SaveTo == "HomeDirectory" || SaveTo == "ProgramDirectory") {
		logger.Infof("error: Failed save to " + SaveTo + " directory")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed save to " + SaveTo + " directory"})
		return
	}
	//如果其他目录已经存在了，就不能保存（暂不支持多配置文件）
	if SaveTo == "WorkingDirectory" && (config.Status.Path.HomeDirectory != "" || config.Status.Path.ProgramDirectory != "") {
		logger.Infof("error: Find config in " + config.Status.Path.HomeDirectory + " " + config.Status.Path.ProgramDirectory)
		c.JSON(http.StatusBadRequest, gin.H{"error": "error: Find config in " + config.Status.Path.HomeDirectory + " " + config.Status.Path.ProgramDirectory})
		return
	}
	if SaveTo == "HomeDirectory" && (config.Status.Path.WorkingDirectory != "" || config.Status.Path.ProgramDirectory != "") {
		logger.Infof("error: Find config in " + config.Status.Path.WorkingDirectory + " " + config.Status.Path.ProgramDirectory)
		c.JSON(http.StatusBadRequest, gin.H{"error": "error: Find config in " + config.Status.Path.WorkingDirectory + " " + config.Status.Path.ProgramDirectory})
		return
	}
	if SaveTo == "ProgramDirectory" && (config.Status.Path.WorkingDirectory != "" || config.Status.Path.HomeDirectory != "") {
		logger.Infof("error: Find config in " + config.Status.Path.WorkingDirectory + " " + config.Status.Path.HomeDirectory)
		c.JSON(http.StatusBadRequest, gin.H{"error": "error: Find config in " + config.Status.Path.WorkingDirectory + " " + config.Status.Path.HomeDirectory})
		return
	}
	//保存配置
	err := config.SaveConfig(SaveTo)
	if err != nil {
		logger.Infof("%s", err.Error())
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Failed to save config"})
		return
	}
	// 返回成功消息
	GetConfigStatus(c)
}
