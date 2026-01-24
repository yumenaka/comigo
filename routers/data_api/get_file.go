package data_api

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	fileutil "github.com/yumenaka/comigo/tools/file"
	"github.com/yumenaka/comigo/tools/logger"
)

// GetFile 示例 URL： 127.0.0.1:1234/get_file?id=2b17a13&filename=1.jpg
// 缩放文件，会转化为jpeg：http://127.0.0.1:1234/api/get-file?id=2b17a13&resize_width=300&resize_height=400&id=597e06&filename=01.jpeg
// 相关参数：
// id：书籍的ID，必须参数       							&id=2B17a
// filename:获取的文件名，必须参数   							&filename=01.jpg
// 可选参数：
// resize_width:指定宽度，缩放图片  							&resize_width=300
// resize_height:指定高度，缩放图片 							&resize_height=300
// resize_max_width:指定宽度上限，图片宽度大于这个上限时缩小图片  	&resize_max_width=740
// resize_max_height:指定高度上限，图片高度大于这个上限时缩小图片  	&resize_max_height=300
// auto_crop:自动切白边，数字是切白边的阈值，范围是0~100 			&auto_crop=10
// gray:黑白化												&gray=true
// blurhash:获取对应图片的blurhash，不是原始图片 				&blurhash=3
// blurhash_image:获取对应图片的blurhash图片，不是原始图片  	    &blurhash_image=3
// base64:返回Base64编码的图片								    &base64=true
func GetFile(c echo.Context) error {
	id := c.QueryParam("id")
	needFile := c.QueryParam("filename")
	// 必须指定 id 和 filename
	if id == "" || needFile == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id and filename are required"})
	}

	// 读取查询参数
	disableCache := getBoolQueryParam(c, "no-cache", false)

	// 读取图片处理参数
	resizeWidth := getIntQueryParam(c, "resize_width", 0)
	resizeHeight := getIntQueryParam(c, "resize_height", 0)
	autoCrop := getIntQueryParam(c, "auto_crop", 0)
	resizeMaxWidth := getIntQueryParam(c, "resize_max_width", 0)
	resizeMaxHeight := getIntQueryParam(c, "resize_max_height", 0)
	blurhash := getIntQueryParam(c, "blurhash", 0)
	blurhashImage := getIntQueryParam(c, "blurhash_image", 0)
	gray := getBoolQueryParam(c, "gray", false)

	// 当有任何图片处理参数生效时，禁用缓存，避免后续返回错误的缓存结果
	if resizeWidth > 0 || resizeHeight > 0 || autoCrop > 0 ||
		resizeMaxWidth > 0 || resizeMaxHeight > 0 ||
		blurhash > 0 || blurhashImage > 0 || gray {
		disableCache = true
	}

	// 如果启用了本地缓存
	if config.GetCfg().UseCache && !disableCache {
		// 获取所有的参数键值对
		query := c.Request().URL.Query()
		// 如果有缓存，读取本地获取缓存文件并返回
		cacheData, ct, err := fileutil.GetFileFromCache(
			id,
			needFile,
			fileutil.GetQueryString(query),
			config.GetCfg().CacheDir,
			config.GetCfg().Debug,
		)
		if err == nil && cacheData != nil {
			return c.Blob(http.StatusOK, ct, cacheData)
		}
	}
	// 获取书籍信息
	book, err := model.IStore.GetBook(id)
	if err != nil {
		logger.Infof("GetBook error: %s", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Book not found"})
	}

	// 获取图片数据的选项
	option := fileutil.GetPictureDataOption{
		PictureName:      needFile,
		BookIsPDF:        book.Type == model.TypePDF,
		BookIsDir:        book.Type == model.TypeDir,
		BookIsNonUTF8Zip: book.NonUTF8Zip,
		BookPath:         book.BookPath,
		Debug:            config.GetCfg().Debug,
		UseCache:         config.GetCfg().UseCache,
		ResizeWidth:      resizeWidth,
		ResizeHeight:     resizeHeight,
		ResizeMaxWidth:   resizeMaxWidth,
		ResizeMaxHeight:  resizeMaxHeight,
		AutoCrop:         autoCrop,
		Gray:             gray,
		BlurHash:         blurhash,
		BlurHashImage:    blurhashImage,
	}

	// 获取图片数据
	imgData, contentType, err := fileutil.GetPictureData(option)
	if err != nil {
		logger.Infof("Get file error: %s", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Get file error: " + err.Error()})
	}

	// 缓存文件到本地，避免重复解压。如果书中的图片，来自本地目录，就不需要缓存。
	if config.GetCfg().UseCache && !disableCache && book.Type != model.TypeDir {
		// 获取所有的参数键值对
		query := c.Request().URL.Query()
		errSave := fileutil.SaveFileToCache(
			id,
			needFile,
			imgData,
			fileutil.GetQueryString(query),
			contentType,
			config.GetCfg().CacheDir,
			config.GetCfg().Debug,
		)
		if errSave != nil {
			logger.Infof("SaveFileToCache error: %s", errSave)
		}
	}
	// 返回图片数据
	return c.Blob(http.StatusOK, contentType, imgData)
}

// getIntQueryParam 从查询参数中获取整数值，带默认值
func getIntQueryParam(c echo.Context, key string, defaultValue int) int {
	valueStr := c.QueryParam(key)
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
func getBoolQueryParam(c echo.Context, key string, defaultValue bool) bool {
	valueStr := c.QueryParam(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
