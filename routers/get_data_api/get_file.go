package get_data_api

import (
	"encoding/base64"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	fileutil "github.com/yumenaka/comigo/tools/file"
	"github.com/yumenaka/comigo/tools/logger"
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
func GetFile(c echo.Context) error {
	id := c.QueryParam("id")
	needFile := c.QueryParam("filename")
	// 必须指定 id 和 filename
	if id == "" || needFile == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id and filename are required"})
	}

	// 读取查询参数
	noCache := getBoolQueryParam(c, "no-cache", false)
	base64Htmx := getBoolQueryParam(c, "base64_htmx", false)
	base64String := getBoolQueryParam(c, "base64_string", false)

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
		query := c.Request().URL.Query()
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
			// 如果启用了 Base64 编码
			if base64Htmx {
				return sendBase64Htmx(c, cacheData, needFile)
			}
			if base64String {
				return sendBase64String(c, cacheData, needFile)
			}
			return c.Blob(http.StatusOK, ct, cacheData)
		}
	}

	// 获取书籍信息
	bookByID, err := model.MainStoreGroup.GetBookByID(id, "")
	if err != nil {
		logger.Infof("GetBookByID error: %s", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Book not found"})
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "GetPictureData error: " + err.Error()})
	}

	// 缓存文件到本地，避免重复解压。如果书中的图片，来自本地目录，就不需要缓存。
	if config.GetUseCache() && !noCache && bookByID.Type != model.TypeDir {
		// 获取所有的参数键值对
		query := c.Request().URL.Query()
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
	if base64Htmx {
		return sendBase64Htmx(c, imgData, needFile)
	}
	if base64String {
		return sendBase64String(c, imgData, needFile)
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

// sendBase64Htmx 将数据编码为 Base64 并通过 HTMX 接口替换 img 标签
func sendBase64String(c echo.Context, data []byte, filename string) error {
	mimeType := mime.TypeByExtension(filepath.Ext(filename))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	dataURI := "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(data)
	// 纯文本返回，Content-Type 设 text/plain 保证浏览器不当成 HTML
	return c.String(http.StatusOK, dataURI)
}

// sendBase64Htmx 将数据编码为 Base64 并通过 HTMX 接口替换 img 标签
func sendBase64Htmx(c echo.Context, data []byte, filename string) error {
	mimeType := mime.TypeByExtension(filepath.Ext(filename))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	dataURI := "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(data)

	htmlContent := `<img
		x-data="{ isDoublePage: false }"
		class="w-full manga_image min-h-16 text-center"
		draggable="false"
		src="` + dataURI + `"
		@load="isDoublePage=$event.target.naturalWidth > $event.target.naturalHeight;"
		:style="{ width: $store.global.isLandscape?($store.scroll.widthUseFixedValue? (isDoublePage ? $store.scroll.doublePageWidth_PX +'px': $store.scroll.singlePageWidth_PX +'px'): (isDoublePage ? $store.scroll.doublePageWidth_Percent + '%' : $store.scroll.singlePageWidth_Percent + '%')): $store.scroll.portraitWidthPercent+'%', maxWidth: '100%' }"
		alt=""
	/>`
	// 直接返回 HTML 片段，让 HTMX 用 outerHTML 替换
	return c.HTML(http.StatusOK, htmlContent)
}
