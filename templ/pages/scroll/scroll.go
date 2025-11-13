package scroll

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/templ/common"
	"github.com/yumenaka/comigo/templ/pages/error_page"
	"github.com/yumenaka/comigo/tools/logger"
)

// ScrollModeHandler 阅读界面（卷轴模式）
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
	book.SortPages(sortBy)
	readMode := "infinite_scroll"
	// 读取分页索引
	pagedIndex := -1
	page := c.QueryParam("page")
	if page != "" {
		readMode = "paged_scroll"
		index, err := strconv.Atoi(page)
		if err == nil {
			pagedIndex = index
		}
	}

	startIndex, err := strconv.Atoi(c.QueryParam("start"))
	if err != nil {
		startIndex = 0
	}

	// 定义模板主体内容。
	scrollPage := ScrollPage(c, book, readMode, pagedIndex, startIndex)
	// 拼接页面
	indexHtml := common.Html(
		c,
		scrollPage, // define body content
		[]string{"script/scroll.js"},
	)
	// 渲染页面
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
		// 渲染失败，返回 HTTP 500 错误。
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}

// 跳转用分页链接 /scroll/4cTOjFm?page=1
func getScrollPaginationURL(book *model.Book, page int) string {
	readURL := `/scroll/` + book.BookID + `?page=` + strconv.Itoa(page)
	// href="javascript:void(0)" 是“点了什么也不发生”的老式写法
	if page < 1 {
		return `javascript:showToast(i18next.t('hint_first_page'), 'warning');`
	}
	if page > (book.PageCount/32 + 1) {
		return `javascript:showToast(i18next.t('hint_last_page'), 'warning')`
	}
	return readURL
}

// 自动书签脚本
func intersectScript(pageIndex int) string {
	return fmt.Sprintf(`
    $nextTick(() => {
	if(!loaded || counter < 1){
        return;
    }
	//console.log({loaded});
    if (loaded && !updateBookmarkCompleted) {
        $store.global.UpdateBookmark(
            {
                type: 'auto',
                bookId: book.id,
                pageIndex: %d,
                label: '自动书签'
            }
        );
        updateBookmarkCompleted = true;
    }
  })
	`, pageIndex)
}
