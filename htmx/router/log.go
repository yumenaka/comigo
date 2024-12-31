package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/util/logger"
)

func SetGinLogger(engine *gin.Engine) {
	//设置log中间件
	logger.ReportCaller = config.GetDebug()
	//TODO:输出到tui。
	//logger.SetOutput(logger.NewLogBuffer())
	ginLogHandler := logger.GinLogHandler(config.GetLogToFile(), config.GetLogFilePath(), config.GetLogFileName(), config.GetDebug())
	engine.Use(ginLogHandler)
	if config.GetLogToFile() {
		// 关闭 log 打印的字体颜色。输出到文件不需要颜色
		gin.DisableConsoleColor()
	}
}
