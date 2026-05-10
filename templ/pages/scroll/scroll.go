package scroll

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/templ/common"
	"github.com/yumenaka/comigo/templ/pages/error_page"
	"github.com/yumenaka/comigo/tools/logger"
)

const (
	scrollLoadModeInfinite = "infinite"
	scrollLoadModePaged    = "paged"
	defaultScrollPageLimit = 32
	maxScrollPageLimit     = 512
)

// ScrollModeHandler 阅读界面（卷轴阅读）
func ScrollModeHandler(c echo.Context) error {
	// 图片排序方式
	sortBy := "default"
	sortBookBy, err := c.Cookie("ScrollSortBy")
	if err == nil {
		// c.Cookie("key") 没找到，那么就会取到空值（nil），
		// 没处判断就直接访问 .Value 属性，会导致空指针引用错误。
		sortBy = sortBookBy.Value
		// logger.Info("Scroll Mode Sort By:" + sortBy)
	}
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
	book.SortPages(sortBy)
	loadMode := parseScrollLoadMode(c)
	pageLimit := parseScrollPageLimit(c)
	pagedIndex := parseScrollPageIndex(c, loadMode)
	if loadMode == scrollLoadModePaged && pagedIndex > scrollTotalPages(book, pageLimit) {
		pagedIndex = scrollTotalPages(book, pageLimit)
	}

	// 定义模板主体内容。
	scrollPage := ScrollPage(c, book, loadMode, pagedIndex, pageLimit)
	// 拼接页面
	indexHtml := common.Html(
		c,
		scrollPage, // define body content
		[]string{
			"static/js/ws_sync.js",
			"static/js/flip_modules/interaction_utils.js",
			"static/js/scroll.js",
		},
	)
	// 渲染页面
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
		// 渲染失败，返回 HTTP 500 错误。
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}

func parseScrollLoadMode(c echo.Context) string {
	switch c.QueryParam("load") {
	case scrollLoadModePaged:
		return scrollLoadModePaged
	default:
		return scrollLoadModeInfinite
	}
}

func parseScrollPageLimit(c echo.Context) int {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 {
		return defaultScrollPageLimit
	}
	if limit > maxScrollPageLimit {
		return maxScrollPageLimit
	}
	return limit
}

func parseScrollPageIndex(c echo.Context, loadMode string) int {
	if loadMode != scrollLoadModePaged {
		return -1
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		return 1
	}
	return page
}

func scrollTotalPages(book *model.Book, pageLimit int) int {
	if book == nil || pageLimit < 1 {
		return 1
	}
	total := (book.PageCount + pageLimit - 1) / pageLimit
	if total < 1 {
		return 1
	}
	return total
}

// 跳转用分页链接 /scroll/4cTOjFm?load=paged&page=1&limit=32
func getScrollPaginationURL(book *model.Book, page int, pageLimit int) string {
	readURL := `/scroll/` + book.BookID + `?load=paged&page=` + strconv.Itoa(page) + `&limit=` + strconv.Itoa(pageLimit)
	// href="javascript:void(0)" 是“点了什么也不发生”的老式写法
	if page < 1 {
		return `javascript:showToast(i18next.t('hint_first_page'), 'warning');`
	}
	if page > scrollTotalPages(book, pageLimit) {
		return `javascript:showToast(i18next.t('hint_last_page'), 'warning')`
	}
	return readURL
}

// 自动书签脚本，同时更新当前页码
func intersectScript(pageIndex int) string {
	return fmt.Sprintf(`
	    $nextTick(() => {
		if ($store.scroll.loadMode !== 'paged') {
			return;
		}
	if(!loaded || counter < 1){
        return;
    }
    // 更新当前页码
    $store.global.nowPageNum = %d;
    if (loaded && !updateBookmarkCompleted) {
        $store.global.UpdateBookmark(
            {
                type: 'auto',
                bookId: book.id,
                pageIndex: %d,
            }
        );
        updateBookmarkCompleted = true;
    }
  })
	`, pageIndex, pageIndex)
}
