package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/config"
	"net/http"
)

// GetConfig 获取json格式的当前配置
func GetConfig(c *gin.Context) {
	//golang结构体默认深拷贝（但是基本类型浅拷贝）
	tempConfig := config.Config
	// Respond with the current server settings
	c.JSON(http.StatusOK, tempConfig)
}
