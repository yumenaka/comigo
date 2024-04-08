package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"log"

	"golang.org/x/image/font/gofont/goregular"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
)

type game struct {
	ui *ebitenui.UI
}

func main() {
	ebiten.SetWindowSize(900, 800)
	ebiten.SetWindowTitle("Hello World")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	// 为此 UI 创建根容器。
	// 所有其他 UI 元素都必须添加到此容器中。
	rootContainer :=
		widget.NewContainer(
			// 根容器的布局设置。
			widget.ContainerOpts.Layout(
				// 使用 GridLayout 布局来排列子元素。
				widget.NewGridLayout(
					widget.GridLayoutOpts.Columns(1),
					//使用 Stretch 参数来定义行的布局方式。
					widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
					// Padding defines how much space to put around the outside of the grid.
					widget.GridLayoutOpts.Padding(widget.Insets{
						Top:    30,
						Bottom: 20,
					}),
				)))

	//这会将根容器添加到 UI，以便将其展示。
	eui := &ebitenui.UI{
		Container: rootContainer,
	}

	// 这会加载字体并创建字体。
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal("Error Parsing Font", err)
	}
	fontFace := truetype.NewFace(ttfFont, &truetype.Options{
		Size: 32,
	})

	// 这将创建一个文本小部件，上面写着“Hello World！”
	helloWorldLabel := widget.NewText(
		widget.TextOpts.Text("Hello World!", fontFace, color.White),
	)

	// 要显示文本小部件，我们必须将其添加到根容器中。
	rootContainer.AddChild(helloWorldLabel)

	game := game{
		ui: eui,
	}

	err = ebiten.RunGame(&game)
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
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.ActualFPS()))
}

func (g *game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
