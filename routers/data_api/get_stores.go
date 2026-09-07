package data_api

import (
	"encoding/base64"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools"
)

// StoreInfo 书库资源；ID 使用配置 URL 的 base64url 编码，远程 URL 不转成本地路径。
type StoreInfo struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Name      string `json:"name"`
	Remote    bool   `json:"remote"`
	Exists    *bool  `json:"exists"`
	BookCount int    `json:"bookCount"`
}

// StoreURLFromID 只允许访问已配置的书库，不能通过刷新接口扫描任意路径。
func StoreURLFromID(id string) (string, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(id)
	if err != nil {
		return "", echo.ErrBadRequest
	}
	target, _, err := tools.NormalizeStoreURLForCompare(string(decoded))
	if err != nil {
		return "", echo.ErrBadRequest
	}
	// 同一个本地书库的相对路径和绝对路径应解析到同一资源。
	for _, storeURL := range config.GetCfg().StoreUrls {
		key, _, err := tools.NormalizeStoreURLForCompare(storeURL)
		if err == nil && key == target {
			return storeURL, nil
		}
	}
	return "", echo.ErrNotFound
}

// storeInfo 复用统一的路径规范化来统计本地和远程书库。
func storeInfo(storeURL string, books []*model.Book) StoreInfo {
	normalized, remote, _ := tools.NormalizeStoreURLForCompare(storeURL)
	info := StoreInfo{ID: base64.RawURLEncoding.EncodeToString([]byte(storeURL)), URL: storeURL, Name: filepath.Base(normalized), Remote: remote}
	if !remote {
		_, err := os.Stat(normalized)
		exists := err == nil
		info.Exists = &exists
	}
	for _, book := range books {
		key, _, _ := tools.NormalizeStoreURLForCompare(book.StoreUrl)
		if key == normalized && book.Type != model.TypeBooksGroup {
			info.BookCount++
		}
	}
	return info
}

// GetStores 返回配置中的书库集合。
func GetStores(c echo.Context) error {
	books, err := model.IStore.ListBooks()
	if err != nil {
		return err
	}
	stores := make([]StoreInfo, 0, len(config.GetCfg().StoreUrls))
	for _, storeURL := range config.GetCfg().StoreUrls {
		stores = append(stores, storeInfo(storeURL, books))
	}
	return c.JSON(http.StatusOK, echo.Map{"stores": stores})
}

// GetStore 返回指定书库的状态及书籍索引。
func GetStore(c echo.Context) error {
	storeURL, err := StoreURLFromID(c.Param("id"))
	if err != nil {
		return err
	}
	books, err := model.IStore.ListBooks()
	if err != nil {
		return err
	}
	selected := make([]*model.Book, 0)
	target, _, _ := tools.NormalizeStoreURLForCompare(storeURL)
	for _, book := range books {
		key, _, _ := tools.NormalizeStoreURLForCompare(book.StoreUrl)
		if key == target {
			selected = append(selected, book.CloneForView())
		}
	}
	return c.JSON(http.StatusOK, echo.Map{"store": storeInfo(storeURL, selected), "books": selected})
}
