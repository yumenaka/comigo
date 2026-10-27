//go:build wails && !js

package wails_systray

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
)

// 原生绑定只能生成已知地址的二维码，且不能把完整配置或登录密码交给页面。
func TestPanelSnapshot(t *testing.T) {
	oldStore := model.IStore
	t.Cleanup(func() { model.IStore = oldStore })
	model.IStore = &store.StoreInRam{}
	old := config.CopyCfg()
	t.Cleanup(func() { config.ReplaceConfig(old) })
	next := config.CopyCfg()
	next.TemporaryReaderMode, next.DisableLAN = true, true
	next.Password = "panel-test-secret"
	next.StoreUrls = []string{"https://reader:library-secret@example.com/books?token=private-token"}
	config.ReplaceConfig(next)
	state, err := (&PanelService{}).Snapshot("https://untrusted.example/")
	if err != nil {
		t.Fatal(err)
	}
	if len(state.URLs) != 1 || state.ReadingURL != config.GetQrcodeURL() || !strings.HasPrefix(state.QR, "data:image/png;base64,") {
		t.Fatalf("错误的二维码状态: %+v", state.URLs)
	}
	data, err := json.Marshal(state)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), next.Password) || strings.Contains(string(data), "untrusted.example") || strings.Contains(string(data), "library-secret") || strings.Contains(string(data), "private-token") {
		t.Fatal("快照包含敏感字段或未经验证的地址")
	}
}

// 语言操作必须保存配置，并在只读模式下拒绝修改；无托盘宿主不应影响保存结果。
func TestPanelLanguageAndReadOnly(t *testing.T) {
	old := config.CopyCfg()
	t.Cleanup(func() { config.ReplaceConfig(old); locale.InitLanguageFromConfig(old.Language) })
	next := config.CopyCfg()
	next.TemporaryReaderMode, next.ReadOnlyMode = true, false
	config.ReplaceConfig(next)
	panel := &PanelService{tray: &Tray{}}
	if err := panel.SetLanguage("ja-JP"); err != nil {
		t.Fatal(err)
	}
	if config.GetCfg().Language != "ja-JP" {
		t.Fatal("语言未写回配置")
	}
	next = config.CopyCfg()
	next.ReadOnlyMode = true
	config.ReplaceConfig(next)
	for _, lang := range []string{"en-US", "unsupported"} {
		if err := panel.SetLanguage(lang); err == nil {
			t.Fatalf("错误接受语言修改: %s", lang)
		}
	}
	if err := panel.Action("rescan", 0); err == nil {
		t.Fatal("只读模式允许重新扫描")
	}
	if config.GetCfg().Language != "ja-JP" {
		t.Fatal("只读配置发生变化")
	}
}

// 绑定拒绝任意命令与越界索引，不访问原生浏览器或剪贴板。
func TestPanelRejectsInvalidActions(t *testing.T) {
	panel := &PanelService{}
	for _, input := range []struct {
		action string
		index  int
	}{{"exec", 0}, {"store", -1}, {"store", len(config.GetCfg().StoreUrls)}, {"copy", -1}} {
		if err := panel.Action(input.action, input.index); err == nil {
			t.Fatalf("错误接受操作: %+v", input)
		}
	}
}

// 面板资源使用本地 Wails 绑定，不引入远程脚本或 HTTP 控制请求。
func TestPanelAsset(t *testing.T) {
	rec := httptest.NewRecorder()
	ServePanel(rec, httptest.NewRequest(http.MethodGet, "/panel.html", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "Call.ByName") {
		t.Fatal("面板资源未正确提供")
	}
	if strings.Contains(rec.Body.String(), "fetch(") || strings.Contains(rec.Body.String(), "<script src=\"http") {
		t.Fatal("面板不应通过网络执行控制操作")
	}
}
