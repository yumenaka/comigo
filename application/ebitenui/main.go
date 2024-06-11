package main

import (
	"flag"
	"fmt"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yumenaka/comi/application/ebitenui/components"
	"github.com/yumenaka/comi/application/ebitenui/model"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
)

// Game object used by ebiten
type game struct {
	ui *ebitenui.UI
}

func main() {
	//如果参数当中有--debug，则不启动UI
	debugMode := flag.Bool("debug", false, "Disable UI by debug mode.")
	// 解析命令行参数
	flag.Parse()
	if *debugMode {
		fmt.Println("Debug Mode, no thing to do.")
		select {}
		return
	}
	// 启动Comigo Web服务器
	// 如果参数当中有--debug，则不启动UI，且什么都不做（UI抢焦点，所以用另一台机器rsync同步代码跑效果）
	StartComigoWebserver()

	// 创建一个新的 ReaderConfig 对象，用于配置阅读器的设置。
	readerConfig := model.NewReaderConfig()
	readerConfig.SetTitle("Comigo v0.9").
		// 阅读器模式。
		SetReaderMode(model.ScrollMode).
		// 窗口是否全屏。
		SetWindowFullScreen(true).
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

	// 创建UI
	eui, closeUI, err := createUI(readerConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer closeUI()
	// 创建游戏对象
	g := game{
		ui: eui,
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
	// 2.根容器与布局设置，
	//参考了 https://github.com/ebitenui/ebitenui/tree/master/_examples/widget_demos/gridlayout/main.go
	// 为此 UI 创建根容器。所有其他 UI 元素都必须添加到此容器中。
	// 在 Ebiten UI 中，小部件通常不是手动布局的。相反，它们作为子部件被组合在一个容器中，容器的布局器负责布局这些部件。
	rootContainer :=
		widget.NewContainer(
			//使用纯色作为背景
			widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{A: 255})),
			// 容器将使用锚布局来布局其单个子窗口小部件
			widget.ContainerOpts.Layout(
				//GridLayout 可以在网格中布置任意数量的小部件。它可以为每个网格单元的小部件设置不同的位置，并且还可以拉伸它们。
				// 根据小部件的实现方式，某些选项是必须指定的（例如按钮的图片），而其他选项则是可选的。选项的顺序通常无关紧要。有些选项可以被多次指定。
				widget.NewGridLayout(
					// 使用 Columns 参数来定义列的数量。
					widget.GridLayoutOpts.Columns(1),
					// 使用 ColumnStretch 和 RowStretch 参数来分别定义列和行的拉伸因子。
					// 只支持布尔值，true表示拉伸，false表示不拉伸。
					widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
					//网格布局的间距，c 列间距，r行间距。
					widget.GridLayoutOpts.Spacing(20, 0),
					// Padding 定义了网格块的外间距大小。
					widget.GridLayoutOpts.Padding(widget.Insets{
						Top:    0,
						Left:   0,
						Bottom: 0,
						Right:  0,
					}),
				)))

	//这会将根容器添加到 UI，以便将其展示。
	eui := &ebitenui.UI{
		Container: rootContainer,
	}

	rootContainer.AddChild(components.HeaderContainer(readerConfig))

	rootContainer.AddChild(components.BodyContainer())

	rootContainer.AddChild(components.FooterContainer())

	return eui, func() {
		// TODO:在结束时，关闭资源。
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
