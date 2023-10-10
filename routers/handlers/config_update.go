package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/arch/scan"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/types"
	"github.com/yumenaka/comi/util"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
)

// HandlerConfigUpdate 修改服务器配置(post json)
func HandlerConfigUpdate(c *gin.Context) {
	// 读取请求体中的JSON数据
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	// 将JSON数据转换为字符串并打印
	jsonString := string(body)
	fmt.Printf("Received JSON data: %s \n", jsonString)

	// 解析JSON数据并更新服务器配置
	newConfig, err := types.UpdateConfig(config.Config, jsonString)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Failed to parse JSON data"})
		return
	}
	BeforeConfigUpdate(&config.Config, &newConfig)

	// 返回成功消息
	c.JSON(http.StatusOK, gin.H{"message": "Server settings updated successfully"})
}

// BeforeConfigUpdate 根据配置的变化，判断是否需要打开浏览器重新扫描等
func BeforeConfigUpdate(oldConfig *types.ServerConfig, newConfig *types.ServerConfig) {
	if (newConfig.OpenBrowser == true) && (oldConfig.OpenBrowser == false) {
		protocol := "http://"
		if newConfig.EnableTLS {
			protocol = "https://"
		}
		util.OpenBrowser(protocol + "127.0.0.1:" + strconv.Itoa(newConfig.Port))
	}
	needScan := false
	reScanFile := false
	if !reflect.DeepEqual(oldConfig.StoresPath, newConfig.StoresPath) {
		needScan = true
		oldConfig.StoresPath = newConfig.StoresPath
	}
	if oldConfig.MaxScanDepth != newConfig.MaxScanDepth {
		needScan = true
		oldConfig.MaxScanDepth = newConfig.MaxScanDepth
	}
	if !reflect.DeepEqual(oldConfig.SupportMediaType, newConfig.SupportMediaType) {
		needScan = true
		reScanFile = true
		oldConfig.SupportMediaType = newConfig.SupportMediaType
	}
	if !reflect.DeepEqual(oldConfig.SupportFileType, newConfig.SupportFileType) {
		needScan = true
		oldConfig.SupportFileType = newConfig.SupportFileType
	}
	if oldConfig.MinImageNum != newConfig.MinImageNum {
		needScan = true
		reScanFile = true
		oldConfig.MinImageNum = newConfig.MinImageNum
	}
	if !reflect.DeepEqual(oldConfig.ExcludePath, newConfig.ExcludePath) {
		needScan = true
		oldConfig.ExcludePath = newConfig.ExcludePath
	}
	if oldConfig.EnableDatabase != newConfig.EnableDatabase {
		needScan = true
		oldConfig.EnableDatabase = newConfig.EnableDatabase
	}
	if needScan {
		option := scan.NewScanOption(
			reScanFile,
			newConfig.StoresPath,
			newConfig.MaxScanDepth,
			newConfig.MinImageNum,
			newConfig.TimeoutLimitForScan,
			newConfig.ExcludePath,
			newConfig.SupportMediaType,
			newConfig.SupportFileType,
			newConfig.ZipFileTextEncoding,
			newConfig.EnableDatabase,
			newConfig.ClearDatabaseWhenExit,
			newConfig.Debug,
		)
		if err := scan.ScanStorePath(option); err != nil {
			log.Printf("Failed to scan store path: %v", err)
		}
		// 保存扫描结果到数据库 //TODO:这里有问题，启用数据库会报错
		if oldConfig.EnableDatabase {
			if err := scan.SaveResultsToDatabase(config.Config.UploadPath, config.Config.ClearDatabaseWhenExit, config.Config.Debug); err != nil {
				log.Printf("Failed to save results to database: %v", err)
			}
		}
	} else {
		if oldConfig.Debug {
			log.Printf("oldConfig.StoresPath == newConfig.StoresPath,skip scan store path")
		}
	}
	oldConfig = newConfig
}
