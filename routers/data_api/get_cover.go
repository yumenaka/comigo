package data_api

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	taglib "github.com/dhowden/tag"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers/apiresp"
	"github.com/yumenaka/comigo/tools"
	fileutil "github.com/yumenaka/comigo/tools/file"
	"github.com/yumenaka/comigo/tools/logger"
)

type coverRequest struct {
	bookID       string
	resizeHeight int
}

type coverCacheContext struct {
	configDir string
	metaPath  string
}

// GetCover 获取书籍封面
// 相关参数：
// id：书籍的ID，必须参数 &id=2B17a
// resize_height：可选参数，指定封面高度，默认值为352 &resize_height=500
// 示例 URL： http://127.0.0.1:1234/api/get-cover?id=2b17a13
// 示例 URL（自定义高度）： http://127.0.0.1:1234/api/get-cover?id=2b17a13&resize_height=500
func GetCover(c echo.Context) error {
	req, err := parseCoverRequest(c)
	if err != nil {
		return err
	}
	book, err := model.IStore.GetBook(req.bookID)
	if err != nil {
		logger.Infof(locale.GetString("log_getbook_error_common"), err)
		return apiresp.Error(c, http.StatusNotFound, "book_not_found", "Book not found", map[string]string{"id": req.bookID})
	}

	// 获取封面信息
	cover := book.GetCover()
	cacheCtx := getCoverCacheContext(book)
	if handled, err := serveCachedCover(c, cacheCtx, req.bookID); handled || err != nil {
		return err
	}
	// 尝试从 MP3 的 ID3(APIC) 内嵌图片中读取封面（在本地缓存未命中时）
	// 说明：
	// - 优先使用第三方库提高兼容性（dhowden/tag）
	// - 成功读取后统一转为 JPEG，并复用现有的封面缓存（bookID.jpg）
	if handled, err := serveMP3Cover(c, book, cacheCtx, req); handled || err != nil {
		return err
	}

	needFile := cover.Name
	coverURLPath := config.StripBasePath(cover.Url)
	// 处理 TypeBooksGroup 的情况：封面来自子书籍
	coverBook := book
	if book.Type == model.TypeBooksGroup && strings.HasPrefix(coverURLPath, "/api/get-file") {
		coverBook, needFile, err = resolveBookGroupCover(c, coverURLPath, cover.Name)
		if err != nil {
			return err
		}
	}

	// 如果封面URL是内嵌图片, 通过封面文件获取
	if strings.HasPrefix(coverURLPath, "/images/") {
		return serveEmbeddedCover(c, cover.Name, cover.Url, cacheCtx, req)
	}

	// 获取图片数据
	imgData, contentType, err := fileutil.GetPictureData(buildCoverPictureDataOption(coverBook, needFile, req.resizeHeight))
	if err != nil {
		logger.Infof(locale.GetString("log_get_file_error"), err)
		return apiresp.BadRequest(c, "get_cover_failed", "Get file error: "+err.Error(), nil)
	}
	// 缓存封面到 configDir
	saveCoverCache(cacheCtx, req.bookID, imgData)
	// 返回图片数据
	return serveCoverBytes(c, contentType, imgData)
}

func parseCoverRequest(c echo.Context) (coverRequest, error) {
	req := coverRequest{
		bookID:       c.QueryParam("id"),
		resizeHeight: getIntQueryParam(c, "resize_height", 352),
	}
	if req.bookID == "" {
		return req, apiresp.BadRequest(c, "missing_param", "id is required", map[string]string{"param": "id"})
	}
	return req, nil
}

func getCoverCacheContext(book *model.Book) coverCacheContext {
	configDir, err := config.GetConfigDir()
	if err != nil {
		logger.Errorf(locale.GetString("err_failed_to_get_config_dir"), err)
		return coverCacheContext{}
	}
	return coverCacheContext{
		configDir: configDir,
		metaPath:  filepath.Join(configDir, "metadata", book.GetStoreID()),
	}
}

func serveCachedCover(c echo.Context, cacheCtx coverCacheContext, bookID string) (bool, error) {
	if cacheCtx.metaPath == "" || !fileutil.CoverFileCacheExists(cacheCtx.metaPath, bookID) {
		return false, nil
	}
	coverData, err := fileutil.GetCoverFromLocal(cacheCtx.metaPath, bookID)
	if err != nil {
		return false, nil
	}
	return true, serveCoverBytes(c, "image/jpeg", coverData)
}

func serveMP3Cover(c echo.Context, book *model.Book, cacheCtx coverCacheContext, req coverRequest) (bool, error) {
	if book.Type != model.TypeAudio || !strings.EqualFold(filepath.Ext(book.BookPath), ".mp3") {
		return false, nil
	}
	imgData, err := extractMP3CoverByTag(book.BookPath)
	if err != nil || len(imgData) == 0 {
		return false, nil
	}
	imgData = tools.ImageResizeByHeight(imgData, req.resizeHeight)
	saveCoverCache(cacheCtx, req.bookID, imgData)
	return true, serveCoverBytes(c, "image/jpeg", imgData)
}

func resolveBookGroupCover(c echo.Context, coverURLPath string, fallbackName string) (*model.Book, string, error) {
	// URL 格式：/api/get-file?id=子书籍ID&filename=文件名。
	parsedURL, err := url.Parse(coverURLPath)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_parse_cover_url"), err)
		return nil, "", apiresp.BadRequest(c, "invalid_cover_url", "Invalid cover URL", nil)
	}
	childID := parsedURL.Query().Get("id")
	if childID == "" {
		logger.Infof(locale.GetString("log_child_book_id_missing_in_cover_url"))
		return nil, "", apiresp.BadRequest(c, "missing_child_book_id", "Child book ID is required in cover URL", nil)
	}
	childBook, err := model.IStore.GetBook(childID)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_get_child_book"), err)
		return nil, "", apiresp.Error(c, http.StatusNotFound, "child_book_not_found", "Child book not found", map[string]string{"id": childID})
	}
	needFile := fallbackName
	if filename := parsedURL.Query().Get("filename"); filename != "" {
		needFile = filename
	}
	return childBook, needFile, nil
}

func serveEmbeddedCover(c echo.Context, coverName string, coverURL string, cacheCtx coverCacheContext, req coverRequest) error {
	imgData := assets.GetImageData(coverName)
	if len(imgData) == 0 {
		logger.Infof(locale.GetString("log_failed_to_read_embedded_image"), coverURL)
		return apiresp.Error(c, http.StatusNotFound, "embedded_image_not_found", "Embedded image not found", nil)
	}
	imgData = tools.ImageResizeByHeight(imgData, req.resizeHeight)
	contentType := tools.GetContentTypeByFileName(".jpg")
	saveCoverCache(cacheCtx, req.bookID, imgData)
	return serveCoverBytes(c, contentType, imgData)
}

func buildCoverPictureDataOption(book *model.Book, needFile string, resizeHeight int) fileutil.GetPictureDataOption {
	return fileutil.GetPictureDataOption{
		PictureName:      needFile,
		BookID:           book.BookID,
		BookPath:         book.BookPath,
		BookIsPDF:        book.Type == model.TypePDF,
		BookIsDir:        book.Type == model.TypeDir,
		BookIsNonUTF8Zip: book.NonUTF8Zip,
		Debug:            config.GetCfg().Debug,
		ResizeHeight:     resizeHeight,
		// 远程书籍支持
		IsRemote:  book.IsRemote,
		RemoteURL: book.RemoteURL,
	}
}

func saveCoverCache(cacheCtx coverCacheContext, bookID string, imgData []byte) {
	if cacheCtx.configDir == "" || cacheCtx.metaPath == "" {
		return
	}
	if err := fileutil.SaveCoverToLocal(cacheCtx.metaPath, bookID, imgData); err != nil {
		logger.Infof(locale.GetString("log_save_cover_to_local_error"), err)
	}
}

func serveCoverBytes(c echo.Context, contentType string, imgData []byte) error {
	return c.Blob(http.StatusOK, contentType, imgData)
}

// extractMP3CoverByTag 使用 github.com/dhowden/tag 从 MP3 中提取内嵌封面（APIC）
func extractMP3CoverByTag(mp3Path string) ([]byte, error) {
	f, err := os.Open(mp3Path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	meta, err := taglib.ReadFrom(f)
	if err != nil {
		return nil, err
	}
	pic := meta.Picture()
	if pic == nil || len(pic.Data) == 0 {
		return nil, errors.New("no picture")
	}
	return pic.Data, nil
}
