package pages

import (
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"image/color"
)

// ScrollMode  滚动模式
func ScrollMode() *widget.Container {
	// 2.根容器与布局设置，
	//参考了 https://github.com/ebitenui/ebitenui/tree/master/_examples/widget_demos/gridlayout/main.go
	// 为此 UI 创建根容器。所有其他 UI 元素都必须添加到此容器中。
	// 在 Ebiten UI 中，小部件通常不是手动布局的。相反，它们作为子部件被组合在一个容器中，容器的布局器负责布局这些部件。
	ScrollMode := widget.NewContainer(
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
	return ScrollMode
}
