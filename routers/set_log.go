package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools/logger"
)

// SetLogger 设置日志中间件
func SetEchoLogger(e *echo.Echo) {
	// 设置log中间件
	logger.ReportCaller = config.GetCfg().Debug
	// TUI 会通过 logger.SetMirrorOutput 接收同一份 Echo 访问日志。
	echoLogHandler := logger.EchoLogHandler(config.GetCfg().LogToFile, config.GetCfg().LogFilePath, config.GetCfg().LogFileName, config.GetCfg().Debug)
	e.Use(echoLogHandler)
}
