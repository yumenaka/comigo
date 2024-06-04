package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yumenaka/comi/cmd"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/routers"
	"golang.org/x/image/font/gofont/goregular"
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
		SetWindowFullScreen(true).                                      //SetWindowFullScreen 设置窗口是否全屏。
		SetWindowDecorated(true).                                       //SetWindowDecorated 设置是否有边框和标题栏
		SetWindowResizingModeEnabled(ebiten.WindowResizingModeEnabled). //SetWindowResizingModeEnabled 设置窗口是否可以调整大小。
		SetWindowSize(1280, 800).
		SetRunOptions(ebiten.RunGameOptions{
			ScreenTransparent: false,
		})
	// 设置窗口大小和标题等
	ebiten.SetWindowSize(readerConfig.Width, readerConfig.Height)
	ebiten.SetWindowTitle(readerConfig.Title)
	ebiten.SetWindowResizingMode(readerConfig.WindowResizingModeEnabled)
	// SetWindowDecorated 设置窗口是否有边框和标题栏。
	ebiten.SetWindowDecorated(readerConfig.WindowDecorated)
	// SetWindowFullScreen 设置窗口是否全屏。
	ebiten.SetFullscreen(readerConfig.WindowFullScreen)
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
					widget.GridLayoutOpts.Columns(1),
					// 使用 ColumnStretch 和 RowStretch 参数来分别定义列和行的拉伸因子。
					// 只支持布尔值，true表示拉伸，false表示不拉伸。
					widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
					//间距配置网格布局，以通过间距 c 分隔列，并通过间距 r 分隔行。
					widget.GridLayoutOpts.Spacing(20, 20),
					// Padding 定义了网格块的外间距大小。
					widget.GridLayoutOpts.Padding(widget.Insets{
						Top:    10,
						Left:   10,
						Bottom: 10,
						Right:  10,
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
	headerContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{255, 0, 0, 255})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(50, 50),
		),
		// 容器有布局的概念。这就是这个容器的子级小部件的布局设置：
		// 它们将被放置在容器的边界内。
		// 容器将使用锚布局来布局其单个子窗口小部件
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			//外边与内部内容之间的的填充(padding)的大小
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(40)),
		)),
	)
	headerText := widget.NewText(
		widget.TextOpts.Text("Header", fontFace, textColor),
		// WidgetOpts 用于设置小部件的各种属性。这里用来设置文本的锚点布局。
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			// 指定网格单元内的水平锚定位置。
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			// 指定网格单元格内的垂直锚定位置。
			VerticalPosition: widget.AnchorLayoutPositionCenter,
		})),
	)
	headerContainer.AddChild(headerText)
	rootContainer.AddChild(headerContainer)

	bodyContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{77, 77, 77, 255})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(50, 50),
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				//指定网格单元内的水平锚定位置。
				//HorizontalPosition: widget.GridLayoutPositionCenter,
				// 指定网格单元格内的垂直锚定位置。
				//VerticalPosition: widget.GridLayoutPositionStart,
				// 限制最大大小。
				//MaxWidth:  300,
				//MaxHeight: 300,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	bodyText := widget.NewText(
		widget.TextOpts.Text("Body", fontFace, textColor),
		// WidgetOpts 用于设置小部件的各种属性。这里用来设置文本的锚点布局。
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			// 指定网格单元内的水平锚定位置。
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			// 指定网格单元格内的垂直锚定位置。
			VerticalPosition: widget.AnchorLayoutPositionCenter,
		})),
	)
	bodyContainer.AddChild(bodyText)
	rootContainer.AddChild(bodyContainer)

	bottomContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{R: 0, G: 0, B: 255, A: 255})),
		widget.ContainerOpts.WidgetOpts(
			//此单元格中的小部件的 MaxHeight 和 MaxWidth 小于
			//网格单元的大小，因此它将使用下面的位置字段
			//确定小部件应在该网格单元中显示的位置。
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				// 设置水平和垂直位置。
				//HorizontalPosition: widget.GridLayoutPositionCenter,
				//VerticalPosition:   widget.GridLayoutPositionCenter,
				// 限制最大大小。
				//MaxWidth:  300,
				//MaxHeight: 300,
			}),
			widget.WidgetOpts.MinSize(100, 50),
		),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	bottomText := widget.NewText(
		widget.TextOpts.Text("Bottom", fontFace, textColor),
		// WidgetOpts 用于设置小部件的各种属性。这里用来设置文本的锚点布局。
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			// 指定网格单元内的水平锚定位置。
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			// 指定网格单元格内的垂直锚定位置。
			VerticalPosition: widget.AnchorLayoutPositionCenter,
		})),
	)
	bottomContainer.AddChild(bottomText)
	rootContainer.AddChild(bottomContainer)

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
	// 在屏幕左上角显示当前 FPS。
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.ActualFPS()))
}

func (g *game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
