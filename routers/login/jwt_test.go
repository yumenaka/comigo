package login

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"net/http/httptest"
	"strings"
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

// JSON 登录只签发 Bearer，避免客户端 Cookie 跨端口串用或退出后继续生效。
func TestJSONLoginIssuesBearerToken(t *testing.T) {
	old := config.CopyCfg()
	defer func() { *config.GetCfg() = old }()
	config.GetCfg().Username = "reader"
	config.GetCfg().Password = "test-password"
	config.GetCfg().Timeout = 60
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(`{"username":"reader","password":"test-password"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	if err := Login(e.NewContext(req, rec)); err != nil {
		t.Fatal(err)
	}
	var body struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	token, err := jwt.Parse(body.Token, func(token *jwt.Token) (interface{}, error) { return []byte(config.GetJwtSigningKey()), nil }, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || !token.Valid {
		t.Fatalf("invalid login token: %v", err)
	}
	cookies := rec.Result().Cookies()
	if len(cookies) != 0 {
		t.Fatal("REST login must not create a browser cookie")
	}
}
