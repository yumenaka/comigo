package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"

	"github.com/yumenaka/comi/common"
)

// 下载服务器配置
func GetConfigHandler(c *gin.Context) {
	//golang结构体默认深拷贝（但是基本类型浅拷贝）
	tempConfig := common.Config
	tempConfig.GenerateConfig = false
	tempConfig.LogFilePath = ""
	common.Config.EnableDatabase = true
	tempConfig.StoresPath = []string{"C:\\test\\Comic", "D:\\some_path\\book", "/home/username/download"}
	tempConfig.CacheFilePath = ".comigo"
	bytes, err := toml.Marshal(tempConfig)
	if err != nil {
		fmt.Println("toml.Marshal Error")
	}
	//在命令行打印
	fmt.Println(string(bytes))
	//用gin实现下载文件的功能，只需要在接口返回时设置Response-Header中的Content-Type为文件类型，并设置Content-Disposition指定默认的文件名，然后将文件数据返回浏览器即可
	fileContentDisposition := "attachment;filename=\"" + "config.toml" + "\""
	c.Header("Content-Type", "application/octet-stream") // 这里是压缩文件类型 .zip
	c.Header("Content-Disposition", fileContentDisposition)
	c.Data(http.StatusOK, "application/octet-stream", bytes)
}
