package cmd

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/locale"
	"github.com/yumenaka/comigo/util/logger"
)

// SetShutdownHandler TODO:退出时清理临时文件的函数
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
	log.Println(locale.GetString("shutdown_hint"))
	// 清理临时文件
	if config.GetClearCacheExit() {
		logger.Infof("\r"+locale.GetString("start_clear_file")+" CachePath:%s ", config.GetCachePath())
		model.ClearTempFilesALL(config.GetDebug(), config.GetCachePath())
		logger.Infof("%s", locale.GetString("clear_temp_file_completed"))
	}
	// 上下文用于通知服务器它有 5 秒的时间来完成它当前正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 只能通过http.Server.Shutdown()/http.Server.Close()等http包里的方法去实现,没办法自己实现.
	// 因为这样的设计即使你给自定义Server接口的实现类设计了Shutdown()方法,也调用不到.
	// 本质上还是因为从端口启动开始,后续的所有工作都是http包来完成的,我们无法干涉这其中的步骤
	if err := config.Srv.Shutdown(ctx); err != nil {
		// logger.Infof("Comigo Server forced to shutdown: ", err)
		// time.Sleep(3 * time.Second)
		log.Fatal("Comigo Server forced to shutdown: ", err)
	}
	log.Println("Comigo Server exit.")
}
