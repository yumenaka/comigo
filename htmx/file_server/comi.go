package file_server

import (
	"fmt"
	"os"

	"github.com/yumenaka/comi/cmd"
	"github.com/yumenaka/comi/config"
)

// StartComigoServer 启动Comigo Web服务器
func StartComigoServer() {
	fmt.Println("Start Comigo File Server.")
	config.Config.OpenBrowser = false
	//解析命令，扫描文件
	cmd.StartScan(os.Args)
	//设置临时文件夹
	config.SetTempDir()
}
