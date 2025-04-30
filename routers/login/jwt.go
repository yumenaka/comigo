package login

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
)

// JwtCustomClaims 扩展默认的“JWT声明”。更多示例，请参见https://github.com/golang-jwt/jwt
type JwtCustomClaims struct {
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
	jwt.RegisteredClaims
}

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// 如果未设置密码或密码错误，则不生成 JWT
	if config.GetPassword() == "" && username != config.GetUsername() && password != config.GetPassword() {
		return echo.ErrUnauthorized
	}

	// 设置自定义“JWT声明”
	claims := &JwtCustomClaims{
		username,
		true, // TODO：账号管理
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		},
	}

	// 用“JWT声明”创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 生成编码的jwt token并将其作为响应发送。加密用密钥是用户名和密码的组合
	t, err := token.SignedString([]byte(config.GetJwtSigningKey()))
	if err != nil {
		return err
	}
	// 返回 token
	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func Accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func Restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	username := claims.Username
	return c.String(http.StatusOK, "Welcome "+username+"!")
}
