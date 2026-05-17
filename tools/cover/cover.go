package cover

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	taglib "github.com/dhowden/tag"
	"github.com/yumenaka/comigo/assets"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools"
	fileutil "github.com/yumenaka/comigo/tools/file"
	"github.com/yumenaka/comigo/tools/logger"
)

var (
	ErrMissingBookID = errors.New("book id is required")
	ErrBookNotFound  = errors.New("book not found")
)

type Request struct {
	BookID       string
	ResizeHeight int
}

type Result struct {
	Data        []byte
	ContentType string
}

type cacheContext struct {
	configDir string
	metaPath  string
}

// GetBookCover 按书籍 ID 获取封面字节，供 HTTP 接口和 TUI 复用同一套封面解析逻辑。
func GetBookCover(req Request) (Result, error) {
	if req.BookID == "" {
		return Result{}, ErrMissingBookID
	}
	if req.ResizeHeight <= 0 {
		req.ResizeHeight = 352
	}

	book, err := model.IStore.GetBook(req.BookID)
	if err != nil {
		logger.Infof(locale.GetString("log_getbook_error_common"), err)
		return Result{}, fmt.Errorf("%w: %s", ErrBookNotFound, req.BookID)
	}

	bookCover := book.GetCover()
	cacheCtx := getCacheContext(book)
	if data, ok := getCached(cacheCtx, req); ok {
		return Result{Data: data, ContentType: "image/jpeg"}, nil
	}
	if data, ok := getMP3Cover(book, cacheCtx, req); ok {
		return Result{Data: data, ContentType: "image/jpeg"}, nil
	}

	needFile := bookCover.Name
	coverURLPath := config.StripBasePath(bookCover.Url)
	coverBook := book
	if book.Type == model.TypeBooksGroup && strings.HasPrefix(coverURLPath, "/api/get-file") {
		coverBook, needFile, err = resolveBookGroupCover(coverURLPath, bookCover.Name)
		if err != nil {
			return Result{}, err
		}
	}

	if strings.HasPrefix(coverURLPath, "/images/") {
		return getEmbeddedCover(bookCover.Name, bookCover.Url, cacheCtx, req)
	}

	imgData, contentType, err := fileutil.GetPictureData(buildPictureDataOption(coverBook, needFile, req.ResizeHeight))
	if err != nil {
		logger.Infof(locale.GetString("log_get_file_error"), err)
		return Result{}, err
	}
	saveCache(cacheCtx, req, imgData)
	return Result{Data: imgData, ContentType: contentType}, nil
}

// getCacheContext 根据书籍所属书库计算封面缓存目录；失败时返回空上下文并跳过缓存。
func getCacheContext(book *model.Book) cacheContext {
	configDir, err := config.GetConfigDir()
	if err != nil {
		logger.Errorf(locale.GetString("err_failed_to_get_config_dir"), err)
		return cacheContext{}
	}
	return cacheContext{
		configDir: configDir,
		metaPath:  filepath.Join(configDir, "metadata", book.GetStoreID()),
	}
}

// getCached 优先读取已生成的封面缓存，避免 TUI 频繁切换选中项时重复解包图片。
func getCached(cacheCtx cacheContext, req Request) ([]byte, bool) {
	if cacheCtx.metaPath == "" || !fileutil.CoverFileCacheExists(cacheCtx.metaPath, req.BookID, req.ResizeHeight) {
		return nil, false
	}
	coverData, err := fileutil.GetCoverFromLocal(cacheCtx.metaPath, req.BookID, req.ResizeHeight)
	if err != nil {
		return nil, false
	}
	return coverData, true
}

// getMP3Cover 处理音频文件内嵌封面；普通图片/漫画封面不会走这个分支。
func getMP3Cover(book *model.Book, cacheCtx cacheContext, req Request) ([]byte, bool) {
	if book.Type != model.TypeAudio || !strings.EqualFold(filepath.Ext(book.BookPath), ".mp3") {
		return nil, false
	}
	imgData, err := extractMP3CoverByTag(book.BookPath)
	if err != nil || len(imgData) == 0 {
		return nil, false
	}
	imgData = tools.ImageResizeByHeight(imgData, req.ResizeHeight)
	saveCache(cacheCtx, req, imgData)
	return imgData, true
}

// resolveBookGroupCover 解析书籍组封面 URL，定位真正提供封面的子书籍和文件名。
func resolveBookGroupCover(coverURLPath string, fallbackName string) (*model.Book, string, error) {
	parsedURL, err := url.Parse(coverURLPath)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_parse_cover_url"), err)
		return nil, "", fmt.Errorf("invalid cover url: %w", err)
	}
	childID := parsedURL.Query().Get("id")
	if childID == "" {
		logger.Infof(locale.GetString("log_child_book_id_missing_in_cover_url"))
		return nil, "", errors.New("child book id is required in cover url")
	}
	childBook, err := model.IStore.GetBook(childID)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_get_child_book"), err)
		return nil, "", fmt.Errorf("%w: %s", ErrBookNotFound, childID)
	}
	needFile := fallbackName
	if filename := parsedURL.Query().Get("filename"); filename != "" {
		needFile = filename
	}
	return childBook, needFile, nil
}

// getEmbeddedCover 读取内置占位封面，例如没有可用封面时的默认图片。
func getEmbeddedCover(coverName string, coverURL string, cacheCtx cacheContext, req Request) (Result, error) {
	imgData := assets.GetImageData(coverName)
	if len(imgData) == 0 {
		logger.Infof(locale.GetString("log_failed_to_read_embedded_image"), coverURL)
		return Result{}, errors.New("embedded image not found")
	}
	imgData = tools.ImageResizeByHeight(imgData, req.ResizeHeight)
	contentType := tools.GetContentTypeByFileName(".jpg")
	saveCache(cacheCtx, req, imgData)
	return Result{Data: imgData, ContentType: contentType}, nil
}

// buildPictureDataOption 将书籍模型转换为底层解包/读取图片接口需要的参数。
func buildPictureDataOption(book *model.Book, needFile string, resizeHeight int) fileutil.GetPictureDataOption {
	return fileutil.GetPictureDataOption{
		PictureName:      needFile,
		BookID:           book.BookID,
		BookPath:         book.BookPath,
		BookIsPDF:        book.Type == model.TypePDF,
		BookIsDir:        book.Type == model.TypeDir,
		BookIsNonUTF8Zip: book.NonUTF8Zip,
		Debug:            config.GetCfg().Debug,
		ResizeHeight:     resizeHeight,
		IsRemote:         book.IsRemote,
		RemoteURL:        book.RemoteURL,
	}
}

// saveCache 写入封面缓存；缓存失败只记录日志，不影响本次封面显示。
func saveCache(cacheCtx cacheContext, req Request, imgData []byte) {
	if cacheCtx.configDir == "" || cacheCtx.metaPath == "" {
		return
	}
	if err := fileutil.SaveCoverToLocal(cacheCtx.metaPath, req.BookID, req.ResizeHeight, imgData); err != nil {
		logger.Infof(locale.GetString("log_save_cover_to_local_error"), err)
	}
}

// extractMP3CoverByTag 使用 github.com/dhowden/tag 从 MP3 中提取内嵌封面（APIC）。
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
