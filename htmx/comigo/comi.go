package comigo

import (
	"github.com/yumenaka/comigo/util/logger"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/resource"
	"github.com/yumenaka/comigo/routers"
)

// StartComigoServer 启动Comigo Web服务器
func StartComigoServer(engine *gin.Engine) {
	logger.Info("Start Comigo File Server.")
	config.GetCfg().OpenBrowser = false
	//解析命令，扫描文件
	cmd.StartScan(os.Args)
	routers.BindAPI(engine)
	// Admin界面 TODO：用 Htmx 重写
	resource.EmbedAdmin(engine)
	//设置临时文件夹
	config.AutoSetCachePath()
}
