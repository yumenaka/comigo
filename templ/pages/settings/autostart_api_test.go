package settings

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools/autostart"
)

// 自启动必须显式赋值；只读和无效输入不能触发系统操作，失败不能报成功。
func TestAutostartWrites(t *testing.T) {
	old := config.CopyCfg()
	oldStatus, oldSet := startupStatus, setStartup
	t.Cleanup(func() { *config.GetCfg() = old; startupStatus = oldStatus; setStartup = oldSet })
	calls, enabled := 0, false
	startupStatus = func() (autostart.State, error) {
		return autostart.State{Enabled: enabled, Supported: true, Scope: "user"}, nil
	}
	setStartup = func(value bool) error { calls++; enabled = value; return nil }
	request := func(body string) error {
		req := httptest.NewRequest("PUT", "/api/autostart", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		return SetAutostartHandler(echo.New().NewContext(req, httptest.NewRecorder()))
	}
	config.GetCfg().ReadOnlyMode = true
	if request(`{"enabled":true}`) == nil || calls != 0 {
		t.Fatal("read-only write reached service manager")
	}
	config.GetCfg().ReadOnlyMode = false
	for _, body := range []string{`{}`, `{"enabled":null}`, `{"enabled":"false"}`} {
		if request(body) == nil || calls != 0 {
			t.Fatal("invalid input reached service manager")
		}
	}
	for _, body := range []string{`{"enabled":true}`, `{"enabled":false}`} {
		if err := request(body); err != nil {
			t.Fatal(err)
		}
	}
	if calls != 2 || enabled {
		t.Fatal("explicit false was lost")
	}
	setStartup = func(bool) error { return errors.New("permission denied") }
	if request(`{"enabled":true}`) == nil {
		t.Fatal("service failure reported success")
	}
}
