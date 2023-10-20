package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/logger"
	"net/http"
)

// HandlerGetConfigToml 下载服务器配置(toml)
func HandlerGetConfigToml(c *gin.Context) {
	//golang结构体默认深拷贝（但是基本类型浅拷贝）
	tempConfig := config.Config
	tempConfig.LogFilePath = ""
	config.Config.OpenBrowser = false
	config.Config.EnableDatabase = true
	tempConfig.StoresPath = []string{"C:\\test\\Comic", "D:\\some_path\\book", "/home/user/download"}
	tempConfig.Username = "comigo"
	tempConfig.Password = ""
	bytes, err := toml.Marshal(tempConfig)
	if err != nil {
		logger.Info("toml.Marshal Error")
	}
	//在命令行打印
	logger.Info(string(bytes))
	//用gin实现下载文件的功能，只需要在接口返回时设置Response-Header中的Content-Type为文件类型，并设置Content-Disposition指定默认的文件名，然后将文件数据返回浏览器即可
	fileContentDisposition := "attachment;filename=\"" + "config.toml" + "\""
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fileContentDisposition)
	c.Data(http.StatusOK, "application/octet-stream", bytes)
}
