package cmd

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/routers"
	"log"
	"os"
	"time"
)

//Go 每日一库之 fsnotify  https://darjun.github.io/2020/01/19/godailylib/fsnotify/

//优雅地重启或停止  https://learnku.com/docs/gin-gonic/1.7/examples-graceful-restart-or-stop/11376
func configReloadHandler(e fsnotify.Event) {
	//打印配置文件路径与触发事件
	fmt.Printf("配置文件改变，Comigo将在5秒后重启:%s Op:%s\n", e.Name, e.Op)

	//重新读取改变后的配置文件
	if err := viperInstance.ReadInConfig(); err != nil {
		if common.ConfigFile == "" && common.Config.Debug {
			fmt.Println(err)
		}
	}
	// 把设定文件的内容，解析到构造体里面。
	if err := viperInstance.Unmarshal(&common.Config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// 上下文用于通知服务器它有 5 秒的时间来完成它当前正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := common.Srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	<-ctx.Done()
	//重启 web 服务器
	routers.StartWebServer()
}
