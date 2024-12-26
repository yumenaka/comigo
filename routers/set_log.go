package routers

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/util/logger"
)

// SetLogger 设置日志中间件
func setLogger(engine *gin.Engine) {
	//禁止控制台输出
	gin.DefaultWriter = io.Discard
	//设置log中间件 TODO:输出到tui界面。
	engine.Use(logger.HandlerLog(config.Cfg.LogToFile, config.Cfg.LogFilePath, config.Cfg.LogFileName))
	if config.Cfg.LogToFile {
		// 关闭 log 打印的字体颜色。输出到文件不需要颜色
		gin.DisableConsoleColor()
	}
}
