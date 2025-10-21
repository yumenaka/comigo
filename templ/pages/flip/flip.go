package flip

import (
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/templ/common"
	"github.com/yumenaka/comigo/templ/pages/error_page"
	"github.com/yumenaka/comigo/tools/logger"
)

// PageHandler 阅读界面（翻页模式）
func PageHandler(c echo.Context) error {
	bookID := c.Param("id")
	logger.Info("Flip Mode Book ID:" + bookID)
	// 图片排序方式
	sortBy := "default"
	// c.Cookie("key") 没找到，那么就会取到空值（nil），没判断nil就直接访问 .Value 属性，会导致空指针引用错误。
	sortBookBy, err := c.Cookie("FlipSortBy")
	if err == nil {
		sortBy = sortBookBy.Value
	}
	// // 给cookie设置默认值
	// if err != nil {
	// 	sortBookBy.Value = "default"
	// 	cookie := new(http.Cookie)
	// 	cookie.Name = "FlipSortBy"
	// 	cookie.Value = sortBookBy.Value
	// 	cookie.MaxAge = 3600000
	// 	cookie.Path = "/"
	// 	cookie.Domain = domain
	// 	cookie.Secure = false
	// 	cookie.HttpOnly = true
	// 	c.SetCookie(cookie)
	// }
	// 读取url参数，获取书籍ID
	book, err := model.IStore.GetBookAndSort(bookID, sortBy)
	if err != nil {
		logger.Infof("GetBook: %v", err)
		// 没有找到书，显示 HTTP 404 错误
		indexHtml := common.Html(
			c,
			error_page.NotFound404(c),
			[]string{},
		)
		// 渲染 404 页面
		if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
			// 渲染失败，返回 HTTP 500 错误。
			return c.NoContent(http.StatusInternalServerError)
		}
		return nil
	}

	// // TODO：加密链接的时候，设置Secure为true
	// readingProgressStr.Value = `{"nowPageNum":0,"nowChapterNum":0,"readingTime":0}`
	// cookie := new(http.Cookie)
	// cookie.Name = "bookID:" + bookID
	// cookie.Value = readingProgressStr.Value
	// cookie.MaxAge = 60 * 60 * 24 * 356
	// cookie.Path = "/"
	// cookie.Domain = domain
	// cookie.Secure = false
	// cookie.HttpOnly = false
	// c.SetCookie(cookie)

	//// 如果是静态页面，就把所有图片都转为 base64 编码，嵌入到 HTML 里
	//// 适合导出为单个 HTML 文件的场景
	//// 但是会非常占用内存和 CPU，尤其是大文件
	//// 使用方法：在 URL 后面加上 ?static=true
	//// 比如：http://
	//if c.QueryParam("static") == "true" {
	//	newBook := *book
	//	for i, file := range newBook.Pages.Images {
	//		// 获取图片数据的选项
	//		option := fileutil.GetPictureDataOption{
	//			PictureName:      file.Name,
	//			BookIsPDF:        newBook.Type == model.TypePDF,
	//			BookIsDir:        newBook.Type == model.TypeDir,
	//			BookIsNonUTF8Zip: newBook.NonUTF8Zip,
	//			BookFilePath:     newBook.FilePath,
	//			Debug:            config.GetDebug(),
	//			UseCache:         config.GetUseCache(),
	//			ResizeWidth:      -1,
	//			ResizeHeight:     -1,
	//			ResizeMaxWidth:   -1,
	//			ResizeMaxHeight:  -1,
	//			ThumbnailMode:    false,
	//			AutoCrop:         -1,
	//			Gray:             false,
	//			BlurHash:         0,
	//			BlurHashImage:    0,
	//		}
	//
	//		// 获取图片数据
	//		imgData, _, err := fileutil.GetPictureData(option)
	//		if err != nil {
	//			logger.Infof("GetPictureData error: %s", err)
	//			return c.JSON(http.StatusBadRequest, map[string]string{"error": "GetPictureData error: " + err.Error()})
	//		}
	//		mimeType := mime.TypeByExtension(filepath.Ext(file.Name))
	//		if mimeType == "" {
	//			mimeType = "application/octet-stream"
	//		}
	//		newBook.Pages.Images[i].Url = "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(imgData)
	//	}
	//	// 翻页模式页面
	//	indexHtml := common.Html(
	//		c,
	//		FlipPage(c, &newBook),
	//		[]string{"script/flip.js", "script/flip_sketch.js"})
	//	// 渲染翻页模式阅读页面
	//	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
	//		// 如果渲染失败，返回 HTTP 500 错误
	//		return c.NoContent(http.StatusInternalServerError)
	//	}
	//}

	// 翻页模式页面
	indexHtml := common.Html(
		c,
		FlipPage(c, book),
		[]string{"script/flip.js", "script/flip_sketch.js"})
	// 渲染翻页模式阅读页面
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
		// 如果渲染失败，返回 HTTP 500 错误
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}
