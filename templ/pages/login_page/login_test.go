package login_page

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// 登录错误来自查询参数和网络响应，必须始终作为文字而不是 HTML 展示。
func TestLoginErrorsRenderAsText(t *testing.T) {
	var output bytes.Buffer
	if err := LoginMainArea().Render(context.Background(), &output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	if !strings.Contains(html, `x-text="error"`) || strings.Contains(html, `x-html=`) {
		t.Fatal("login errors must not use an HTML sink")
	}
}
