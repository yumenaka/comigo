package token

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/util/logger"
)

// User 定义 User 结构体，用于接受登录的用户名与密码
type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// JWTCustomClaims  自定义 JWT Claims
type JWTCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NewJwtMiddleware 返回一个新的 JWT 中间件
func NewJwtMiddleware() (echo.MiddlewareFunc, error) {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 跳过登录路由的验证
			if c.Path() == "/api/login" {
				return next(c)
			}

			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return echo.ErrUnauthorized
			}

			// 移除 "Bearer " 前缀
			tokenString := strings.TrimPrefix(auth, "Bearer ")

			// 解析 token
			token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.GetUsername() + config.GetPassword()), nil
			})

			if err != nil || !token.Valid {
				return echo.ErrUnauthorized
			}

			// 将用户信息存储到上下文中
			if claims, ok := token.Claims.(*JWTCustomClaims); ok {
				c.Set("username", claims.Username)
			}

			return next(c)
		}
	}, nil
}

// Login 处理登录请求
func Login(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return echo.ErrBadRequest
	}

	logger.Infof("username is %s, password is %s, cfg is %s@%s",
		user.Username, user.Password, config.GetUsername(), config.GetPassword())

	// 如果未设置密码，直接通过
	if config.GetPassword() == "" {
		return generateToken(c, user)
	}

	// 验证用户名和密码
	if user.Username == config.GetUsername() && user.Password == config.GetPassword() {
		return generateToken(c, user)
	}

	// 设置跨域 cookie
	c.SetCookie(&http.Cookie{
		Name:     "SameSite",
		Value:    "None",
		MaxAge:   3600,
		Path:     "/",
		Domain:   c.Request().Host,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	return echo.ErrUnauthorized
}

// generateToken 生成 JWT token
func generateToken(c echo.Context, user *User) error {
	// 设置自定义 claims
	claims := &JWTCustomClaims{
		user.Username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.GetTimeout()) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 生成签名字符串
	t, err := token.SignedString([]byte(config.GetUsername() + config.GetPassword()))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":   t,
		"expires": claims.ExpiresAt,
	})
}

// GetUserFromContext 从上下文中获取用户名
func GetUserFromContext(c echo.Context) string {
	username := c.Get("username")
	if username == nil {
		return ""
	}
	return username.(string)
}
