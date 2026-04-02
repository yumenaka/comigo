package login

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"golang.org/x/oauth2"
)

const (
	// OAuthStateCookieName OAuth 登录流程的 CSRF 状态 Cookie。
	OAuthStateCookieName = "oauth_state"
)

// buildOAuthConfig 根据当前配置构造 OAuth2 配置。
func buildOAuthConfig(c echo.Context) *oauth2.Config {
	cfg := config.GetCfg()
	scopes := append([]string{}, cfg.OAuthScopes...)
	if len(scopes) == 0 {
		scopes = []string{"openid", "profile", "email"}
	}
	return &oauth2.Config{
		ClientID:     strings.TrimSpace(cfg.OAuthClientID),
		ClientSecret: strings.TrimSpace(cfg.OAuthClientSecret),
		RedirectURL:  oauthRedirectURL(c),
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  strings.TrimSpace(cfg.OAuthAuthURL),
			TokenURL: strings.TrimSpace(cfg.OAuthTokenURL),
		},
	}
}

// oauthRedirectURL 解析 OAuth 回调地址。
func oauthRedirectURL(c echo.Context) string {
	if redirectURL := strings.TrimSpace(config.GetCfg().OAuthRedirectURL); redirectURL != "" {
		return redirectURL
	}
	scheme := c.Scheme()
	host := c.Request().Host
	if forwardedProto := strings.TrimSpace(c.Request().Header.Get("X-Forwarded-Proto")); forwardedProto != "" {
		scheme = strings.Split(forwardedProto, ",")[0]
	}
	if forwardedHost := strings.TrimSpace(c.Request().Header.Get("X-Forwarded-Host")); forwardedHost != "" {
		host = strings.Split(forwardedHost, ",")[0]
	}
	return fmt.Sprintf("%s://%s/api/oauth/callback", scheme, host)
}

// randomState 生成 OAuth state。
func randomState() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func setOAuthStateCookie(c echo.Context, state string) {
	c.SetCookie(&http.Cookie{
		Name:     OAuthStateCookieName,
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   c.Scheme() == "https" || c.Request().TLS != nil,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(10 * time.Minute),
	})
}

func clearOAuthStateCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     OAuthStateCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   c.Scheme() == "https" || c.Request().TLS != nil,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(-1 * time.Hour),
	})
}

func loginRedirectWithError(c echo.Context, errorCode string) error {
	return c.Redirect(http.StatusFound, "/login?error="+url.QueryEscape(errorCode))
}

// StartOAuthLogin 发起 OAuth 登录跳转。
func StartOAuthLogin(c echo.Context) error {
	if !config.GetCfg().LoginProtection || !config.GetCfg().HasOAuthLoginConfigured() {
		return loginRedirectWithError(c, "oauth_not_configured")
	}
	state, err := randomState()
	if err != nil {
		return loginRedirectWithError(c, "oauth_login_failed")
	}
	setOAuthStateCookie(c, state)
	return c.Redirect(http.StatusFound, buildOAuthConfig(c).AuthCodeURL(state))
}

// OAuthCallback 处理 OAuth 回调。
func OAuthCallback(c echo.Context) error {
	if !config.GetCfg().LoginProtection || !config.GetCfg().HasOAuthLoginConfigured() {
		return loginRedirectWithError(c, "oauth_not_configured")
	}
	if c.QueryParam("error") != "" {
		clearOAuthStateCookie(c)
		return loginRedirectWithError(c, "oauth_login_failed")
	}
	state := strings.TrimSpace(c.QueryParam("state"))
	code := strings.TrimSpace(c.QueryParam("code"))
	if state == "" || code == "" {
		clearOAuthStateCookie(c)
		return loginRedirectWithError(c, "oauth_login_failed")
	}

	stateCookie, err := c.Cookie(OAuthStateCookieName)
	if err != nil || stateCookie == nil || stateCookie.Value != state {
		clearOAuthStateCookie(c)
		return loginRedirectWithError(c, "oauth_state_mismatch")
	}
	clearOAuthStateCookie(c)

	token, err := buildOAuthConfig(c).Exchange(c.Request().Context(), code)
	if err != nil {
		return loginRedirectWithError(c, "oauth_login_failed")
	}

	client := buildOAuthConfig(c).Client(c.Request().Context(), token)
	req, err := http.NewRequestWithContext(c.Request().Context(), http.MethodGet, strings.TrimSpace(config.GetCfg().OAuthUserInfoURL), nil)
	if err != nil {
		return loginRedirectWithError(c, "oauth_login_failed")
	}
	resp, err := client.Do(req)
	if err != nil {
		return loginRedirectWithError(c, "oauth_login_failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return loginRedirectWithError(c, "oauth_login_failed")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return loginRedirectWithError(c, "oauth_login_failed")
	}
	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return loginRedirectWithError(c, "oauth_login_failed")
	}
	username := oauthDisplayName(userInfo)
	if username == "" {
		return loginRedirectWithError(c, "oauth_userinfo_invalid")
	}
	if err := issueLoginCookie(c, username); err != nil {
		return loginRedirectWithError(c, "oauth_login_failed")
	}
	return c.Redirect(http.StatusFound, "/")
}

// oauthDisplayName 从 userinfo 响应中提取显示名。
func oauthDisplayName(userInfo map[string]interface{}) string {
	for _, key := range []string{"preferred_username", "email", "name", "login", "sub"} {
		value, ok := userInfo[key]
		if !ok {
			continue
		}
		if str, ok := value.(string); ok && strings.TrimSpace(str) != "" {
			return strings.TrimSpace(str)
		}
	}
	return ""
}
