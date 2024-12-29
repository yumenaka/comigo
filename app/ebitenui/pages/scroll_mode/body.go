package scroll_mode

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util"
	fileutil "github.com/yumenaka/comigo/util/file"
)

func BodyContainer() widget.PreferredSizeLocateableWidget {
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

	randomBook, err := model.GetRandomBook()
	if err != nil {
		log.Print(err)
	}
	for _, i := range randomBook.Pages.Images {
		println(fmt.Sprintf("GetFile: %s", i.NameInArchive))
		option := fileutil.GetPictureDataOption{
			PictureName:      i.NameInArchive,
			BookIsPDF:        randomBook.Type == model.TypePDF,
			BookIsDir:        randomBook.Type == model.TypeDir,
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
	return bodyContainer
}
