package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/util/logger"
)

// SetLogger 设置日志中间件
func SetEchoLogger(e *echo.Echo) {
	// 设置log中间件
	logger.ReportCaller = config.GetDebug()
	// TODO:输出到tui
	echoLogHandler := logger.EchoLogHandler(config.GetLogToFile(), config.GetLogFilePath(), config.GetLogFileName(), config.GetDebug())
	e.Use(echoLogHandler)
}
