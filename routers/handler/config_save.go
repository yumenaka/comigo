package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/yumenaka/comi/common"
	"io"
	"net/http"
)

// ConfigSaveHandler 保存服务器配置到文件
func ConfigSaveHandler(c *gin.Context) {
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
		common.Config.ConfigSaveTo = ConfigSaveTo.String()
		common.SaveConfig()
		// 返回成功消息
		c.JSON(http.StatusOK, gin.H{"message": "Config yaml save successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
	}
}
