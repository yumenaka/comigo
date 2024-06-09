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
	"github.com/yumenaka/comi/util"
	fileutil "github.com/yumenaka/comi/util/file"
	_ "golang.org/x/image/bmp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	_ "golang.org/x/image/webp"
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
		config.Config.OpenBrowser = false
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
	readerConfig.SetTitle("Comigo v0.9").
		// 阅读器模式。
		SetReaderMode(ScrollMode).
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

	// 加载按钮文字所需的字体
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal("Error Parsing Font", err)
	}
	// 设置字体大小
	fontFace := truetype.NewFace(ttfFont, &truetype.Options{
		Size: 20,
	})
	// 设置文本颜色
	textColor := color.RGBA{R: 0, G: 0, B: 0, A: 0xff}
	// 加载按钮状态的图片：静止、悬停和按下(idle, hover, and pressed)。
	buttonImage, _ := loadButtonImage()
	// 加载按钮文字字体
	face, _ := loadFont(20)
	// headerContainer 是一个新的容器，用于包含标题文本和按钮。
	headerContainer := widget.NewContainer(
		// header容器的背景颜色
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{R: 245, G: 245, B: 228, A: 255})),
		// 设置容器的布局
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			// 使用 Columns 参数来定义列的数量。
			widget.GridLayoutOpts.Columns(7),
			// 使用 ColumnStretch 和 RowStretch 参数来分别定义列和行的拉伸因子。
			// 只支持布尔值，true表示拉伸，false表示不拉伸。
			widget.GridLayoutOpts.Stretch([]bool{false, false, false, true, false, false, false}, []bool{true}),
			//网格布局的间距，c 列间距，r行间距。
			widget.GridLayoutOpts.Spacing(2, 0),
		)),
		// 设置容器的通用选项
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(10, 10),
		),
	)
	// 服务器设置按钮
	serverButton := widget.NewButton(
		// 指定要使用的图像
		widget.ButtonOpts.Image(buttonImage),
		// 指定按钮的文本、字体和颜色
		widget.ButtonOpts.Text(fmt.Sprintf("Server"), face, &widget.ButtonTextColor{
			Idle: color.NRGBA{R: 0xdf, G: 0xf4, B: 0xff, A: 0xff},
		}),
		// 指定按钮的文本需要一些填充才能正确显示
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   10,
			Right:  10,
			Top:    10,
			Bottom: 10,
		}),
		// 添加一个处理程序以响应点击按钮事件
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println(fmt.Sprintf("Server Button Clicked!"))
		}),
		// 设置按钮的通用选项
		widget.ButtonOpts.WidgetOpts(
			// 布局设置，将按钮水平和垂直居中
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)
	headerContainer.AddChild(serverButton)

	// Upload按钮
	uploadButton := widget.NewButton(
		// 指定要使用的图像
		widget.ButtonOpts.Image(buttonImage),
		// 指定按钮的文本、字体和颜色
		widget.ButtonOpts.Text(fmt.Sprintf("Upload"), face, &widget.ButtonTextColor{
			Idle: color.NRGBA{R: 0xdf, G: 0xf4, B: 0xff, A: 0xff},
		}),
		// 指定按钮的文本需要一些填充才能正确显示
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   10,
			Right:  10,
			Top:    10,
			Bottom: 10,
		}),
		// 添加一个处理程序以响应点击按钮事件
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println(fmt.Sprintf("Upload Button Clicked!"))
		}),
		// 设置按钮的通用选项
		widget.ButtonOpts.WidgetOpts(
			// 布局设置，将按钮水平和垂直居中
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)
	headerContainer.AddChild(uploadButton)

	// Sort按钮
	sortButton := widget.NewButton(
		// 指定要使用的图像
		widget.ButtonOpts.Image(buttonImage),
		// 指定按钮的文本、字体和颜色
		widget.ButtonOpts.Text(fmt.Sprintf("Sort"), face, &widget.ButtonTextColor{
			Idle: color.NRGBA{R: 0xdf, G: 0xf4, B: 0xff, A: 0xff},
		}),
		// 指定按钮的文本需要一些填充才能正确显示
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   10,
			Right:  10,
			Top:    10,
			Bottom: 10,
		}),
		// 添加一个处理程序以响应点击按钮事件
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println(fmt.Sprintf("Sort Button Clicked!"))
		}),
		// 设置按钮的通用选项
		widget.ButtonOpts.WidgetOpts(
			// 布局设置，将按钮水平和垂直居中
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)
	headerContainer.AddChild(sortButton)

	// 一个新的文本小部件，用于显示文本。
	titleText := widget.NewText(
		widget.TextOpts.Text("Title", fontFace, textColor),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		//要配置单个小部件与其兄弟小部件有不同的布局，可以在小部件上设置一个可选的“布局数据”。
		//布局数据的类型取决于所使用的布局实现。例如，RowLayout 需要使用 RowLayoutData。
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{}),
		),
	)
	headerContainer.AddChild(titleText)

	// QRCode按钮
	qrcodeButton := widget.NewButton(
		// 指定要使用的图像
		widget.ButtonOpts.Image(buttonImage),
		// 指定按钮的文本、字体和颜色
		widget.ButtonOpts.Text(fmt.Sprintf("QRCode"), face, &widget.ButtonTextColor{
			Idle: color.NRGBA{R: 0xdf, G: 0xf4, B: 0xff, A: 0xff},
		}),
		// 指定按钮的文本需要一些填充才能正确显示
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   10,
			Right:  10,
			Top:    10,
			Bottom: 10,
		}),
		// 添加一个处理程序以响应点击按钮事件
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println(fmt.Sprintf("Sort Button Clicked!"))
		}),
		// 设置按钮的通用选项
		widget.ButtonOpts.WidgetOpts(
			// 布局设置，将按钮水平和垂直居中
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)
	headerContainer.AddChild(qrcodeButton)

	// FullScreen按钮
	fullScreenButton := widget.NewButton(
		// 指定要使用的图像
		widget.ButtonOpts.Image(buttonImage),
		// 指定按钮的文本、字体和颜色
		widget.ButtonOpts.Text(fmt.Sprintf("FullScreen"), face, &widget.ButtonTextColor{
			Idle: color.NRGBA{R: 0xdf, G: 0xf4, B: 0xff, A: 0xff},
		}),
		// 指定按钮的文本需要一些填充才能正确显示
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   10,
			Right:  10,
			Top:    10,
			Bottom: 10,
		}),
		// 添加一个处理程序以响应点击按钮事件
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println(fmt.Sprintf("FullScreen Button Clicked!"))
			readerConfig.WindowFullScreen = !readerConfig.WindowFullScreen
			ebiten.SetFullscreen(readerConfig.WindowFullScreen)
		}),
		// 设置按钮的通用选项
		widget.ButtonOpts.WidgetOpts(
			// 布局设置，将按钮水平和垂直居中
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)
	headerContainer.AddChild(fullScreenButton)

	// 设置按钮
	settingButton := widget.NewButton(
		// 指定要使用的图像
		widget.ButtonOpts.Image(buttonImage),
		// 指定按钮的文本、字体和颜色
		widget.ButtonOpts.Text(fmt.Sprintf("Setting"), face, &widget.ButtonTextColor{
			Idle: color.NRGBA{R: 0xdf, G: 0xf4, B: 0xff, A: 0xff},
		}),
		// 指定按钮的文本需要一些填充才能正确显示
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   10,
			Right:  10,
			Top:    10,
			Bottom: 10,
		}),
		// 添加一个处理程序以响应点击按钮事件
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println(fmt.Sprintf("Setting Button Clicked!"))
		}),
		// 设置按钮的通用选项
		widget.ButtonOpts.WidgetOpts(
			// 锚点布局设置，将按钮水平和垂直居中
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)
	headerContainer.AddChild(settingButton)

	rootContainer.AddChild(headerContainer)

	bodyContainer := widget.NewContainer(
		// 设置容器的背景图像。#E0D9CD
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{R: 0xE0, G: 0xD9, B: 0xCD, A: 0xff})),
		// 设置容器的通用选项
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
	// RowLayout 可以在一行或一列中布置任意数量的小部件。它还可以根据需要对小部件进行不同的定位和拉伸。
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
		widget.RowLayoutOpts.Spacing(5),
	)))

	randomBook, err := entity.GetRandomBook()
	if err != nil {
		log.Print(err)
	}
	for _, i := range randomBook.Pages.Images {
		println(fmt.Sprintf("GetFile: %s", i.NameInArchive))
		option := fileutil.GetPictureDataOption{
			PictureName:      i.NameInArchive,
			BookIsPDF:        randomBook.Type == entity.TypePDF,
			BookIsDir:        randomBook.Type == entity.TypeDir,
			BookIsNonUTF8Zip: randomBook.NonUTF8Zip,
			BookFilePath:     randomBook.FilePath,
			ResizeMaxWidth:   800,
			ResizeMaxHeight:  1000,
		}
		imgData, _, err := fileutil.GetPictureData(option)
		if err != nil {
			println(fmt.Sprintf("GetPictureData error:%s", err))
		}
		// 限制图片大小(宽度)
		tempData, limitErr := util.ImageResizeByMaxWidth(imgData, option.ResizeMaxWidth)
		if limitErr != nil {
			println(fmt.Sprintf(limitErr.Error()))
		} else {
			imgData = tempData
		}
		// 限制图片大小(高度)
		tempData, limitErr = util.ImageResizeByMaxHeight(imgData, option.ResizeMaxHeight)
		if limitErr != nil {
			println(fmt.Sprintf(limitErr.Error()))
		} else {
			imgData = tempData
		}
		// 从图像数据中创建一个新的 Reader
		reader := bytes.NewReader(imgData)
		// 从图像数据中创建一个新的图像
		eImage, _, err := ebitenutil.NewImageFromReader(reader)
		if err != nil {
			println(fmt.Sprintf("NewImageFromReader error:%s", err))
		}
		g := widget.NewGraphic(
			// 设置图像
			widget.GraphicOpts.Image(eImage),
			// 设置小部件的通用选项
			widget.GraphicOpts.WidgetOpts(
				//要配置单个小部件与其兄弟小部件有不同的布局，可以在小部件上设置一个可选的“布局数据”。
				//布局数据的类型取决于所使用的布局实现。例如，AnchorLayout 需要使用 AnchorLayoutData。
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					MaxWidth: 800,
					//MaxHeight: 480,
					// 锚点布局，水平和垂直居中
					Position: widget.RowLayoutPositionCenter,
				}),
			),
		)
		// 获取首选大小
		w, h := g.PreferredSize()
		fmt.Println(w, h)
		bodyContent.AddChild(g)
	}

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
