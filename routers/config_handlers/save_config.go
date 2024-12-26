package config_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/util/logger"
)

// SaveConfigHandler 保存服务器配置到文件
func SaveConfigHandler(c *gin.Context) {
	SaveTo := c.Param("to")
	//如果不是三个目录之一，就不能保存
	if !(SaveTo == "WorkingDirectory" || SaveTo == "HomeDirectory" || SaveTo == "ProgramDirectory") {
		logger.Infof("error: Failed save to %s directory", SaveTo)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed save to " + SaveTo + " directory"})
		return
	}
	//如果其他目录有配置文件，就不能保存（暂不支持多配置文件）
	if SaveTo == "WorkingDirectory" && (config.CfgStatus.Path.HomeDirectory != "" || config.CfgStatus.Path.ProgramDirectory != "") {
		logger.Infof("error: Find config in %s %s", config.CfgStatus.Path.HomeDirectory, config.CfgStatus.Path.ProgramDirectory)
		c.JSON(http.StatusBadRequest, gin.H{"error": "error: Find config in " + config.CfgStatus.Path.HomeDirectory + " " + config.CfgStatus.Path.ProgramDirectory})
		return
	}
	if SaveTo == "HomeDirectory" && (config.CfgStatus.Path.WorkingDirectory != "" || config.CfgStatus.Path.ProgramDirectory != "") {
		logger.Infof("error: Find config in  %s %s", config.CfgStatus.Path.WorkingDirectory, config.CfgStatus.Path.ProgramDirectory)
		c.JSON(http.StatusBadRequest, gin.H{"error": "error: Find config in " + config.CfgStatus.Path.WorkingDirectory + " " + config.CfgStatus.Path.ProgramDirectory})
		return
	}
	if SaveTo == "ProgramDirectory" && (config.CfgStatus.Path.WorkingDirectory != "" || config.CfgStatus.Path.HomeDirectory != "") {
		logger.Infof("error: Find config in  %s %s", config.CfgStatus.Path.WorkingDirectory, config.CfgStatus.Path.HomeDirectory)
		c.JSON(http.StatusBadRequest, gin.H{"error": "error: Find config in " + config.CfgStatus.Path.WorkingDirectory + " " + config.CfgStatus.Path.HomeDirectory})
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
