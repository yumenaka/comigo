package data_api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
)

// StoreInfo 书库信息
type StoreInfo struct {
	URL    string `json:"url"`    // 书库路径
	Name   string `json:"name"`   // 书库名称
	Exists bool   `json:"exists"` // 路径是否存在
}

// GetStores 获取书库列表的 API 处理函数
func GetStores(c echo.Context) error {
	storeUrls := config.GetCfg().StoreUrls
	stores := make([]StoreInfo, 0, len(storeUrls))

	// 遍历所有配置的书库路径
	for _, url := range storeUrls {
		// 获取绝对路径
		absPath, err := filepath.Abs(url)
		if err != nil {
			absPath = url
		}

		// 检查路径是否存在
		exists := false
		if _, err := os.Stat(absPath); err == nil {
			exists = true
		}

		stores = append(stores, StoreInfo{
			URL:    absPath,
			Name:   filepath.Base(absPath),
			Exists: exists,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"stores": stores,
	})
}
