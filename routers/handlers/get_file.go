package handlers

import (
	"encoding/base64"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	fileutil "github.com/yumenaka/comigo/util/file"
	"github.com/yumenaka/comigo/util/logger"
)

// GetFile 示例 URL： 127.0.0.1:1234/get_file?id=2b17a13&filename=1.jpg
// 缩放文件，会转化为jpeg：http://127.0.0.1:1234/api/get_file?resize_width=300&resize_height=400&id=597e06&filename=01.jpeg
// 相关参数：
// id：书籍的ID，必须参数       							&id=2B17a
// filename:获取的文件名，必须参数   							&filename=01.jpg
// 可选参数：
// resize_width:指定宽度，缩放图片  							&resize_width=300
// resize_height:指定高度，缩放图片 							&resize_height=300
// thumbnail_mode:缩略图模式，同时指定宽高的时候要不要剪切图片		&thumbnail_mode=true
// resize_max_width:指定宽度上限，图片宽度大于这个上限时缩小图片  	&resize_max_width=740
// resize_max_height:指定高度上限，图片高度大于这个上限时缩小图片  	&resize_max_height=300
// auto_crop:自动切白边，数字是切白边的阈值，范围是0~100 			&auto_crop=10
// gray:黑白化												&gray=true
// blurhash:获取对应图片的blurhash，不是原始图片 				&blurhash=3
// blurhash_image:获取对应图片的blurhash图片，不是原始图片  	    &blurhash_image=3
// base64:返回Base64编码的图片								    &base64=true
func GetFile(c *gin.Context) {
	id := c.Query("id")
	needFile := c.Query("filename")
	// 必须指定 id 和 filename
	if id == "" || needFile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id and filename are required"})
		return
	}
	// 读取查询参数
	noCache := getBoolQueryParam(c, "no-cache", false)
	base64Encode := getBoolQueryParam(c, "base64", false)

	// 读取图片处理参数
	resizeWidth := getIntQueryParam(c, "resize_width", 0)
	resizeHeight := getIntQueryParam(c, "resize_height", 0)
	autoCrop := getIntQueryParam(c, "auto_crop", 0)
	resizeMaxWidth := getIntQueryParam(c, "resize_max_width", 0)
	resizeMaxHeight := getIntQueryParam(c, "resize_max_height", 0)
	blurhash := getIntQueryParam(c, "blurhash", 0)
	blurhashImage := getIntQueryParam(c, "blurhash_image", 0)
	thumbnailMode := getBoolQueryParam(c, "thumbnail_mode", false)
	gray := getBoolQueryParam(c, "gray", false)

	// 如果启用了本地缓存
	if config.GetUseCache() && !noCache {
		// 获取所有的参数键值对
		query := c.Request.URL.Query()
		// 如果有缓存，直接读取本地获取缓存文件并返回
		cacheData, ct, err := fileutil.GetFileFromCache(
			id,
			needFile,
			fileutil.GetQueryString(query),
			thumbnailMode,
			config.GetConfigPath(),
			config.GetDebug(),
		)
		if err == nil && cacheData != nil {
			if base64Encode {
				sendBase64Data(c, cacheData, needFile)
				return
			}
			c.Data(http.StatusOK, ct, cacheData)
			return
		}
	}
	// 获取书籍信息
	bookByID, err := model.GetBookByID(id, "")
	if err != nil {
		logger.Infof("GetBookByID error: %s", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	// 获取图片数据的选项
	option := fileutil.GetPictureDataOption{
		PictureName:      needFile,
		BookIsPDF:        bookByID.Type == model.TypePDF,
		BookIsDir:        bookByID.Type == model.TypeDir,
		BookIsNonUTF8Zip: bookByID.NonUTF8Zip,
		BookFilePath:     bookByID.FilePath,
		Debug:            config.GetDebug(),
		UseCache:         config.GetUseCache(),
		ResizeWidth:      resizeWidth,
		ResizeHeight:     resizeHeight,
		ResizeMaxWidth:   resizeMaxWidth,
		ResizeMaxHeight:  resizeMaxHeight,
		ThumbnailMode:    thumbnailMode,
		AutoCrop:         autoCrop,
		Gray:             gray,
		BlurHash:         blurhash,
		BlurHashImage:    blurhashImage,
	}
	// 获取图片数据
	imgData, contentType, err := fileutil.GetPictureData(option)
	if err != nil {
		logger.Infof("GetPictureData error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "GetPictureData error: " + err.Error()})
		return
	}
	// 缓存文件到本地，避免重复解压。如果书中的图片，来自本地目录，就不需要缓存。
	if config.GetUseCache() && !noCache && bookByID.Type != model.TypeDir {
		// 获取所有的参数键值对
		query := c.Request.URL.Query()
		errSave := fileutil.SaveFileToCache(
			id,
			needFile,
			imgData,
			fileutil.GetQueryString(query),
			contentType,
			thumbnailMode,
			config.GetConfigPath(),
			config.GetDebug(),
		)
		if errSave != nil {
			logger.Infof("SaveFileToCache error: %s", errSave)
		}
	}
	// 如果启用了 Base64 编码
	if base64Encode {
		sendBase64Data(c, imgData, needFile)
		return
	}
	// 返回图片数据
	c.Data(http.StatusOK, contentType, imgData)
}

// getIntQueryParam 从查询参数中获取整数值，带默认值
func getIntQueryParam(c *gin.Context, key string, defaultValue int) int {
	valueStr := c.Query(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getBoolQueryParam 从查询参数中获取布尔值，带默认值
func getBoolQueryParam(c *gin.Context, key string, defaultValue bool) bool {
	valueStr := c.Query(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// sendBase64Data 将数据编码为 Base64 并发送
func sendBase64Data(c *gin.Context, data []byte, filename string) {
	mimeType := mime.TypeByExtension(filepath.Ext(filename))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	base64Data := base64.StdEncoding.EncodeToString(data)
	dataURI := "data:" + mimeType + ";base64," + base64Data
	c.String(http.StatusOK, dataURI)
}
