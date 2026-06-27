//go:build wails

package routers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
)

// TestWailsAndroidFetchReplaysPostRequest 验证 Android fetch 桥不会丢失 method、query、headers 和 JSON body。
func TestWailsAndroidFetchReplaysPostRequest(t *testing.T) {
	cfg := config.GetCfg()
	oldBasePath := cfg.BasePath
	cfg.BasePath = ""
	t.Cleanup(func() {
		cfg.BasePath = oldBasePath
	})

	e := echo.New()
	api := e.Group("/api")
	bindWailsAPI(api)
	api.POST("/echo", func(c echo.Context) error {
		var body map[string]string
		if err := c.Bind(&body); err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, map[string]string{
			"method": c.Request().Method,
			"query":  c.QueryParam("q"),
			"header": c.Request().Header.Get("X-Test"),
			"body":   body["value"],
		})
	})

	rawPayload, err := json.Marshal(wailsAndroidFetchRequest{
		Method: http.MethodPost,
		Path:   "/api/echo?q=ok",
		Headers: map[string]string{
			echo.HeaderContentType: echo.MIMEApplicationJSON,
			"X-Test":               "from-bridge",
		},
		Body: `{"value":"kept"}`,
	})
	if err != nil {
		t.Fatalf("marshal bridge payload: %v", err)
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(rawPayload)
	req := httptest.NewRequest(http.MethodGet, "/api/wails/android-fetch/"+encodedPayload, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var got map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	for key, want := range map[string]string{
		"method": http.MethodPost,
		"query":  "ok",
		"header": "from-bridge",
		"body":   "kept",
	} {
		if got[key] != want {
			t.Fatalf("%s = %q, want %q", key, got[key], want)
		}
	}
}
