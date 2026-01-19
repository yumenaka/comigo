package player

import (
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/templ/common"
	"github.com/yumenaka/comigo/templ/pages/error_page"
	"github.com/yumenaka/comigo/tools/logger"
)

// PlayerModeHandler 播放器页面处理器
func PlayerModeHandler(c echo.Context) error {
	// 读取url参数，获取书籍ID
	bookID := c.Param("id")
	book, err := model.IStore.GetBook(bookID)
	if err != nil {
		logger.Infof(locale.GetString("log_getbook_error_scroll"), err)
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

	// 获取同文件夹下的所有书籍作为播放列表
	playlist := common.QuickJumpBarBooks(book)

	// 过滤出音视频文件
	var mediaList model.BookInfos
	if playlist != nil {
		for _, item := range *playlist {
			if item.Type == model.TypeVideo || item.Type == model.TypeAudio {
				mediaList = append(mediaList, item)
			}
		}
	}

	// 定义模板主体内容。
	playerPage := PlayerPage(c, book, &mediaList)
	// 拼接页面
	indexHtml := common.Html(
		c,
		playerPage, // define body content
		[]string{"script/player.js"},
	)
	// 渲染页面
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
		// 渲染失败，返回 HTTP 500 错误。
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}
