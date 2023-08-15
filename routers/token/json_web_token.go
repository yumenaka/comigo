package token

import (
	"fmt"
	"github.com/yumenaka/comi/common"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// NewJwtMiddleware returns a new JWT middleware.
// sample: https://github.com/appleboy/gin-jwt/blob/master/_example/basic/server.go
func NewJwtMiddleware() (*jwt.GinJWTMiddleware, error) {
	jwtMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "comigo server",                                         //标识
		SigningAlgorithm: "HS256",                                                 //加密算法
		Key:              []byte(common.Config.Username + common.Config.Password), //JWT服务端密钥，需要确保别人不知道
		//time.Duration类型 不能直接和 int类型相乘，需要先将变量转换为time.Duration类型
		Timeout:       time.Minute * time.Duration(common.Config.Timeout), //jwt过期时间
		MaxRefresh:    time.Minute * time.Duration(common.Config.Timeout), //刷新时，最大能延长多少时间
		IdentityKey:   "id",                                               //指定cookie的id
		Authenticator: Authenticator,                                      //认证器：根据登录信息进行用户认证。须返回用户数据作为用户标识符，它将被存储在Claim Array中。// 必须
		//Authorizator:     Authorizator,                                            //授权者： 应执行已验证用户授权的回调函数。	// 可选，默认为成功。
		SendCookie: true, //是否发送cookie
		//验证失败后的函数调用，可用于自定义返回的 JSON 格式之类
		Unauthorized: func(c *gin.Context, code int, message string) {
			fmt.Println(code)
			fmt.Println(message)
			c.JSON(code, gin.H{
				"code Unauthorized":    code,
				"message Unauthorized": message,
			})
		},
		//定义登录成功后用户名储存以及传递用户名到 Authorization
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return claims["username"]
		},
		//添加额外业务相关的信息
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if user, ok := data.(User); ok {
				return jwt.MapClaims{"username": user.Username}
			}
			return jwt.MapClaims{}
		},
		// 指定从哪里获取token 其格式为："<source>:<name>" 如有多个，用逗号隔开，可用值：
		//header: Authorization, query: token, cookie: jwt
		TokenLookup: "header: Authorization,query: token, cookie: jwt",
		//Header 中 token 的头部字段，默认值 Bearer
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		return nil, err
	}
	err = jwtMiddleware.MiddlewareInit()
	if err != nil {
		return nil, err
	}
	return jwtMiddleware, nil
}

// User 定义 User 结构体，用于接受登录的用户名与密码  //jwt中payload的数据结构
type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Authenticator 认证器：login接口会经过这段逻辑
func Authenticator(c *gin.Context) (interface{}, error) {
	// TODO : データベースやストレージ、SaaSからuserIDを元にユーザー情報を取得する
	//解析Body数据，格式是JSON 或 XML。根据 "Content-Type" header 来判断。
	user := User{}
	if err := c.ShouldBind(&user); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	fmt.Printf("username is %s, password is %s,Config is %s@%s\",", user.Username, user.Password, common.Config.Username, common.Config.Password)
	if "" == common.Config.Username || "" == common.Config.Password {
		return user, nil
	}
	//登录验证函数,打印用户信息和错误信息
	if user.Username == common.Config.Username && user.Password == common.Config.Password {
		return user, nil
	}
	//解决跨域问题
	c.SetCookie("SameSite", "None", 3600, "/", c.Request.Host, true, true)
	return nil, jwt.ErrFailedAuthentication
}

func Authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*User); ok {
		fmt.Println(v)
		return true
	}
	return false
}

func userIdInJwt(c *gin.Context) string {
	claims := jwt.ExtractClaims(c)
	userID := claims[jwt.IdentityKey]
	return userID.(string)
}
