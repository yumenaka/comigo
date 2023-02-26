package token

import (
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// User 定义 User 结构体，用于接受登录的用户名与密码
type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Authenticator 认证器：登录验证函数,在这里写登录验证逻辑
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

// Authorizator 登录后权限验证函数 //当用户通过token请求受限接口时，会经过这段逻辑
func Authorizator(data interface{}, _ *gin.Context) bool {
	return strings.Contains(data.(string), "admin")
}
