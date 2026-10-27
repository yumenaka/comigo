package epub

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/templ/common"
	"github.com/yumenaka/comigo/tools/comigo_remote"
)

// PageHandler 只提供 EPUB 阅读界面和原文件地址，解包与排版在浏览器独立完成。
func PageHandler(c echo.Context) error {
	source := ""
	if id := c.Param("id"); id != "" {
		book, remote, err := comigo_remote.ResolveBook(config.GetCfg().StoreUrls, c.QueryParam(comigo_remote.RemoteStoreQuery), id, "default", config.GetCfg().TimeoutLimitForScan)
		if !remote {
			book, err = model.IStore.GetBook(id)
		}
		if err != nil || book == nil || book.Type != model.TypeEpub {
			return echo.ErrNotFound
		}
		source = common.RawBookURL(book)
	}
	c.Response().Header().Set("Referrer-Policy", "no-referrer")
	return common.RenderHTML(c, common.Html(c, Page(c, source), []string{
		"static/lib/epub/jszip.min.js",
		"static/lib/epub/epub.min.js",
		"static/js/epub_reader.js",
	}))
}

// Redirect 将既有漫画阅读链接交给 EPUB 专用页面，保留远程书库参数。
func Redirect(c echo.Context) error {
	target := config.PrefixPath("/epub/" + c.Param("id"))
	if query := c.QueryParams().Encode(); query != "" {
		target += "?" + query
	}
	return c.Redirect(http.StatusTemporaryRedirect, target)
}
