package common

import (
	"fmt"
	"github.com/yumenaka/comigo/entity"
	"github.com/yumenaka/comigo/htmx/state"
	"strconv"
)

// ServerHostBindStr  传递给前端，现实QRCode用的“主机域名”字符串
func ServerHostBindStr(serverHost string) string {
	//"{ serverHost: 'abc.com' }"
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
	if BookID == "" {
		return "/"
	}
	for _, book := range state.Global.TopBooks.BookInfos {
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

//// ShowContentAPIHandler htmx：一个用于显示内容的 API，未使用 htmx 时返回 HTTP 400 错误。
//func ShowContentAPIHandler(c *gin.Context) {
//	// 检查当前请求是否有 'HX-Request' 头部。
//	// 更多信息请见 https://htmx.org/docs/#request-headers”
//	if !htmx.IsHTMX(c.Request) {
//		// If not, return HTTP 400 error.
//		err := c.AbortWithError(http.StatusBadRequest, errors.New("non-htmx request"))
//		if err != nil {
//			log.Println(err)
//		}
//		return
//	}
//	// 编写 HTML内容。
//	_, err := c.Writer.Write([]byte("<p>🎉 Yes, <strong>htmx</strong> is ready to use! (<code>GET /api/hello-world</code>)</p>"))
//	if err != nil {
//		log.Println(err)
//	}
//
//	// 发送 htmx 响应。
//	err = htmx.NewResponse().Write(c.Writer)
//	if err != nil {
//		log.Println(err)
//	}
//}
