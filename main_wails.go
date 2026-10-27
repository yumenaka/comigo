//go:build wails && !js && !bindings

//go:generate go tool goversioninfo -icon=build/windows/icon.ico -manifest=goversioninfo.exe.manifest

package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers"
	"github.com/yumenaka/comigo/tools/wails_systray"
)

// main 保持 desktop 配置与管理命令行为，窗口和托盘共用 Wails 原生事件循环。
func main() {
	config.UseDesktopConfigProfile()
	if handled, err := cmd.RunProcessCommand(os.Args[1:], os.Stdout); handled {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	for _, arg := range os.Args[1:] {
		if arg == "-v" || arg == "--version" || arg == "-h" || arg == "--help" {
			cmd.Execute()
			return
		}
	}
	cmd.Execute()
	if err := ensureWailsAndroidImportStore(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cmd.AddStoreUrls(cmd.Args)
	if err := cmd.LoadMetadata(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := routers.StartWebServer(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cmd.Shutdown()
	// 资源始终走内嵌 Echo，不能被开发服务器环境变量替换。
	_ = os.Unsetenv("FRONTEND_DEVSERVER_URL")
	panel := &wails_systray.PanelService{}
	app := application.New(application.Options{
		Name: "Comigo", Description: "Comigo desktop reader",
		Services:   []application.Service{application.NewService(NewApp()), application.NewService(panel)},
		Assets:     application.AssetOptions{Handler: http.HandlerFunc(serveWailsAsset), DisableLogging: true},
		OnShutdown: cmd.Shutdown,
	})
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name: "main", Title: wailsWindowTitle(), Width: 1024, Height: 768, URL: config.PrefixPath("/"),
		BackgroundColour: application.NewRGB(27, 38, 54), DefaultContextMenuDisabled: true,
	})
	routers.SetWailsRuntime(app, window)
	tray := wails_systray.New(app, window)
	wails_systray.AttachPanel(tray, panel)
	var scanOnce sync.Once
	window.OnWindowEvent(events.Common.WindowRuntimeReady, func(*application.WindowEvent) { scanOnce.Do(func() { go finishWailsStartupScan(window) }) })
	cmd.LoadUserPlugins()
	config.StartOrStopAutoRescan()
	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// wailsWindowTitle 保留版本号，便于确认当前运行的构建。
func wailsWindowTitle() string { return "Comigo " + config.GetVersion() }

// finishWailsStartupScan 后台扫描，主窗口无需等待整个书库读完。
func finishWailsStartupScan(window application.Window) {
	cmd.ScanStore()
	cmd.SaveMetadata()
	window.Reload()
}

// serveWailsAsset 由 Wails 注入运行时，业务页面交给现有 Echo 路由。
func serveWailsAsset(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/panel.html" {
		wails_systray.ServePanel(w, r)
		return
	}
	if config.Server == nil || config.Server.Handler == nil {
		http.NotFound(w, r)
		return
	}
	config.Server.Handler.ServeHTTP(w, r)
}
