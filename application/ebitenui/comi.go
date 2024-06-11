package main

import (
	"fmt"
	"github.com/yumenaka/comi/cmd"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/routers"
	"os"
)

// StartComigoWebserver 启动Comigo Web服务器
func StartComigoWebserver() {
	fmt.Println("UI is enabled.")
	config.Config.OpenBrowser = false
	//解析命令，扫描文件
	cmd.StartScan(os.Args)
	//设置临时文件夹
	config.SetTempDir()
	//SetWebServerPort
	routers.SetWebServerPort()
	//设置书籍API
	routers.StartWebServer()
}
