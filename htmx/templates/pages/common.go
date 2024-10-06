package pages

import (
	"fmt"
	"github.com/yumenaka/comigo/entity"
	"github.com/yumenaka/comigo/htmx/state"
	"strconv"
)

func getImageAlt(key int) string {
	return strconv.Itoa(key)
}

// 阅读或书架页面，返回按钮实际使用的链接
func getReturnUrl(BookID string) string {
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
