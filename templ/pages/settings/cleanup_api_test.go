package settings

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
)

// TestCleanupRequiresConfirmation 验证只读模式与缺少预览指纹时都不会进入删除逻辑。
func TestCleanupRequiresConfirmation(t *testing.T) {
	old := config.CopyCfg()
	t.Cleanup(func() { config.ReplaceConfig(old) })
	for _, readOnly := range []bool{false, true} {
		cfg := old
		cfg.ReadOnlyMode = readOnly
		config.ReplaceConfig(cfg)
		e := echo.New()
		body := `{}`
		wantMessage := locale.GetString("cleanup_changed")
		if readOnly {
			body = `{"token":"nonempty-preview-token"}`
			wantMessage = locale.GetString("err_config_locked")
		}
		req := httptest.NewRequest(http.MethodPost, "/api/stores/cleanup", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		err := CleanupDataHandler(e.NewContext(req, httptest.NewRecorder()))
		httpErr, ok := err.(*echo.HTTPError)
		want := http.StatusBadRequest

		if !ok || httpErr.Code != want || httpErr.Message != wantMessage {
			t.Fatalf("readonly=%v: got %v, want %d", readOnly, err, want)
		}
	}
}
