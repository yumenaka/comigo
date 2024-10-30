package comigo

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/resource"
	"github.com/yumenaka/comigo/routers"
)

// StartComigoServer 启动Comigo Web服务器
func StartComigoServer(engine *gin.Engine) {
	fmt.Println("Start Comigo File Server.")
	config.Config.OpenBrowser = false
	//解析命令，扫描文件
	cmd.StartScan(os.Args)
	routers.BindAPI(engine)
	// Admin界面 TODO：用 Htmx 重写
	resource.EmbedAdmin(engine)
	//设置临时文件夹
	config.SetTempDir()
}
