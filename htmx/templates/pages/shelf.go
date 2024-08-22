package pages

import (
	"github.com/angelofallars/htmx-go"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/entity"
	"github.com/yumenaka/comi/htmx/state"
	"github.com/yumenaka/comi/htmx/templates/components"
	"github.com/yumenaka/comi/util/logger"
	"net/http"
)

// ShelfHandler 书架页面的处理程序。
func ShelfHandler(c *gin.Context) {
	//书籍排列的方式，默认name
	sortBy := c.DefaultQuery("sort_by", "default")

	// 获取书架信息。
	id := c.Param("id")
	var err error
	if id == "" {
		// 获取顶层书架信息。
		state.Global.TopBooks, err = entity.TopOfShelfInfo(sortBy)
		if err != nil {
			logger.Infof("TopOfShelfInfo: %v", err)
			//TODO: 处理没有图书的情况（上传压缩包或远程下载示例漫画）
		}
	}
	if id != "" {
		// 通过书架ID获取书架信息。
		state.Global.NowBookID = id
		state.Global.TopBooks, err = entity.GetBookInfoListByID(id, sortBy)
		if err != nil {
			logger.Infof("GetBookShelf: %v", err)
		}
	}

	// 网页meta标签。
	metaTags := components.MetaTags(
		"Comigo  Comic Manga Reader 在线漫画 阅读器",         // define meta keywords
		"Simple Manga Reader in Linux，Windows，Mac OS", // define meta description
	)

	// 为首页定义模板布局。
	indexTemplate := components.MainLayout(
		"Comigo "+state.Global.Version, // define title text
		metaTags,                       // define meta tags
		ShelfPage(&state.Global),       // define body content
	)

	// 渲染索引页模板。
	if err := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, indexTemplate); err != nil {
		// 如果不是，返回 HTTP 500 错误。
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
