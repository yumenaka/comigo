package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/util/logger"
	"io"
)

// SetLogger 设置日志中间件
func setLogger(engine *gin.Engine) {
	//禁止控制台输出
	gin.DefaultWriter = io.Discard
	//设置log中间件 TODO:输出到tui界面。
	engine.Use(logger.HandlerLog(config.Config.LogToFile, config.Config.LogFilePath, config.Config.LogFileName))
	if config.Config.LogToFile {
		// 关闭 log 打印的字体颜色。输出到文件不需要颜色
		gin.DisableConsoleColor()
	}
}
