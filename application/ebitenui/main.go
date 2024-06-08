package main

import (
	"bytes"
	"flag"
	"fmt"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yumenaka/comi/cmd"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/entity"
	"github.com/yumenaka/comi/routers"
	fileutil "github.com/yumenaka/comi/util/file"
	"golang.org/x/image/font"
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
	readerConfig.SetTitle("Sample v1.0").
		// 阅读器模式。
		SetReaderMode(ScrollMode).
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

	// 2.根容器与布局设置，
	//参考了 https://github.com/ebitenui/ebitenui/tree/master/_examples/widget_demos/gridlayout/main.go
	// 为此 UI 创建根容器。
	// 所有其他 UI 元素都必须添加到此容器中。
	rootContainer :=
		widget.NewContainer(
			//使用纯色作为背景
			widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{A: 255})),
			// 容器将使用锚布局来布局其单个子窗口小部件
			widget.ContainerOpts.Layout(
				//GridLayout 网格布局模式，将小部件放置在网格中。
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

	// 加载字体
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal("Error Parsing Font", err)
	}
	// 字体大小
	fontFace := truetype.NewFace(ttfFont, &truetype.Options{
		Size: 20,
	})
	// 文本颜色
	textColor := color.RGBA{R: 0, G: 0, B: 0, A: 0xff}

	// 一个新的文本小部件，用于显示一些文本。
	headerContainer := widget.NewContainer(
		// header容器的背景颜色
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{R: 245, G: 245, B: 228, A: 255})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(50, 30),
		),
		// 容器有布局的概念。这就是这个容器的子级小部件的布局设置：
		// 它们将被放置在容器的边界内。
		// widget.NewAnchorLayout() 创建一个新的锚点布局。锚点布局不会处理重叠。要自定义重叠顺序的话，可以用widget.NewStackedLayout()
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			//外边与内部内容之间的的填充(padding)的大小
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(15)),
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
		// 设置容器的背景图像。#E0D9CD
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{R: 0xE0, G: 0xD9, B: 0xCD, A: 0xff})),
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
		// 容器将使用网格布局来配置 ScrollableContainer 和 Slider
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Spacing(2, 0),
			widget.GridLayoutOpts.Stretch([]bool{true, false}, []bool{true}),
		)),
	)

	// 创建一个包含滚动内容的容器
	bodyContent := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewRowLayout(
		// 设置行布局的方向
		widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		// 设置行布局的填充
		widget.RowLayoutOpts.Padding(widget.Insets{
			Left:   0,
			Right:  0,
			Top:    30,
			Bottom: 30,
		}),
		// 设置行布局的间距
		widget.RowLayoutOpts.Spacing(10),
	)))

	randomBook, err := entity.GetRandomBook()
	if err != nil {
		log.Print(err)
	}

	for _, i := range randomBook.Pages.Images {
		//if _ > 10 {
		//	break
		//}
		println(fmt.Sprintf("GetFile: %s", i.NameInArchive))
		option := fileutil.GetPictureDataOption{
			PictureName:      i.NameInArchive,
			BookIsPDF:        randomBook.Type == entity.TypePDF,
			BookIsDir:        randomBook.Type == entity.TypeDir,
			BookIsNonUTF8Zip: randomBook.NonUTF8Zip,
			BookFilePath:     randomBook.FilePath,
		}
		imgData, _, err := fileutil.GetPictureData(option)
		if err != nil {
			println(fmt.Sprintf("GetPictureData error:%s", err))
		}
		reader := bytes.NewReader(imgData)
		// 从图像数据中创建一个新的图像
		eImage, _, err := ebitenutil.NewImageFromReader(reader)
		if err != nil {
			println(fmt.Sprintf("NewImageFromReader error:%s", err))
		}
		g := widget.NewGraphic(
			// 设置小部件的通用选项
			widget.GraphicOpts.WidgetOpts(
				// 锚点布局，水平和垂直居中
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				}),
			),
			widget.GraphicOpts.Image(eImage),
			widget.GraphicOpts.WidgetOpts(
				widget.WidgetOpts.MinSize(600, 480),
			),
		)
		// 获取首选大小
		w, h := g.PreferredSize()
		fmt.Println(w, h)

		bodyContent.AddChild(g)
	}

	//// 加载按钮状态的图片：静止、悬停和按下(idle, hover, and pressed)。
	//buttonImage, _ := loadButtonImage()
	//// 加载按钮文字字体
	//face, _ := loadFont(20)
	//// 将N个按钮添加到可滚动内容容器中
	//for x := 0; x < 100; x++ {
	//	// Capture x for use in callback
	//	x := x
	//	// construct a button
	//	button := widget.NewButton(
	//		// 设置小部件的通用选项
	//		widget.ButtonOpts.WidgetOpts(
	//			// 指示容器的锚点布局，将按钮水平和垂直居中
	//			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
	//				Position: widget.RowLayoutPositionCenter,
	//			}),
	//		),
	//
	//		// 指定要使用的图像
	//		widget.ButtonOpts.Image(buttonImage),
	//
	//		// 指定按钮的文本、字体和颜色
	//		widget.ButtonOpts.Text(fmt.Sprintf("Awesome! %d", x), face, &widget.ButtonTextColor{
	//			Idle: color.NRGBA{R: 0xdf, G: 0xf4, B: 0xff, A: 0xff},
	//		}),
	//
	//		// 指定按钮的文本需要一些填充才能正确显示
	//		widget.ButtonOpts.TextPadding(widget.Insets{
	//			Left:   300,
	//			Right:  300,
	//			Top:    120,
	//			Bottom: 120,
	//		}),
	//
	//		// 添加一个处理程序以响应点击按钮事件
	//		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
	//			println(fmt.Sprintf("Button %d Clicked!", x))
	//		}),
	//	)
	//
	//	// 将按钮添加为容器的子元素
	//	bodyContent.AddChild(button)
	//}

	// 创建新的 ScrollContainer 对象
	scrollContainer := widget.NewScrollContainer(
		// 设置将要滚动的内容
		widget.ScrollContainerOpts.Content(bodyContent),
		// 让容器将内容宽度拉伸以匹配可用空间
		widget.ScrollContainerOpts.StretchContentWidth(),
		// 为可滚动容器设置背景图像。
		widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
			Idle: image.NewNineSliceColor(color.NRGBA{0xE0, 0xD9, 0xCD, 0xff}),
			Mask: image.NewNineSliceColor(color.NRGBA{0xE0, 0xD9, 0xCD, 0xff}),
		}),
	)
	// 将可滚动容器添加到左侧网格单元中
	bodyContainer.AddChild(scrollContainer)

	// 创建一个函数来返回滑块使用的页面大小
	pageSizeFunc := func() int {
		return int(math.Round(float64(scrollContainer.ViewRect().Dy()) / float64(bodyContent.GetWidget().Rect.Dy()) * 1000))
	}
	// 创建一个垂直滑块来控制可滚动容器
	vSlider := widget.NewSlider(
		widget.SliderOpts.Direction(widget.DirectionVertical),
		widget.SliderOpts.MinMax(0, 1000),
		widget.SliderOpts.PageSizeFunc(pageSizeFunc),
		// 根据滑块的值更新滚动位置
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			scrollContainer.ScrollTop = float64(args.Slider.Current) / 1000
		}),
		widget.SliderOpts.Images(
			// 设置轨道图片
			&widget.SliderTrackImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			// 设置滑动块图片
			&widget.ButtonImage{
				Idle:    image.NewNineSliceColor(color.NRGBA{150, 150, 235, 255}),
				Hover:   image.NewNineSliceColor(color.NRGBA{150, 150, 235, 255}),
				Pressed: image.NewNineSliceColor(color.NRGBA{150, 150, 235, 255}),
			},
		),
	)
	// 如果滚动容器通过滑块以外的其他方式进行滚动，设置滑块的位置。
	scrollContainer.GetWidget().ScrolledEvent.AddHandler(func(args interface{}) {
		a := args.(*widget.WidgetScrolledEventArgs)
		p := pageSizeFunc() / 3
		if p < 1 {
			p = 1
		}
		vSlider.Current -= int(math.Round(a.Y * float64(p)))
	})

	// 将滑块添加到bodyContainer容器的第二个列插槽中
	bodyContainer.AddChild(vSlider)
	rootContainer.AddChild(bodyContainer)

	footerContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{245, 245, 228, 255})),
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
	footerText := widget.NewText(
		widget.TextOpts.Text("Footer", fontFace, textColor),
		// WidgetOpts 用于设置小部件的各种属性。这里用来设置文本的锚点布局。
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			// 指定网格单元内的水平锚定位置。
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			// 指定网格单元格内的垂直锚定位置。
			VerticalPosition: widget.AnchorLayoutPositionCenter,
		})),
	)
	footerContainer.AddChild(footerText)
	rootContainer.AddChild(footerContainer)

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

func loadButtonImage() (*widget.ButtonImage, error) {
	idle := image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})

	hover := image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})

	pressed := image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}

func loadFont(size float64) (font.Face, error) {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}
