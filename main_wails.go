//go:build wails && !js && !bindings

//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers"
)

// main 是 Wails 桌面壳入口；普通 CLI 入口保留在 main.go，减少合并冲突。
func main() {
	for _, arg := range os.Args {
		if arg == "-v" || arg == "--version" || arg == "-h" || arg == "--help" {
			// 仅打印版本或帮助信息时，不启动 WebView。
			cmd.Execute()
			return
		}
	}

	app := NewApp()
	err := wails.Run(&options.App{
		Title:  "Comigo",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Handler: http.HandlerFunc(serveWailsAsset),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			if err := startComigoForWails(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// startComigoForWails 启动桌面壳内嵌的 Comigo Web 服务。
func startComigoForWails() error {
	cmd.Execute()
	if err := routers.StartWebServer(); err != nil {
		return err
	}
	routers.StartTailscale()
	cmd.LoadUserPlugins()
	cmd.AddStoreUrls(cmd.Args)
	cmd.SetCwdAsScanPathIfNeed()
	cmd.LoadMetadata()
	cmd.ScanStore()
	model.GenerateBookGroup()
	cmd.SaveMetadata()
	config.StartOrStopAutoRescan()
	return nil
}

// serveWailsAsset 在 Web 服务就绪后把 Wails 的资源请求转给 Echo。
func serveWailsAsset(w http.ResponseWriter, r *http.Request) {
	if config.Server == nil || config.Server.Handler == nil {
		http.NotFound(w, r)
		return
	}
	config.Server.Handler.ServeHTTP(w, r)
}
