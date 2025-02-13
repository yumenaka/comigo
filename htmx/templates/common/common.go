package common

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/htmx/state"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/logger"
)

// ServerHostBindStr  传递给前端，显示QRCode用的"主机域名"字符串
func ServerHostBindStr(serverHost string) string {
	return "{ serverHost: '" + serverHost + "' }"
}

// GetPageTitle 获取页面标题
func GetPageTitle(bookID string) string {
	if bookID == "" {
		return "Comigo " + state.Global.Version
	}
	groupBook, err := model.GetBookByID(bookID, "")
	if err != nil {
		fmt.Printf("GetBookByID: %v", err)
		return "Comigo " + state.Global.Version
	}
	return groupBook.Title
}

// GetImageAlt 获取图片的 alt 属性
func GetImageAlt(key int) string {
	return strconv.Itoa(key)
}

// GetReturnUrl 阅读或书架页面，返回按钮实际使用的链接
func GetReturnUrl(BookID string) string {
	if BookID == "" || state.Global.ShelfBookList == nil {
		return "/"
	}
	for _, book := range state.Global.ShelfBookList.BookInfos {
		if book.BookID == BookID {
			return "/"
		}
	}
	// 如果是书籍组，就跳转到子书架
	info, err := model.GetBookGroupInfoByChildBookID(BookID)
	if err != nil {
		fmt.Println("ParentBookInfo not found")
		return "/"
	}
	if info.Depth <= 0 {
		return "/"
	}
	return "/shelf/" + info.BookID
}

func AddQuery(c echo.Context, key, value string) string {
	// 获取当前请求的URL
	req := c.Request()
	currentURL := req.URL

	// 解析URL参数
	params := currentURL.Query()

	// 替换或添加新的查询参数 key=value
	params.Set(key, value)

	// 将修改后的查询参数重新编码
	currentURL.RawQuery = params.Encode()

	// 返回修改后的URL
	return currentURL.String()
}

func ShowQuickJumpBar(b *model.Book) (QuickJumpBar bool) {
	_, err := model.GetBookInfoListByParentFolder(b.ParentFolder, "")
	if err != nil {
		logger.Infof("%s", err)
		return false
	}
	return true
}

func QuickJumpBarBooks(b *model.Book) (list *model.BookInfoList) {
	list, err := model.GetBookInfoListByParentFolder(b.ParentFolder, "")
	if err != nil {
		logger.Infof("%s", err)
		return nil
	}
	return list
}
