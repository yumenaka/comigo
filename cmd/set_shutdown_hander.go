package cmd

import (
	"context"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers"
	"github.com/yumenaka/comigo/tools/logger"
)

// SetShutdownHandler 退出时清理临时文件的函数
func SetShutdownHandler() {
	// 优雅地停止或重启： https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-with-context/server.go
	// 创建侦听来自操作系统的中断信号的上下文。
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	defer stop()
	// Listen for the interrupt signal.
	// 监听中断信号。
	<-ctx.Done()
	// 恢复中断信号的默认行为并通知用户关机。
	stop()
	logger.Info(locale.GetString("shutdown_hint"))
	// 清理临时文件
	if config.GetCfg().ClearCacheExit {
		logger.Infof("\r"+locale.GetString("start_clear_file")+" CacheDir:%s ", config.GetCfg().CacheDir)
		allBooks, err := model.IStore.ListBooks()
		if err != nil {
			logger.Infof(locale.GetString("log_error_listing_books"), err)
		}
		for _, book := range allBooks {
			//清理某一本书的缓存
			cachePath := path.Join(config.GetCfg().CacheDir, book.BookID)
			err := os.RemoveAll(cachePath)
			if err != nil {
				logger.Infof(locale.GetString("log_error_clearing_temp_files"), cachePath)
			} else if config.GetCfg().Debug {
				logger.Infof(locale.GetString("log_cleared_temp_files"), cachePath)
			}
		}
		logger.Infof("%s", locale.GetString("clear_temp_file_completed"))
	}
	// 统一走 routers.StopWebServer，确保 SSE、Tailscale 和 HTTP Server 按同一套流程关闭。
	if err := routers.StopWebServer(); err != nil {
		logger.Infof(locale.GetString("err_server_shutdown_failed")+": %v", err)
	}
	logger.Info("Comigo Server exit.")
}
