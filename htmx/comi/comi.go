package comi

import (
	"fmt"
	"os"

	"github.com/yumenaka/comi/cmd"
	"github.com/yumenaka/comi/config"
)

// StartComigoWebserver 启动Comigo Web服务器
func StartComigoWebserver() {
	fmt.Println("Start Comigo Server.")
	config.Config.OpenBrowser = false
	//解析命令，扫描文件
	cmd.StartScan(os.Args)
	//设置临时文件夹
	config.SetTempDir()
}
