package components

import (
	"fmt"
	"github.com/yumenaka/comigo/entity"
	"github.com/yumenaka/comigo/htmx/state"
)

func serverHostBindStr(serverHost string) string {
	//"{ serverHost: 'abc.com' }"
	return "{ serverHost: '" + serverHost + "' }"
}

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
