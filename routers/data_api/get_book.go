package data_api

import (
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers/apiresp"
	"github.com/yumenaka/comigo/tools/comigo_remote"
	"github.com/yumenaka/comigo/tools/file"
	"github.com/yumenaka/comigo/tools/logger"
	"net/http"
)

// GetBook 相关参数：
// id：     书籍的ID，必须项目  &id=2b17a130
// sort_by：页面排序方法，可选	 &sort_by=filename
// 示例 URL： http://127.0.0.1:1234/api/books/{id}?sort_by=filename
func GetBook(c echo.Context) error {
	sortBy := c.QueryParam("sort_by")
	if sortBy == "" {
		sortBy = "default"
	}
	id := c.Param("id")
	if id == "" {
		return apiresp.BadRequest(c, "missing_param", "not set id param", map[string]string{"param": "id"})
	}
	if localBook, client, storeURL, ok, err := remoteComigoBookFromRequest(c, id); ok {
		if err != nil {
			logger.Infof("%s", err)
			return writeRemoteComigoError(c, err)
		}
		remoteBook, err := client.GetBook(localBook.RemoteBookID, sortBy)
		if err != nil {
			logger.Infof("%s", err)
			return writeRemoteComigoError(c, err)
		}
		return apiresp.Success(c, "ok", "book retrieved", comigo_remote.LocalizeBookInShelf(storeURL, remoteBook, localBook.RemoteShelfKey, localBook.RemoteShelfName))
	}
	model.ClearBookWhenStoreUrlNotExist(config.GetCfg().StoreUrls)
	model.ClearBookNotExist()
	// 获取书籍信息
	b, err := model.IStore.GetBook(id)
	if err != nil {
		logger.Infof("%s", err)
		return apiresp.Error(c, http.StatusNotFound, "book_not_found", "id not found", map[string]string{"id": id})
	}
	b = b.CloneForView()
	b.SortPages(sortBy)
	// 如果是epub文件，重新按照Epub信息排序
	if b.Type == model.TypeEpub && sortBy == "epub_info" {
		imageList, err := file.GetImageListFromEpubFile(b.BookPath)
		if err != nil {
			logger.Infof("%s", err)
			return apiresp.Success(c, "ok", "book retrieved", b)
		}
		b.SortPagesByImageList(imageList)
	}
	return apiresp.Success(c, "ok", "book retrieved", b)
}

// GetBooks 返回当前书籍索引；完整分页信息通过单本资源获取。
func GetBooks(c echo.Context) error {
	books, err := model.IStore.ListBooks()
	if err != nil {
		return err
	}
	infos := make([]model.BookInfo, 0, len(books))
	for _, book := range books {
		infos = append(infos, book.CloneForView().BookInfo)
	}
	return c.JSON(http.StatusOK, echo.Map{"books": infos})
}
