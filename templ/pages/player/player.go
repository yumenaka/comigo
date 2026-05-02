package player

import (
	"net/http"
	"path/filepath"

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

	// 过滤出当前目录下的音视频文件。
	mediaList := buildMediaPlaylist(book, playlist)

	// 定义模板主体内容。
	playerPage := PlayerPage(c, book, &mediaList)
	// 拼接页面
	indexHtml := common.Html(
		c,
		playerPage, // define body content
		[]string{"static/js/player.js"},
	)
	// 渲染页面
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
		// 渲染失败，返回 HTTP 500 错误。
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}

// buildMediaPlaylist 从同目录书籍列表中提取播放器列表。
// ParentFolder 只保存目录名，同名目录会混入候选列表；播放器须按真实路径再次过滤
func buildMediaPlaylist(current *model.Book, playlist *model.BookInfos) model.BookInfos {
	if playlist == nil {
		return nil
	}

	currentDir := ""
	currentParent := ""
	if current != nil {
		currentParent = current.ParentFolder
		if current.BookPath != "" {
			currentDir = filepath.Dir(current.BookPath)
		}
	}

	mediaList := make(model.BookInfos, 0, len(*playlist))
	for _, item := range *playlist {
		if item.Type != model.TypeVideo && item.Type != model.TypeAudio {
			continue
		}
		if !samePlayerFolder(currentDir, currentParent, item) {
			continue
		}
		mediaList = append(mediaList, item)
	}
	return mediaList
}

func samePlayerFolder(currentDir, currentParent string, item model.BookInfo) bool {
	if currentDir != "" && item.BookPath != "" {
		return filepath.Dir(item.BookPath) == currentDir
	}
	if currentParent != "" {
		return item.ParentFolder == currentParent
	}
	return true
}
