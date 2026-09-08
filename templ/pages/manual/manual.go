package manual

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/templ/common"
)

var validPages = map[string]bool{
	"index": true, "install": true, "quick-start": true, "reading": true,
	"library": true, "desktop": true, "comigo-omarchy": true, "deployment": true, "faq": true,
}

// Handler 渲染公开的内置手册；正文由独立 JSON 在浏览器中载入。
func Handler(c echo.Context) error {
	language, page, ok := parsePath(c.Request().URL.Path)
	if !ok {
		return echo.ErrNotFound
	}
	return common.RenderHTML(c, ManualDocument(c, language, page))
}

// parsePath 校验固定的手册深链接，避免未知路径返回一个空白的 200 页面。
func parsePath(path string) (language, page string, ok bool) {
	parts := strings.FieldsFunc(strings.TrimPrefix(path, "/manual"), func(r rune) bool { return r == '/' })
	language, page = "zh", "index"
	if len(parts) > 0 && (parts[0] == "ja-JP" || parts[0] == "en-US") {
		language, parts = parts[0], parts[1:]
	}
	if len(parts) > 1 {
		return "", "", false
	}
	if len(parts) == 1 {
		page = parts[0]
	}
	if !validPages[page] {
		return "", "", false
	}
	return language, page, true
}
