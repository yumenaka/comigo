package jwt

import (
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 定义登录验证函数并添加到 ginJWTMiddleware 中的 Authenticator 字段
//
//	在这里写登录验证逻辑
func Authenticator(c *gin.Context) (interface{}, error) {
	user := User{}
	if err := c.ShouldBind(&user); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	if user.Username == "admin" && user.Password == "admin" {
		return user, nil
	}
	return nil, jwt.ErrFailedAuthentication
}

// 定义登录后权限验证函数并添加到 ginJWTMiddleware 中的 Authorizator 字段
func Authorizator(data interface{}, ctx *gin.Context) bool {
	return strings.Contains(data.(string), "admin")
}
