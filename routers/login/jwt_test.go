package login

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
)

// TestTokenAndCookieUseConfiguredTimeout 确认 JWT 与 Cookie 不会使用两套过期时间。
func TestTokenAndCookieUseConfiguredTimeout(t *testing.T) {
	cfg := config.GetCfg()
	oldTimeout := cfg.Timeout
	cfg.Timeout = 90
	defer func() { cfg.Timeout = oldTimeout }()

	before := time.Now().Add(89 * time.Minute)
	claims := newClaims("reader")
	if claims.ExpiresAt.Time.Before(before) || claims.ExpiresAt.Time.After(time.Now().Add(91*time.Minute)) {
		t.Fatalf("unexpected JWT expiry: %v", claims.ExpiresAt.Time)
	}

	e := echo.New()
	rec := httptest.NewRecorder()
	setTokenCookie(e.NewContext(httptest.NewRequest("POST", "/api/login", nil), rec), "token")
	cookies := rec.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Expires.Before(before) || cookies[0].Expires.After(time.Now().Add(91*time.Minute)) {
		t.Fatalf("unexpected cookie expiry: %#v", cookies)
	}
}
