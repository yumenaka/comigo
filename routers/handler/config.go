package handler

import (
	"fmt"
	"github.com/yumenaka/comi/settings"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"

	"github.com/yumenaka/comi/common"
)

// GetJsonConfigHandler 获取json格式的当前配置
func GetJsonConfigHandler(c *gin.Context) {
	//golang结构体默认深拷贝（但是基本类型浅拷贝）
	tempConfig := common.Config
	// Respond with the current server settings
	c.JSON(http.StatusOK, tempConfig)
}

// UpdateConfigHandler 修改服务器配置(post json)
func UpdateConfigHandler(c *gin.Context) {
	var newSettings settings.ServerSettings
	// Bind the JSON data from the request body into the newSettings variable
	if err := c.BindJSON(&newSettings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON data"})
		return
	}
	// Update the server settings with the new values
	common.Config = newSettings
	//扫描配置文件指定的书籍库
	common.ScanStorePathInConfig()
	//TODO:保存扫描结果到数据库
	common.SaveResultsToDatabase()
	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Server settings updated successfully"})
}

// GetTomlConfigHandler 下载服务器配置(toml)
func GetTomlConfigHandler(c *gin.Context) {
	//golang结构体默认深拷贝（但是基本类型浅拷贝）
	tempConfig := common.Config
	tempConfig.LogFilePath = ""
	common.Config.OpenBrowser = false
	common.Config.EnableDatabase = true
	tempConfig.StoresPath = []string{"C:\\test\\Comic", "D:\\some_path\\book", "/home/user/download"}
	tempConfig.UserName = "admin"
	tempConfig.Password = "admin"
	bytes, err := toml.Marshal(tempConfig)
	if err != nil {
		fmt.Println("toml.Marshal Error")
	}
	//在命令行打印
	fmt.Println(string(bytes))
	//用gin实现下载文件的功能，只需要在接口返回时设置Response-Header中的Content-Type为文件类型，并设置Content-Disposition指定默认的文件名，然后将文件数据返回浏览器即可
	fileContentDisposition := "attachment;filename=\"" + "config.toml" + "\""
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fileContentDisposition)
	c.Data(http.StatusOK, "application/octet-stream", bytes)
}
