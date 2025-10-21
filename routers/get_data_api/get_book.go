package get_data_api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/file"
	"github.com/yumenaka/comigo/tools/logger"
)

// GetBook 相关参数：
// id：书籍的ID，必须项目       							&id=2b17a130
// author：书籍的作者，未必存在									&author=佚名
// sort_page：按照自然文件名重新排序							&sort_page=true
// 示例 URL： http://127.0.0.1:1234/api/get_book?id=1215a&sort_by=name
// 示例 URL： http://127.0.0.1:1234/api/get_book?&author=Doe&name=book_name
func GetBook(c echo.Context) error {
	//author := c.QueryParam("author")
	sortBy := c.QueryParam("sort_by")
	if sortBy == "" {
		sortBy = "default"
	}
	id := c.QueryParam("id")

	model.IStore.CheckAllNotExistBooks()

	//if author != "" {
	//	bookList, err := model.IStore.GetBookByAuthor(author, sortBy)
	//	if err != nil {
	//		logger.Infof("%s", err)
	//	}
	//	return c.JSON(http.StatusOK, bookList)
	//}

	if id != "" {
		b, err := model.IStore.GetBookByID(id, sortBy)
		if err != nil {
			logger.Infof("%s", err)
			return c.JSON(http.StatusBadRequest, "id not found")
		}

		// 如果是epub文件，重新按照Epub信息排序
		if b.Type == model.TypeEpub && sortBy == "epub_info" {
			imageList, err := file.GetImageListFromEpubFile(b.FilePath)
			if err != nil {
				logger.Infof("%s", err)
				return c.JSON(http.StatusOK, b)
			}
			b.SortPagesByImageList(imageList)
		}
		return c.JSON(http.StatusOK, b)
	}

	return c.JSON(http.StatusBadRequest, "no valid parameters provided")
}
