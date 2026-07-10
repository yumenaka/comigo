//go:build wails && !js && !bindings

//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=build/windows/icon.ico -manifest=goversioninfo.exe.manifest

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers"
	"github.com/yumenaka/comigo/tools/wails_systray"
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

	config.UseDesktopConfigProfile()
	cmd.Execute()
	tray := wails_systray.Start()
	defer tray.Stop()

	app := NewApp()
	err := wails.Run(&options.App{
		Title:  wailsWindowTitle(),
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Handler: http.HandlerFunc(serveWailsAsset),
		},
		BackgroundColour:         &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		EnableDefaultContextMenu: false,
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			routers.SetWailsContext(ctx)
			tray.SetContext(ctx)
			if err := startComigoForWails(ctx); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				wailsruntime.Quit(ctx)
			}
		},
		OnShutdown: func(context.Context) {
			if err := routers.StopWebServer(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
			}
			tray.Stop()
		},
		OnBeforeClose: func(ctx context.Context) bool {
			return tray.HandleBeforeClose(ctx)
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

// wailsWindowTitle 生成桌面窗口标题，保留完整版本便于用户确认当前构建。
func wailsWindowTitle() string {
	return "Comigo " + config.GetVersion()
}

// startComigoForWails 启动桌面壳内嵌的 Comigo Web 服务。
func startComigoForWails(ctx context.Context) error {
	if err := routers.StartWebServer(); err != nil {
		return err
	}
	routers.StartTailscale()
	cmd.LoadUserPlugins()
	cmd.AddStoreUrls(cmd.Args)
	cmd.LoadMetadata()
	go finishWailsStartupScan(ctx)
	config.StartOrStopAutoRescan()
	return nil
}

// finishWailsStartupScan 后台刷新书库，避免 Wails 首页等扫描完成才出现。
func finishWailsStartupScan(ctx context.Context) {
	cmd.ScanStore()
	cmd.SaveMetadata()
	wailsruntime.WindowReload(ctx)
}

// serveWailsAsset 在 Web 服务就绪后把 Wails 的资源请求转给 Echo。
func serveWailsAsset(w http.ResponseWriter, r *http.Request) {
	if config.Server == nil || config.Server.Handler == nil {
		http.NotFound(w, r)
		return
	}
	config.Server.Handler.ServeHTTP(w, r)
}
