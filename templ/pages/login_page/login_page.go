package login_page

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/templ/common"
)

func getOAuthLoginButtonText() string {
	providerName := strings.TrimSpace(config.GetCfg().OAuthProviderName)
	if providerName == "" {
		return "i18next.t('login_with_oauth')"
	}
	return "i18next.t('login_with_oauth_provider', { provider: " + strconv.Quote(providerName) + " })"
}

// Handler 上传文件页面
func Handler(c echo.Context) error {
	indexHtml := common.Html(
		c,
		LoginPage(),
		[]string{},
	)
	// 渲染页面
	if err := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexHtml); err != nil {
		// 渲染失败，返回 HTTP 500 错误。
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}
