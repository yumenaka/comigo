package data_api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
)

// 验证书库 ID 必须来自配置，详情只含目标书库；远程路径不伪装成本地目录。
func TestStoreResourceScope(t *testing.T) {
	oldConfig, oldStore := config.CopyCfg(), model.IStore
	defer func() { *config.GetCfg() = oldConfig; model.IStore = oldStore }()
	local, other := t.TempDir(), t.TempDir()
	remote := "https://example.invalid/library"
	config.GetCfg().StoreUrls = []string{local, remote}
	model.IStore = &store.StoreInRam{}
	for _, info := range []model.BookInfo{
		{BookID: "first", StoreUrl: local, Type: model.TypeZip},
		{BookID: "second", StoreUrl: other, Type: model.TypeZip},
	} {
		if err := model.IStore.StoreBook(&model.Book{BookInfo: info}); err != nil {
			t.Fatal(err)
		}
	}
	id := base64.RawURLEncoding.EncodeToString([]byte(local + "/."))
	if _, err := StoreURLFromID(base64.RawURLEncoding.EncodeToString([]byte(other))); err != echo.ErrNotFound {
		t.Fatalf("unconfigured store: %v", err)
	}
	e := echo.New()
	rec := httptest.NewRecorder()
	ctx := e.NewContext(httptest.NewRequest(http.MethodGet, "/api/stores/"+id, nil), rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues(id)
	if err := GetStore(ctx); err != nil {
		t.Fatal(err)
	}
	var body struct {
		Store StoreInfo    `json:"store"`
		Books []model.Book `json:"books"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Store.BookCount != 1 || len(body.Books) != 1 || body.Books[0].BookID != "first" {
		t.Fatalf("unexpected scope: %s", rec.Body.String())
	}
	info := storeInfo(remote, nil)
	if !info.Remote || info.URL != remote || info.Exists != nil {
		t.Fatalf("remote store: %+v", info)
	}
}
