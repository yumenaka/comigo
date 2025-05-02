package login

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
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

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// 如果未设置密码或密码错误，则不生成 JWT
	if config.GetPassword() == "" || (username != config.GetUsername() && password != config.GetPassword()) {
		return echo.ErrUnauthorized
	}

	// 设置自定义"JWT声明"
	claims := &JwtCustomClaims{
		username,
		true, // TODO：账号管理
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(config.GetTimeout()))),
		},
	}

	// 用"JWT声明"创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 生成编码的jwt token
	t, err := token.SignedString([]byte(config.GetJwtSigningKey()))
	if err != nil {
		return err
	}

	// 将token设置到Cookie
	cookie := new(http.Cookie)
	cookie.Name = CookieName
	cookie.Value = t
	cookie.Expires = time.Now().Add(24 * time.Hour * 30) // 30天有效期
	cookie.Path = "/"
	cookie.HttpOnly = true                 // 防止JavaScript访问
	cookie.Secure = c.Request().TLS != nil // 如果是HTTPS则设置Secure
	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)

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
