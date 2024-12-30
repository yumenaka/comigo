package routers

import (
	"github.com/yumenaka/comigo/util/logger"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/config"
)

// SetLogger 设置日志中间件
func SetLogger(engine *gin.Engine) {
	//禁止控制台输出
	gin.DefaultWriter = io.Discard
	//设置log中间件
	logger.ReportCaller = config.GetDebug()
	//TODO:输出到tui。
	ginLogHandler := logger.GinLogHandler(config.GetLogToFile(), config.GetLogFilePath(), config.GetLogFileName(), config.GetDebug())
	engine.Use(ginLogHandler)
	if config.GetLogToFile() {
		// 关闭 log 打印的字体颜色。输出到文件不需要颜色
		gin.DisableConsoleColor()
	}
}
