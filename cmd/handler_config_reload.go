package cmd

import (
	"github.com/fsnotify/fsnotify"
	"github.com/yumenaka/comigo/util/logger"
)

//Go 每日一库之 fsnotify  https://darjun.github.io/2020/01/19/godailylib/fsnotify/

// 监听配置文件变更
// 优雅地重启或停止  https://learnku.com/docs/gin-gonic/1.7/examples-graceful-restart-or-stop/11376
func handlerConfigReload(e fsnotify.Event) {
	logger.Infof("检测到Comigo配置文件改变:%s Op:%s", e.Name, e.Op)

	//打印配置文件路径与触发事件
	//logger.Infof("配置文件改变，Comigo重启:%s Op:%s\n", e.Name, e.Op)
	////重新读取改变后的配置文件
	//if err := runtimeViper.ReadInConfig(); err != nil {
	//	if config.Config.ConfigPath == "" && config.Config.Debug {
	//		logger.Infof("%s", err)
	//	}
	//}
	//// 把设定文件的内容，解析到构造体里面。
	//if err := runtimeViper.Unmarshal(&config.Config); err != nil {
	//	logger.Infof("%s", err)
	//	time.Sleep(3 * time.Second)
	//	os.Exit(1)
	//}
	////3、扫描配置文件指定的书籍库
	//option := scan.NewScanOption(
	//	true,
	//	config.Config.LocalStores,
	//	config.Config.MaxScanDepth,
	//	config.Config.MinImageNum,
	//	config.Config.TimeoutLimitForScan,
	//	config.Config.ExcludePath,
	//	config.Config.SupportMediaType,
	//	config.Config.SupportFileType,
	//	config.Config.ZipFileTextEncoding,
	//	config.Config.EnableDatabase,
	//	config.Config.ClearDatabaseWhenExit,
	//	config.Config.Debug,
	//)
	//err := scan.InitStore(option)
	//if err != nil {
	//	logger.Infof("Failed to scan store path: %v", err)
	//}
	////4，保存扫描结果到数据库
	//if config.Config.EnableDatabase {
	//	err := scan.SaveResultsToDatabase(viper.ConfigFileUsed(), config.Config.ClearDatabaseWhenExit)
	//	if err != nil {
	//		return
	//	}
	//}
	////5、通过“可执行文件名”设置部分默认参数,目前不生效
	//config.Config.SetByExecutableFilename()
	////重新设置文件下载链接
	//routers.SetDownloadLink()
	//
	//// 上下文用于通知服务器它有 5 秒的时间来完成它当前正在处理的请求
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//if err := config.Srv.Shutdown(ctx); err != nil {
	//	time.Sleep(5 * time.Second)
	//	log.Fatal("Server forced to shutdown: ", err)
	//}
	////重启 web 服务器
	//routers.StartWebServer()
	//<-ctx.Done()
}
