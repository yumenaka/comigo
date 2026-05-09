package file

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/logger"
)

// cacheKey 用于存储文件信息的键
type cacheKey struct {
	bookID      string
	queryString string
}

// 使用标准库的 sync.Map 作为并发安全的映射
var contentTypeMap sync.Map

// SaveFileToCache 保存文件到缓存，加快读取速度
func SaveFileToCache(id, filename string, data []byte, queryString, contentType string, cachePath string, debug bool) error {
	// 创建缓存目录
	cacheDir := filepath.Join(cachePath, id)
	err := os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		logger.Infof(locale.GetString("log_save_file_to_cache_error"), err)
		return err
	}
	// 特殊字符转义，避免文件名不合法
	escapedFilename := url.PathEscape(filename)
	// 写入文件
	filePath := filepath.Join(cacheDir, escapedFilename)
	err = os.WriteFile(filePath, data, 0o644)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_write_file_to_cache"), err)
		return err
	}
	// 将 ContentType 存入并发安全的映射
	key := cacheKey{bookID: id, queryString: queryString}
	contentTypeMap.Store(key, contentType)
	return nil
}

// coverCacheFilename 根据书籍 ID 和请求尺寸生成封面缓存文件名。
// 说明：同一本书会被不同页面以不同 resize_height 请求，文件名必须包含尺寸，避免小封面污染大封面缓存。
func coverCacheFilename(bookID string, resizeHeight int) string {
	if resizeHeight > 0 {
		return bookID + "_h" + strconv.Itoa(resizeHeight) + ".jpg"
	}
	return bookID + "_original.jpg"
}

// SaveCoverToLocal 保存封面文件到本地文件夹，加快读取速度
func SaveCoverToLocal(metaDataDir string, bookID string, resizeHeight int, data []byte) error {
	// 创建封面meta data目录
	err := os.MkdirAll(metaDataDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf(locale.GetString("log_failed_to_write_file_to_cache"), err)
	}
	// 写入文件
	filename := coverCacheFilename(bookID, resizeHeight)
	filePath := filepath.Join(metaDataDir, filename)
	err = os.WriteFile(filePath, data, 0o644)
	if err != nil {
		return fmt.Errorf(locale.GetString("log_failed_to_write_file_to_cache"), err)
	}
	return nil
}

// CoverFileCacheExists 检查封面文件是否存在
func CoverFileCacheExists(metaDataDir string, bookID string, resizeHeight int) bool {
	// 构建封面文件路径
	filename := coverCacheFilename(bookID, resizeHeight)
	filePath := filepath.Join(metaDataDir, filename)
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// GetCoverFromLocal 从本地文件夹读取封面文件
func GetCoverFromLocal(metaDataDir string, bookID string, resizeHeight int) ([]byte, error) {
	filename := coverCacheFilename(bookID, resizeHeight)
	filePath := filepath.Join(metaDataDir, filename)
	// 读取文件
	return os.ReadFile(filePath)
}

// GetQueryString 根据查询参数生成一个排序后的键字符串
func GetQueryString(query url.Values) string {
	// 复制一份查询参数，避免影响原始数据
	queryCopy := url.Values{}
	for key, values := range query {
		queryCopy[key] = values
	}
	// 使用 Encode 方法会对键进行排序，生成稳定的查询字符串
	return queryCopy.Encode()
}

// DeleteCoverCache 删除封面缓存文件
func DeleteCoverCache(metaDataDir string, bookID string) error {
	cacheFiles := []string{filepath.Join(metaDataDir, bookID+".jpg")}
	if matchedFiles, err := filepath.Glob(filepath.Join(metaDataDir, bookID+"_h*.jpg")); err == nil {
		cacheFiles = append(cacheFiles, matchedFiles...)
	}
	cacheFiles = append(cacheFiles, filepath.Join(metaDataDir, bookID+"_original.jpg"))

	var deleteErrs []error
	for _, filePath := range cacheFiles {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue
		}
		if err := os.Remove(filePath); err != nil {
			deleteErrs = append(deleteErrs, err)
		}
	}
	return errors.Join(deleteErrs...)
}

// DeleteBookCache 删除书籍的图片缓存目录
func DeleteBookCache(cachePath string, bookID string) error {
	cacheDir := filepath.Join(cachePath, bookID)
	// 检查目录是否存在
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		return nil // 目录不存在，无需删除
	}
	return os.RemoveAll(cacheDir)
}

// GetFileFromCache 从缓存读取文件，加快第二次访问的速度
func GetFileFromCache(id, filename, queryString string, cachePath string, debug bool) ([]byte, string, error) {
	// 从映射中读取 ContentType
	key := cacheKey{bookID: id, queryString: queryString}
	value, ok := contentTypeMap.Load(key)
	if !ok {
		if debug {
			logger.Infof(locale.GetString("log_content_type_not_found_in_cache"), key)
		}
		return nil, "", errors.New(locale.GetString("err_content_type_not_found"))
	}
	contentType, _ := value.(string)

	// 特殊字符转义，构建文件路径
	escapedFilename := url.PathEscape(filename)
	filePath := filepath.Join(cachePath, id, escapedFilename)

	// 读取文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		if debug {
			logger.Infof(locale.GetString("log_failed_to_read_file_from_cache"), err)
		}
		return nil, contentType, err
	}
	return data, contentType, nil
}
