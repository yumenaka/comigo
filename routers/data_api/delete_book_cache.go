package data_api

import (
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
	fileutil "github.com/yumenaka/comigo/tools/file"
	"github.com/yumenaka/comigo/tools/logger"
)

// DeleteBookCache 删除书籍的元数据和缓存文件
// 示例 URL： http://127.0.0.1:1234/api/delete_book_cache?id=2b17a13
// 相关参数：
// id：书籍的ID，必须参数  &id=2b17a13
// delete_metadata: 是否删除元数据JSON文件，可选参数，默认 true  &delete_metadata=true
// delete_cover: 是否删除封面缓存文件，可选参数，默认 true  &delete_cover=true
// delete_image_cache: 是否删除图片缓存目录，可选参数，默认 true  &delete_image_cache=true
func DeleteBookCache(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	// 获取可选参数
	deleteMetadata := getBoolQueryParam(c, "delete_metadata", true)
	deleteCover := getBoolQueryParam(c, "delete_cover", true)
	deleteImageCache := getBoolQueryParam(c, "delete_image_cache", true)

	// 获取书籍信息
	book, err := model.IStore.GetBook(id)
	if err != nil {
		logger.Infof("GetBook error: %s", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Book not found"})
	}

	result := map[string]interface{}{
		"book_id": id,
		"deleted": map[string]bool{},
	}
	deletedMap := result["deleted"].(map[string]bool)

	// 获取配置目录
	configDir, err := config.GetConfigDir()
	if err != nil {
		logger.Infof("GetConfigDir error: %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get config directory"})
	}

	// 元数据目录路径
	metaPath := filepath.Join(configDir, "metadata", book.GetStoreID())

	// 删除元数据 JSON 文件
	if deleteMetadata {
		err := store.DeleteBookJson(book)
		if err != nil {
			logger.Infof("DeleteBookJson error: %s", err)
			deletedMap["metadata"] = false
		} else {
			deletedMap["metadata"] = true
		}
	}

	// 删除封面缓存文件
	if deleteCover {
		err := fileutil.DeleteCoverCache(metaPath, id)
		if err != nil {
			logger.Infof("DeleteCoverCache error: %s", err)
			deletedMap["cover"] = false
		} else {
			deletedMap["cover"] = true
		}
	}

	// 删除图片缓存目录
	if deleteImageCache && config.GetCfg().CacheDir != "" {
		err := fileutil.DeleteBookCache(config.GetCfg().CacheDir, id)
		if err != nil {
			logger.Infof("DeleteBookCache error: %s", err)
			deletedMap["image_cache"] = false
		} else {
			deletedMap["image_cache"] = true
		}
	}

	return c.JSON(http.StatusOK, result)
}
