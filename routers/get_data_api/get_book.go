package get_data_api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/file"
	"github.com/yumenaka/comigo/tools/logger"
)

// GetBook 相关参数：
// id：     书籍的ID，必须项目  &id=2b17a130
// sort_by：页面排序方法，可选	 &sort_by=filename
// 示例 URL： http://127.0.0.1:1234/api/get_book?id=1215a&sort_by=filename
func GetBook(c echo.Context) error {
	sortBy := c.QueryParam("sort_by")
	if sortBy == "" {
		sortBy = "default"
	}
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "not set id param")
	}
	model.IStore.ClearBookNotExist()
	// 获取书籍信息
	b, err := model.IStore.GetBookAndSort(id, sortBy)
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
