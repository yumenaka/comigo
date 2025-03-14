package comigo

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/resource"
	"github.com/yumenaka/comigo/routers"
	"github.com/yumenaka/comigo/util/logger"
)

// SetComigoServer 设置Comigo Web服务器
func SetComigoServer(e *echo.Echo) {
	logger.Info("Start Comigo File Server.")
	config.GetCfg().OpenBrowser = true
	// 解析命令，扫描文件
	cmd.SetStore(os.Args)
	routers.BindAPI(e)
	// Admin界面
	resource.EmbedAdmin(e)
	// 设置临时文件夹
	config.AutoSetCachePath()
}
