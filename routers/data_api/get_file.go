package data_api

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers/apiresp"
	fileutil "github.com/yumenaka/comigo/tools/file"
	"github.com/yumenaka/comigo/tools/logger"
)

type getFileRequest struct {
	bookID          string
	fileName        string
	disableCache    bool
	resizeWidth     int
	resizeHeight    int
	autoCrop        int
	resizeMaxWidth  int
	resizeMaxHeight int
	blurhash        int
	blurhashImage   int
	gray            bool
}

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
	req, err := parseGetFileRequest(c)
	if err != nil {
		return writeValidationError(c, err)
	}
	if localBook, client, _, ok, err := remoteComigoBookFromRequest(c, req.bookID); ok {
		if err != nil {
			return writeRemoteComigoError(c, err)
		}
		imgData, contentType, err := client.GetBytes("/api/get-file", remoteComigoQuery(c, localBook.RemoteBookID))
		if err != nil {
			return writeRemoteComigoError(c, err)
		}
		return c.Blob(http.StatusOK, contentType, imgData)
	}

	// 获取书籍信息
	book, err := model.IStore.GetBook(req.bookID)
	if err != nil {
		logger.Infof(locale.GetString("log_getbook_error_common"), err)
		return apiresp.Error(c, http.StatusNotFound, "book_not_found", "Book not found", map[string]string{"id": req.bookID})
	}
	if !bookContainsPage(book, req.fileName) {
		return apiresp.Error(c, http.StatusNotFound, "page_not_found", "Page not found", map[string]string{"filename": req.fileName})
	}
	// 缓存也必须经过书籍与页面校验，不能让旧缓存绕过当前元数据边界。
	if handled, err := serveCachedPicture(c, req); handled || err != nil {
		return err
	}

	// 获取图片数据
	imgData, contentType, err := fileutil.GetPictureData(buildGetPictureDataOption(req, book))
	if err != nil {
		logger.Infof(locale.GetString("log_get_file_error"), err)
		return apiresp.BadRequest(c, "get_file_failed", "Get file error: "+err.Error(), nil)
	}

	// 缓存文件到本地，避免重复解压。如果书中的图片，来自本地目录，就不需要缓存。
	savePictureCache(c, req, book, imgData, contentType)
	// 返回图片数据
	return c.Blob(http.StatusOK, contentType, imgData)
}

// bookContainsPage 只允许读取扫描后公开在书籍元数据中的页面。
func bookContainsPage(book *model.Book, filename string) bool {
	for _, page := range book.PageInfos {
		if page.Name == filename {
			return true
		}
	}
	return false
}

func parseGetFileRequest(c echo.Context) (getFileRequest, error) {
	resizeWidth, err := parseOptionalBoundedInt(c, "resize_width", 0, 1, imageQueryMaxDimension)
	if err != nil {
		return getFileRequest{}, err
	}
	resizeHeight, err := parseOptionalBoundedInt(c, "resize_height", 0, 1, imageQueryMaxDimension)
	if err != nil {
		return getFileRequest{}, err
	}
	autoCrop, err := parseOptionalBoundedInt(c, "auto_crop", 0, 0, imageQueryMaxAutoCrop)
	if err != nil {
		return getFileRequest{}, err
	}
	resizeMaxWidth, err := parseOptionalBoundedInt(c, "resize_max_width", 0, 1, imageQueryMaxDimension)
	if err != nil {
		return getFileRequest{}, err
	}
	resizeMaxHeight, err := parseOptionalBoundedInt(c, "resize_max_height", 0, 1, imageQueryMaxDimension)
	if err != nil {
		return getFileRequest{}, err
	}
	blurhash, err := parseOptionalBoundedInt(c, "blurhash", 0, 0, imageQueryMaxBlurComponents)
	if err != nil {
		return getFileRequest{}, err
	}
	blurhashImage, err := parseOptionalBoundedInt(c, "blurhash_image", 0, 0, imageQueryMaxBlurComponents)
	if err != nil {
		return getFileRequest{}, err
	}
	req := getFileRequest{
		bookID:          c.QueryParam("id"),
		fileName:        c.QueryParam("filename"),
		disableCache:    getBoolQueryParam(c, "no-cache", false),
		resizeWidth:     resizeWidth,
		resizeHeight:    resizeHeight,
		autoCrop:        autoCrop,
		resizeMaxWidth:  resizeMaxWidth,
		resizeMaxHeight: resizeMaxHeight,
		blurhash:        blurhash,
		blurhashImage:   blurhashImage,
		gray:            getBoolQueryParam(c, "gray", false),
	}
	if req.bookID == "" || req.fileName == "" {
		return req, requestValidationError{
			code:    "missing_param",
			message: "id and filename are required",
			details: []string{"id", "filename"},
		}
	}
	if hasImageTransform(req) {
		req.disableCache = true
	}
	return req, nil
}

func hasImageTransform(req getFileRequest) bool {
	return req.resizeWidth > 0 || req.resizeHeight > 0 || req.autoCrop > 0 ||
		req.resizeMaxWidth > 0 || req.resizeMaxHeight > 0 ||
		req.blurhash > 0 || req.blurhashImage > 0 || req.gray
}

func serveCachedPicture(c echo.Context, req getFileRequest) (bool, error) {
	if !config.GetCfg().UseCache || req.disableCache {
		return false, nil
	}
	cacheData, contentType, err := fileutil.GetFileFromCache(
		req.bookID,
		req.fileName,
		fileutil.GetQueryString(c.Request().URL.Query()),
		config.GetCfg().CacheDir,
		config.GetCfg().Debug,
	)
	if err != nil || cacheData == nil {
		return false, nil
	}
	return true, c.Blob(http.StatusOK, contentType, cacheData)
}

func buildGetPictureDataOption(req getFileRequest, book *model.Book) fileutil.GetPictureDataOption {
	return fileutil.GetPictureDataOption{
		PictureName:      req.fileName,
		BookID:           book.BookID,
		BookIsPDF:        book.Type == model.TypePDF,
		BookIsDir:        book.Type == model.TypeDir,
		BookIsNonUTF8Zip: book.NonUTF8Zip,
		BookPath:         book.BookPath,
		Debug:            config.GetCfg().Debug,
		UseCache:         config.GetCfg().UseCache,
		ResizeWidth:      req.resizeWidth,
		ResizeHeight:     req.resizeHeight,
		ResizeMaxWidth:   req.resizeMaxWidth,
		ResizeMaxHeight:  req.resizeMaxHeight,
		AutoCrop:         req.autoCrop,
		Gray:             req.gray,
		BlurHash:         req.blurhash,
		BlurHashImage:    req.blurhashImage,
		// 远程书籍支持
		IsRemote:  book.IsRemote,
		RemoteURL: book.RemoteURL,
	}
}

func savePictureCache(c echo.Context, req getFileRequest, book *model.Book, imgData []byte, contentType string) {
	if !config.GetCfg().UseCache || req.disableCache || book.Type == model.TypeDir {
		return
	}
	if err := fileutil.SaveFileToCache(
		req.bookID,
		req.fileName,
		imgData,
		fileutil.GetQueryString(c.Request().URL.Query()),
		contentType,
		config.GetCfg().CacheDir,
		config.GetCfg().Debug,
	); err != nil {
		logger.Infof(locale.GetString("log_save_file_to_cache_error"), err)
	}
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
