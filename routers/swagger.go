//go:build doc
// +build doc

// 在需要的地方声明swagHandler，该参数不为空时才加入路由，以减少包体积
// 通过go build -tags "doc"来打包带文档的包，直接go build来打包不带文档的包

//文档：https://github.com/swaggo/gin-swagger

//go get -u github.com/swaggo/gin-swagger
//go get -u github.com/swaggo/files

// https://www.cnblogs.com/Ivan-Wu/p/15821288.html
// 注意：每次修改swagger注释或其他参数的时候需要（swag init）重新生成swagger文件才会生效

//访问：http://127.0.0.1:1234/swagger/index.html

package routers

import (
	docs "github.com/yumenaka/comi/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	docs.SwaggerInfo.BasePath = "/api"
	swagHandler = ginSwagger.WrapHandler(swaggerFiles.Handler)
}
