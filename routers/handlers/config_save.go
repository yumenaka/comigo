package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/yumenaka/comi/config"
	"io"
	"net/http"
)

// HandlerConfigSave 保存服务器配置到文件
func HandlerConfigSave(c *gin.Context) {
	// 读取请求体中的JSON数据
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	// 将JSON数据转换为字符串并打印
	jsonString := string(body)
	fmt.Printf("Received JSON data: %s \n", jsonString)
	ConfigSaveTo := gjson.Get(jsonString, "ConfigSaveTo")
	if ConfigSaveTo.Exists() {
		config.Config.ConfigSaveTo = ConfigSaveTo.String()
		config.SaveConfig()
		// 返回成功消息
		c.JSON(http.StatusOK, gin.H{"message": "Config yaml save successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
	}
}
