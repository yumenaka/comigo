package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/settings"
	"github.com/yumenaka/comi/tools"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
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
	// 读取请求体中的JSON数据
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	// 将JSON数据转换为字符串并打印
	jsonString := string(body)
	fmt.Printf("Received JSON data: %s", jsonString)

	// 解析JSON数据并更新服务器配置
	newConfig, err := settings.UpdateConfig(common.Config, jsonString)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Failed to parse JSON data"})
		return
	}
	if (newConfig.OpenBrowser == true) && (common.Config.OpenBrowser == false) {
		protocol := "http://"
		if newConfig.EnableTLS {
			protocol = "https://"
		}
		tools.OpenBrowser(protocol + "127.0.0.1:" + strconv.Itoa(newConfig.Port))
	}
	needScan := false
	reScanFile := false
	if !reflect.DeepEqual(common.Config.StoresPath, newConfig.StoresPath) {
		needScan = true
		common.Config.StoresPath = newConfig.StoresPath
	}
	if common.Config.MaxScanDepth != newConfig.MaxScanDepth {
		needScan = true
		common.Config.MaxScanDepth = newConfig.MaxScanDepth
	}
	if !reflect.DeepEqual(common.Config.SupportMediaType, newConfig.SupportMediaType) {
		needScan = true
		reScanFile = true
		common.Config.SupportMediaType = newConfig.SupportMediaType
	}
	if !reflect.DeepEqual(common.Config.SupportFileType, newConfig.SupportFileType) {
		needScan = true
		common.Config.SupportFileType = newConfig.SupportFileType
	}
	if common.Config.MinImageNum != newConfig.MinImageNum {
		needScan = true
		reScanFile = true
		common.Config.MinImageNum = newConfig.MinImageNum
	}
	if !reflect.DeepEqual(common.Config.ExcludePath, newConfig.ExcludePath) {
		needScan = true
		common.Config.ExcludePath = newConfig.ExcludePath
	}
	if common.Config.EnableDatabase != newConfig.EnableDatabase {
		needScan = true
		common.Config.EnableDatabase = newConfig.EnableDatabase
	}
	if needScan {
		// 扫描配置文件指定的书籍库
		if err := common.ScanStorePath(reScanFile); err != nil {
			log.Printf("Failed to scan store path: %v", err)
		}
		// 保存扫描结果到数据库
		if common.Config.EnableDatabase {
			if err := common.SaveResultsToDatabase(); err != nil {
				log.Printf("Failed to save results to database: %v", err)
			}
		}
	} else {
		fmt.Println("common.Config.StoresPath == newConfig.StoresPath,skip scan store path")
	}
	common.Config = newConfig
	// 返回成功消息
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
	tempConfig.Username = "comigo"
	tempConfig.Password = ""
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
