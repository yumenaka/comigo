package login

import (
	"net/http"
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
		Admin:    true, // 账号管理（未实现）
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(config.GetCfg().Timeout))),
		},
	}
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
	cookie.Expires = time.Now().Add(24 * time.Hour * 30) // 30天有效期
	cookie.Path = "/"
	cookie.HttpOnly = false                                         // 防止JavaScript访问的话需要设置为true
	cookie.Secure = c.Scheme() == "https" || c.Request().TLS != nil // 如果是HTTPS则设置Secure
	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)
}

// issueLoginCookie 为指定用户签发 JWT 并写入 Cookie。
func issueLoginCookie(c echo.Context, username string) error {
	token, err := signedToken(newClaims(username))
	if err != nil {
		return err
	}
	setTokenCookie(c, token)
	return nil
}

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	// 如果未配置密码登录，则不接受表单登录
	if !config.GetCfg().HasPasswordLoginConfigured() {
		return echo.ErrTeapot
	}
	// 如果密码错误，则不生成 JWT
	if username != config.GetCfg().Username || password != config.GetCfg().Password {
		logger.Infof(locale.GetString("log_login_failed")+"\n", username, config.GetCfg().Username, config.GetCfg().Password, password)
		return echo.ErrUnauthorized
	}
	if err := issueLoginCookie(c, username); err != nil {
		return err
	}

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
	cookie.Secure = c.Request().TLS != nil
	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "登出成功",
	})
}
