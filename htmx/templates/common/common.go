package common

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/entity"
	"github.com/yumenaka/comigo/htmx/state"
	"github.com/yumenaka/comigo/util/logger"
)

// ServerHostBindStr  传递给前端，现实QRCode用的“主机域名”字符串
func ServerHostBindStr(serverHost string) string {
	return "{ serverHost: '" + serverHost + "' }"
}

// GetPageTitle 获取页面标题
func GetPageTitle(bookID string) string {
	if bookID == "" {
		return "Comigo " + state.Global.Version
	}
	groupBook, err := entity.GetBookByID(bookID, "")
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
	info, err := entity.GetBookGroupInfoByChildBookID(BookID)
	if err != nil {
		fmt.Println("ParentBookInfo not found")
		return "/"
	}
	if info.Depth <= 0 {
		return "/"
	}
	return "/shelf/" + info.BookID
}

func AddQuery(c *gin.Context, key string, value string) string {
	// 获取当前请求的 URL
	currentUrl := c.Request.URL

	// 解析 URL 参数
	params := currentUrl.Query()

	// 使用 Set 方法替换或添加新的查询参数 key=value
	params.Set(key, value)

	// 将修改后的查询参数重新编码
	currentUrl.RawQuery = params.Encode()

	// 输出修改后的 URL
	return currentUrl.String()
}

func ShowQuickJumpBar(b *entity.Book) (QuickJumpBar bool) {
	_, err := entity.GetBookInfoListByParentFolder(b.ParentFolder, "")
	if err != nil {
		logger.Infof("%s", err)
		return false
	}
	return true
}

func QuickJumpBarBooks(b *entity.Book, readMode string) (list *entity.BookInfoList) {
	list, err := entity.GetBookInfoListByParentFolder(b.ParentFolder, "")
	if err != nil {
		logger.Infof("%s", err)
		return nil
	}
	return list
}
