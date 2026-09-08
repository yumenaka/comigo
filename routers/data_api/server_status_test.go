package data_api

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
)

// 状态响应保持旧字段并添加插件与网页共同消费的地址及流量快照。
func TestServerStatusExtendedFields(t *testing.T) {
	originalStore := model.IStore
	t.Cleanup(func() { model.IStore = originalStore })
	model.IStore = &store.StoreInRam{}
	originalScanner := config.GlobalLibraryScanner
	t.Cleanup(func() { config.GlobalLibraryScanner = originalScanner })
	config.InitLibraryScanner()
	e := echo.New()
	r := httptest.NewRecorder()
	if err := GetServerInfoHandler(e.NewContext(httptest.NewRequest("GET", "/api/server", nil), r)); err != nil {
		t.Fatal(err)
	}
	var result map[string]json.RawMessage
	if err := json.Unmarshal(r.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"Version", "NumberOfBooks", "connections", "readingURL", "localBrowserURL", "localIPs", "traffic", "update", "externalAccess", "listenAddress"} {
		if _, ok := result[key]; !ok {
			t.Errorf("missing %s", key)
		}
	}
	var reading string
	if err := json.Unmarshal(result["readingURL"], &reading); err != nil || reading != config.GetQrcodeURL() {
		t.Fatal("reading URL mismatch")
	}
	var traffic map[string]json.RawMessage
	json.Unmarshal(result["traffic"], &traffic)
	for _, key := range []string{"startedAt", "receivedBytes", "sentBytes", "receiveBytesPerSecond", "sendBytesPerSecond"} {
		if _, ok := traffic[key]; !ok {
			t.Errorf("missing traffic %s", key)
		}
	}
}

// 高频接口只包含流量数据，与完整状态使用同一个计数器。
func TestServerTrafficSnapshot(t *testing.T) {
	e := echo.New()
	r := httptest.NewRecorder()
	if err := GetServerTrafficHandler(e.NewContext(httptest.NewRequest("GET", "/api/server/traffic", nil), r)); err != nil {
		t.Fatal(err)
	}
	var result map[string]json.RawMessage
	if err := json.Unmarshal(r.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if len(result) != 5 || result["startedAt"] == nil || result["sentBytes"] == nil {
		t.Fatalf("unexpected snapshot %s", r.Body.String())
	}
	if r.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("traffic must not be cached")
	}
}
