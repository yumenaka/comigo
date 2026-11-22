package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools/logger"
)

// SetLogger 设置日志中间件
func SetEchoLogger(e *echo.Echo) {
	cfg := config.GetCfg()
	// TUI 会通过 logger.SetMirrorOutput 接收同一份 Echo 访问日志。
	echoLogHandler := logger.EchoLogHandler(cfg.LogToFile, cfg.LogFilePath, cfg.LogFileName, cfg.Debug)
	e.Use(echoLogHandler)
}
