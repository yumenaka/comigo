package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/common"
	"net/http"
)

// GetConfigJsonHandler 获取json格式的当前配置
func GetConfigJsonHandler(c *gin.Context) {
	//golang结构体默认深拷贝（但是基本类型浅拷贝）
	tempConfig := common.Config
	// Respond with the current server settings
	c.JSON(http.StatusOK, tempConfig)
}
