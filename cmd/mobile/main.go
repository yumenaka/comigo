package server

import (
	"os"
	"strconv"

	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers"
	// _ "golang.org/x/mobile/bind"
)

func Start(path string) (string, error) {
	// 初始化命令行flag，环境变量与配置文件
	cmd.Execute()
	// 扫描书库（命令行指定）
	cmd.ScanStore(os.Args)
	// 启动网页服务器（不阻塞）
	routers.StartWebServer()
	return strconv.Itoa(config.GetCfg().Port), nil
}
