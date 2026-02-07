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

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	// 检查是否需要登录
	if config.GetCfg().Username == "" || config.GetCfg().Password == "" {
		logger.Infof(locale.GetString("log_cannot_set_username_or_password") + "\n")
		return echo.ErrTeapot
	}
	// 如果未设置密码或密码错误，则不生成 JWT
	if username != config.GetCfg().Username || password != config.GetCfg().Password {
		logger.Infof(locale.GetString("log_login_failed")+"\n", username, config.GetCfg().Username, config.GetCfg().Password, password)
		return echo.ErrUnauthorized
	}

	// 设置自定义"JWT声明"
	claims := &JwtCustomClaims{
		username,
		true, // 账号管理（未实现）
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(config.GetCfg().Timeout))),
		},
	}

	// 用"JWT声明"创建令牌
	// HS256 是对称的，这要求解码器也具有秘密私钥
	// 可以换成非对称的jwt.SigningMethodRS512，需要两个密钥：一个公钥和一个必须保密的私钥 https://packagemain.tech/p/json-web-tokens-in-go
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
	cookie.HttpOnly = false                // 防止JavaScript访问的话需要设置为true
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
