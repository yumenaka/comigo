package main

import (
	"flag"
	"fmt"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yumenaka/comi/cmd"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/routers"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
	"log"
	"os"
)

// Game object used by ebiten
type game struct {
	ui *ebitenui.UI
}

func main() {
	//0. 如果参数当中有--debug，则不启动UI，且什么都不做（UI抢焦点，所以用另一台机器rsync同步代码跑效果）
	{
		// 如果参数当中有--debug，则不启动UI
		debugMode := flag.Bool("debug", false, "Disable UI by debug mode.")
		// 解析命令行参数
		flag.Parse()
		// 根据 debugMode 的值决定后续逻辑
		if *debugMode {
			fmt.Println("Debug Mode, no thing to do.")
			select {}
			return
		}
		fmt.Println("UI is enabled.")
		//解析命令，扫描文件
		cmd.StartScan(os.Args)
		//设置临时文件夹
		config.SetTempDir()
		//SetWebServerPort
		routers.SetWebServerPort()
		//设置书籍API
		routers.StartWebServer()
	}

	// 1. 创建一个新的 ReaderConfig 对象，用于配置阅读器的设置。
	readerConfig := NewReaderConfig()
	readerConfig.SetTitle("Comigo Reader v0.9.9").
		SetReaderMode(ScrollMode).
		SetWindowFullScreen(false).                                     //SetWindowFullScreen 设置窗口是否全屏。
		SetWindowDecorated(true).                                       //SetWindowDecorated 设置是否有边框和标题栏
		SetWindowResizingModeEnabled(ebiten.WindowResizingModeEnabled). //SetWindowResizingModeEnabled 设置窗口是否可以调整大小。
		SetWindowSize(1024, 800).
		SetRunOptions(ebiten.RunGameOptions{
			ScreenTransparent: false,
		})

	ebiten.SetWindowSize(readerConfig.Width, readerConfig.Height)
	ebiten.SetWindowTitle(readerConfig.Title)
	ebiten.SetWindowResizingMode(readerConfig.WindowResizingModeEnabled)
	// SetWindowDecorated 设置窗口是否有边框和标题栏。
	ebiten.SetWindowDecorated(readerConfig.WindowDecorated)
	ebiten.SetScreenClearedEveryFrame(false)

	// 2.根容器与布局设置，
	//参考了 https://github.com/ebitenui/ebitenui/tree/master/_examples/widget_demos/gridlayout/main.go
	// 为此 UI 创建根容器。
	// 所有其他 UI 元素都必须添加到此容器中。
	rootContainer :=
		widget.NewContainer(
			//使用纯色作为背景
			widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
			// 容器将使用锚布局来布局其单个子窗口小部件
			widget.ContainerOpts.Layout(
				//GridLayout 网格布局模式，将小部件放置在网格中。
				widget.NewGridLayout(
					// 使用 Columns 参数来定义列的数量。
					widget.GridLayoutOpts.Columns(2),
					// 使用 ColumnStretch 和 RowStretch 参数来分别定义列和行的拉伸因子。
					// 只支持布尔值，true表示拉伸，false表示不拉伸。
					widget.GridLayoutOpts.Stretch([]bool{true, true, true, true}, []bool{false, true, true, true, true, true}),
					// Padding 定义了网格块的外间距大小。
					widget.GridLayoutOpts.Padding(widget.Insets{
						Top:    15,
						Left:   20,
						Bottom: 20,
						Right:  20,
					}),
				)))

	//这会将根容器添加到 UI，以便将其展示。
	eui := &ebitenui.UI{
		Container: rootContainer,
	}

	// 加载字体
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal("Error Parsing Font", err)
	}
	// 字体大小
	fontFace := truetype.NewFace(ttfFont, &truetype.Options{
		Size: 24,
	})
	// 文本颜色
	textColor := color.RGBA{R: 255, G: 255, B: 255, A: 0xff}

	// 一个新的文本小部件，用于显示一些文本。
	innerContainer1 := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{255, 0, 0, 255})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(100, 100),
		),
	)
	label1 := widget.NewText(
		widget.TextOpts.Text("innerContainer1", fontFace, textColor),
	)
	innerContainer1.AddChild(label1)
	rootContainer.AddChild(innerContainer1)

	innerContainer2 := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0, 255, 0, 255})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(100, 100),
		),
	)
	label2 := widget.NewText(
		widget.TextOpts.Text("innerContainer2", fontFace, textColor),
	)
	innerContainer1.AddChild(label2)
	rootContainer.AddChild(innerContainer2)

	innerContainer3 := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0, 0, 255, 255})),
		widget.ContainerOpts.WidgetOpts(
			//The widget in this cell has a MaxHeight and MaxWidth less than the
			//Size of the grid cell so it will use the Position fields below to
			//Determine where the widget should be displayed within that grid cell.
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				HorizontalPosition: widget.GridLayoutPositionCenter,
				VerticalPosition:   widget.GridLayoutPositionCenter,
				MaxWidth:           100,
				MaxHeight:          100,
			}),
			widget.WidgetOpts.MinSize(100, 100),
		),
	)
	rootContainer.AddChild(innerContainer3)

	innerContainer4 := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0, 255, 255, 255})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(100, 100),
		),
	)
	rootContainer.AddChild(innerContainer4)

	g := game{
		ui: eui,
	}
	err = ebiten.RunGameWithOptions(
		&g,
		&readerConfig.RunOptions,
	)
	if err != nil {
		log.Print(err)
	}
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
	// 这只是一个调试打印，用于在屏幕左上角显示当前 FPS。
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.ActualFPS()))
}

func (g *game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
