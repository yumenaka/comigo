package comigo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"

	"github.com/yumenaka/comi/cmd"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/routers"
)

// StartComigoServer 启动Comigo Web服务器
func StartComigoServer(engine *gin.Engine) {
	fmt.Println("Start Comigo File Server.")
	config.Config.OpenBrowser = false
	//解析命令，扫描文件
	cmd.StartScan(os.Args)
	routers.BindAPI(engine)
	//设置临时文件夹
	config.SetTempDir()
}
