package components

import (
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/yumenaka/comigo/config"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
	"log"
)

func FooterContainer() widget.PreferredSizeLocateableWidget {
	// 加载按钮文字所需的字体
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal("Error Parsing Font", err)
	}
	// 设置字体大小
	fontFace := truetype.NewFace(ttfFont, &truetype.Options{
		Size: 18,
	})
	// 设置文本颜色
	textColor := color.RGBA{R: 0, G: 0, B: 0, A: 0xff}
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
		widget.TextOpts.Text("Power by Comigo "+config.GetVersion(), fontFace, textColor),
		// WidgetOpts 用于设置小部件的各种属性。这里用来设置文本的锚点布局。
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			// 指定网格单元内的水平锚定位置。
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			// 指定网格单元格内的垂直锚定位置。
			VerticalPosition: widget.AnchorLayoutPositionCenter,
		})),
	)
	footerContainer.AddChild(footerText)
	return footerContainer
}
