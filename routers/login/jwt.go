package login

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools/logger"
)

// JwtCustomClaims 扩展默认的"JWT声明"。更多示例，请参见https://github.com/golang-jwt/jwt
type JwtCustomClaims struct {
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
	jwt.RegisteredClaims
}

const (
	// CookieName JWT令牌的Cookie名称
	CookieName = "jwt_token"
)

// newClaims 创建登录成功后的 JWT 声明。
func newClaims(username string) *JwtCustomClaims {
	return &JwtCustomClaims{
		Username: username,
		Admin:    true, // 当前是单管理员模型，已通过配置页维护管理员账号密码。
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(loginDuration())),
		},
	}
}

func loginDuration() time.Duration {
	return time.Minute * time.Duration(config.GetCfg().Timeout)
}

// signedToken 将声明签名为 JWT 字符串。
func signedToken(claims *JwtCustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetJwtSigningKey()))
}

// setTokenCookie 将 JWT 写入 Cookie。
func setTokenCookie(c echo.Context, token string) {
	cookie := new(http.Cookie)
	cookie.Name = CookieName
	cookie.Value = token
	cookie.Expires = time.Now().Add(loginDuration())
	cookie.Path = "/"
	cookie.HttpOnly = true                                          // JWT 不再暴露给前端 JS，退出登录统一走 /api/logout。
	cookie.Secure = c.Scheme() == "https" || c.Request().TLS != nil // 如果是HTTPS则设置Secure
	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)
}

func Login(c echo.Context) error {
	// 表单用于浏览器登录，JSON 用于签发 REST 客户端的 Bearer 令牌。
	var request struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}
	if err := c.Bind(&request); err != nil {
		return echo.ErrBadRequest
	}
	username, password := request.Username, request.Password
	// 如果未配置密码登录，则不接受表单登录
	if !config.GetCfg().HasPasswordLoginConfigured() {
		return echo.ErrTeapot
	}
	// 如果密码错误，则不生成 JWT
	if username != config.GetCfg().Username || password != config.GetCfg().Password {
		logger.Infof(locale.GetString("log_login_failed"), username)
		return echo.ErrUnauthorized
	}
	token, err := signedToken(newClaims(username))
	if err != nil {
		return err
	}
	if strings.HasPrefix(c.Request().Header.Get(echo.HeaderContentType), echo.MIMEApplicationJSON) {
		return c.JSON(http.StatusOK, echo.Map{"token": token, "token_type": "Bearer", "expires_in": int(loginDuration().Seconds())})
	}
	setTokenCookie(c, token)

	// 返回登录成功信息，不再返回token本身
	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
		"message": "登录成功",
	})
}

// Logout 登出，删除Cookie中的JWT
func Logout(c echo.Context) error {
	// 清除Cookie
	cookie := new(http.Cookie)
	cookie.Name = CookieName
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-1 * time.Hour) // 设置为过期
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = c.Scheme() == "https" || c.Request().TLS != nil
	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "登出成功",
	})
}
