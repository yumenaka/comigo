package main

import (
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers"
)

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// 解析命令行参数
	cmd.InitFlags()
	// 初始化配置文件
	cobra.OnInitialize(cmd.LoadConfigFile)
	// 通过“可执行文件名”设置部分默认参数,目前不生效
	config.SetByExecutableFilename()
	// 设置临时文件夹
	config.AutoSetCachePath()
	// 扫描命令行指定的书库与文件
	cmd.ScanStore(os.Args)

	// 在命令行显示QRCode
	cmd.ShowQRCode()

	// 获取网页服务器（echo）
	echo := routers.GetWebServer()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Comigo",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:  nil,
			Handler: http.HandlerFunc(echo.ServeHTTP),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
