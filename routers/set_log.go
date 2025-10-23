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
	// 输出到tui（未实现）
	echoLogHandler := logger.EchoLogHandler(config.GetCfg().LogToFile, config.GetCfg().LogFilePath, config.GetCfg().LogFileName, config.GetCfg().Debug)
	e.Use(echoLogHandler)
	// // 设置日志级别
	// if config.GetCfg().Debug {
	// 	e.Logger.SetLevel(log.DEBUG)
	// } else {
	// 	e.Logger.SetLevel(log.INFO)
	// }
	//
	// // 如果需要输出到文件
	// if config.GetLogToFile() {
	// 	// 确保日志目录存在
	// 	logDir := config.GetLogFilePath()
	// 	if err := os.MkdirAll(logDir, 0o755); err != nil {
	// 		logger.Errorf("创建日志目录失败: %v", err)
	// 		return
	// 	}
	//
	// 	// 打开日志文件
	// 	logFile := filepath.Join(logDir, config.GetLogFileName())
	// 	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	// 	if err != nil {
	// 		logger.Errorf("打开日志文件失败: %v", err)
	// 		return
	// 	}
	//
	// 	// 如果是调试模式，同时输出到文件和控制台
	// 	if config.GetCfg().Debug {
	// 		e.Logger.SetOutput(io.MultiWriter(f, os.Stdout))
	// 	} else {
	// 		e.Logger.SetOutput(f)
	// 	}
	// } else {
	// 	// 如果不输出到文件，则输出到标准输出
	// 	e.Logger.SetOutput(os.Stdout)
	// }
	//
	// // 设置日志格式
	// e.Logger.SetHeader("${time_rfc3339} ${level} ${short_file}:${line} ${message}")
}
