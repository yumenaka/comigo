package cmd

import (
	"context"
	"fmt"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/types"
	"log"
	"os/signal"
	"syscall"
	"time"
)

// SetShutdownHandler TODO:退出时清理临时文件的函数
func SetShutdownHandler() {
	//优雅地停止或重启： https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-with-context/server.go
	// 创建侦听来自操作系统的中断信号的上下文。
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	defer stop()
	// Listen for the interrupt signal.
	//监听中断信号。
	<-ctx.Done()
	//恢复中断信号的默认行为并通知用户关机。
	stop()
	log.Println(locale.GetString("ShutdownHint"))
	//清理临时文件
	if config.Config.ClearCacheExit {
		fmt.Println("\r" + locale.GetString("start_clear_file") + " CachePath:" + config.Config.CachePath)
		types.ClearTempFilesALL(config.Config.Debug, config.Config.CachePath)
		fmt.Println(locale.GetString("clear_temp_file_completed"))
	}
	// 上下文用于通知服务器它有 5 秒的时间来完成它当前正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := config.Srv.Shutdown(ctx); err != nil {
		//fmt.Println("Comigo Server forced to shutdown: ", err)
		//time.Sleep(3 * time.Second)
		log.Fatal("Comigo Server forced to shutdown: ", err)
	}
	log.Println("Comigo Server exit.")
}
