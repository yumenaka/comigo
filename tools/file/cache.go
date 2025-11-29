package file

import (
	"errors"
	"net/url"
	"os"
	"path"
	"path/filepath"
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
func SaveFileToCache(id, filename string, data []byte, queryString, contentType string, isCover bool, cachePath string, debug bool) error {
	// 创建缓存目录
	cacheDir := filepath.Join(cachePath, id)
	err := os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		logger.Infof(locale.GetString("saveFileToCache_error")+": %v", err)
		return err
	}

	// 特殊字符转义，避免文件名不合法
	escapedFilename := url.PathEscape(filename)
	// 如果是封面，另存为固定文件名
	if isCover {
		escapedFilename = "comigo_cover" + path.Ext(filename)
	}

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

// GetFileFromCache 从缓存读取文件，加快第二次访问的速度
func GetFileFromCache(id, filename, queryString string, isCover bool, cachePath string, debug bool) ([]byte, string, error) {
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
	if isCover {
		escapedFilename = "comigo_cover" + path.Ext(filename)
	}
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
