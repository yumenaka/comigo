package data_api

import (
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools"
	fileutil "github.com/yumenaka/comigo/tools/file"
	"github.com/yumenaka/comigo/tools/logger"
)

// GetCover 获取书籍封面
// 相关参数：
// id：书籍的ID，必须参数 &id=2B17a
// resize_height：可选参数，指定封面高度，默认值为352 &resize_height=500
// 示例 URL： http://127.0.0.1:1234/api/get_cover?id=2b17a13
// 示例 URL（自定义高度）： http://127.0.0.1:1234/api/get_cover?id=2b17a13&resize_height=500
func GetCover(c echo.Context) error {
	// 获取书籍ID
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}
	// 获取 resize_height 参数，默认值为 352
	resizeHeight := getIntQueryParam(c, "resize_height", 352)
	// 获取书籍信息
	book, err := model.IStore.GetBook(id)
	if err != nil {
		logger.Infof("GetBook error: %s", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Book not found"})
	}
	// 获取封面信息
	cover := book.GetCover()
	configDir, err := config.GetConfigDir()
	// 封面元数据目录路径
	metaPath := ""
	// 如果获取配置目录成功，尝试从本地缓存读取封面
	if err == nil {
		//  封面元数据目录路径
		metaPath = filepath.Join(configDir, "metadata", book.GetStoreID())
		// 先从本地缓存读取封面
		coverFileCacheExists := fileutil.CoverFileCacheExists(metaPath, id)
		if coverFileCacheExists {
			coverData, err := fileutil.GetCoverFromLocal(metaPath, id)
			if err == nil {
				return c.Blob(http.StatusOK, "image/jpeg", coverData)
			}
		}
	}
	// 如果本地没有缓存封面，就从压缩文件中获取封面
	if err != nil {
		logger.Errorf(locale.GetString("err_failed_to_get_config_dir"), err)
	}
	needFile := cover.Name
	// 处理 TypeBooksGroup 的情况：封面来自子书籍
	coverBook := book
	if book.Type == model.TypeBooksGroup && strings.HasPrefix(cover.Url, "/api/get_file") {
		// 解析封面 URL 获取子书籍 ID 和文件名
		// URL 格式：/api/get_file?id=子书籍ID&filename=文件名
		parsedURL, err := url.Parse(cover.Url)
		if err != nil {
			logger.Infof("Failed to parse cover URL: %s", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid cover URL"})
		}
		childID := parsedURL.Query().Get("id")
		if childID == "" {
			logger.Infof("Child book ID is missing in cover URL")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Child book ID is required in cover URL"})
		}
		// 获取子书籍信息
		childBook, err := model.IStore.GetBook(childID)
		if err != nil {
			logger.Infof("Failed to get child book: %s", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Child book not found"})
		}
		coverBook = childBook
		// 从 URL 中获取文件名，如果没有则使用 cover.Name
		if filename := parsedURL.Query().Get("filename"); filename != "" {
			needFile = filename
		}
	}
	// 如果封面URL是内嵌图片, 通过封面文件获取
	if strings.HasPrefix(cover.Url, "/images/") {
		// 从内嵌文件系统读取图片数据
		imgData := assets.GetImageData(cover.Name)
		if len(imgData) == 0 {
			logger.Infof("Failed to read embedded image: %s", cover.Url)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Embedded image not found"})
		}
		// 确定MIME类型
		contentType := tools.GetContentTypeByFileName(cover.Name)
		// 缩放图片到指定高度
		imgData = tools.ImageResizeByHeight(imgData, resizeHeight)
		// 缩放后的图片转换为JPEG格式
		contentType = tools.GetContentTypeByFileName(".jpg")
		// 缓存封面到 configDir
		if configDir != "" {
			// 内嵌图片也缓存封面到本地（方便以后做手动指定封面功能）
			err = fileutil.SaveCoverToLocal(metaPath, id, imgData)
			if err != nil {
				logger.Infof("SaveCoverToLocal error: %s", err)
			}
		}
		// 返回图片数据
		return c.Blob(http.StatusOK, contentType, imgData)
	}
	// 获取图片数据的选项
	option := fileutil.GetPictureDataOption{
		PictureName:      needFile,
		BookPath:         coverBook.BookPath,
		BookIsPDF:        coverBook.Type == model.TypePDF,
		BookIsDir:        coverBook.Type == model.TypeDir,
		BookIsNonUTF8Zip: coverBook.NonUTF8Zip,
		Debug:            config.GetCfg().Debug,
		ResizeHeight:     resizeHeight,
	}
	// 获取图片数据
	imgData, contentType, err := fileutil.GetPictureData(option)
	if err != nil {
		logger.Infof("Get file error: %s", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Get file error: " + err.Error()})
	}
	// 缓存封面到 configDir
	if configDir != "" {
		// 缓存封面到本地
		err = fileutil.SaveCoverToLocal(metaPath, id, imgData)
		if err != nil {
			logger.Infof("SaveCoverToLocal error: %s", err)
		}
	}
	// 返回图片数据
	return c.Blob(http.StatusOK, contentType, imgData)
}
