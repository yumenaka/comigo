package main

import (
	"flag"
	"fmt"
	"github.com/yumenaka/comigo/application/ebitenui/pages/scroll_mode"
	"github.com/yumenaka/comigo/application/ebitenui/resources"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yumenaka/comigo/application/ebitenui/comi"
	"github.com/yumenaka/comigo/application/ebitenui/components"
	"github.com/yumenaka/comigo/application/ebitenui/model"
	"github.com/yumenaka/comigo/application/ebitenui/pages/book_shelf"
	"github.com/yumenaka/comigo/config"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
)

// Game object used by ebiten
type game struct {
	ui *ebitenui.UI
}

func main() {
	//如果参数当中有--debug，则不启动UI
	//debugMode := flag.Bool("debug", false, "Disable UI by debug mode.")
	// 解析命令行参数
	flag.Parse()
	//if *debugMode {
	//	fmt.Println("Debug Mode, no thing to do.")
	//	select {}
	//	return
	//}
	// 启动Comigo Web服务器
	// 如果参数当中有--debug，则不启动UI，且什么都不做（UI抢焦点，所以用另一台机器rsync同步代码跑效果）
	comi.StartComigoWebserver()

	// 创建一个新的 ReaderConfig 对象，用于配置阅读器的设置。
	readerConfig := model.NewReaderConfig()
	config.Cfg.OpenBrowser = false
	config.Cfg.UseCache = true
	config.Cfg.ClearCacheExit = false
	readerConfig.SetTitle("Comigo "+config.GetVersion()).
		// 阅读器模式。
		SetReaderMode(model.ScrollMode).
		// 窗口是否全屏。
		SetWindowFullScreen(false).
		// 窗口是否有边框和标题栏。
		SetWindowDecorated(true).
		// 窗口是否可以调整大小。
		SetWindowResizingModeEnabled(ebiten.WindowResizingModeEnabled).
		// 窗口的宽度和高度。
		SetWindowSize(1280, 800).
		// 运行选项。
		SetRunOptions(ebiten.RunGameOptions{
			// 背景透明
			ScreenTransparent: false,
		})
	// 设置窗口大小
	ebiten.SetWindowSize(readerConfig.Width, readerConfig.Height)
	// 设置窗口标题。
	ebiten.SetWindowTitle(readerConfig.Title)
	// 设置窗口是否可以调整大小。
	ebiten.SetWindowResizingMode(readerConfig.WindowResizingModeEnabled)
	// 设置窗口是否有边框和标题栏。
	ebiten.SetWindowDecorated(readerConfig.WindowDecorated)
	//  设置窗口是否全屏。
	ebiten.SetFullscreen(readerConfig.WindowFullScreen)
	// SetScreenClearedEveryFrame 用于启用或禁用每个帧开始时清除屏幕的功能。
	// 默认值为 true，这意味着默认情况下每个帧都会清屏。
	ebiten.SetScreenClearedEveryFrame(false)

	// 构建UI
	ui, closeUI, err := createUI(readerConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer closeUI()
	// 创建游戏对象
	g := game{
		ui: ui,
	}
	// 运行ebitenUI
	err = ebiten.RunGameWithOptions(
		&g,
		&readerConfig.RunOptions,
	)
	// 检查错误
	if err != nil {
		log.Print(err)
	}
}

// 构建UI，需要进一步拆分
func createUI(readerConfig *model.ReaderConfig) (*ebitenui.UI, func(), error) {
	// 创建一个新的根容器，用于包含整个 UI。
	rootContainer := book_shelf.NewPage()
	//这会将根容器添加到 UI，以便将其展示。
	ui := &ebitenui.UI{
		Container: rootContainer,
	}
	// 加载资源
	res, err := resources.NewUIResources()
	if err != nil {
		return nil, nil, err
	}
	rootContainer.AddChild(components.HeaderContainer(res, readerConfig, ui))
	rootContainer.AddChild(scroll_mode.BodyContainer())
	rootContainer.AddChild(components.FooterContainer())

	return ui, func() {
		// 在结束时，关闭资源。
		//res.close()
	}, nil
}

func (g *game) Update() error {
	// ui.Update() 必须在 ebiten Update 函数中调用，以处理用户输入和其他事情
	g.ui.Update()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	// 应在 ebiten Draw 函数中调用 ui.Draw()，以将 UI 绘制到屏幕上。
	// 还应该在游戏的所有其他事件渲染之后调用它，以便它显示在游戏世界的顶部。
	g.ui.Draw(screen)
	// 在屏幕左上角显示当前 FPS。
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.ActualFPS()))
}

func (g *game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
