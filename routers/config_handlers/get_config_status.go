package config_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/config"
)

// GetConfigStatus 获取json格式的当前配置
func GetConfigStatus(c *gin.Context) {
	err := config.Status.SetConfigStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get config"})
		return
	}
	c.IndentedJSON(http.StatusOK, config.Status)
}
