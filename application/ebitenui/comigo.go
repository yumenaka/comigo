package main

import (
	"flag"
	"fmt"
	"github.com/yumenaka/comi/cmd"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/routers"
	"image/color"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font/gofont/goregular"
)

type game struct {
	ui *ebitenui.UI
}

func main() {
	//解析命令，扫描文件
	cmd.StartScan(os.Args)
	//设置临时文件夹
	config.SetTempDir()
	//SetWebServerPort
	routers.SetWebServerPort()
	//设置书籍API
	routers.StartWebServer()
	// 如果参数当中有--debug，则不启动UI
	debugMode := flag.Bool("debug", false, "Disable UI by debug mode.")
	// 解析命令行参数
	flag.Parse()
	// 根据 debugMode 的值决定后续逻辑
	if *debugMode {
		fmt.Println("Debug Mode, UI is disabled.")
		select {}
		return
	}
	fmt.Println("UI is enabled.")

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

	// 为此 UI 创建根容器。
	// 所有其他 UI 元素都必须添加到此容器中。
	rootContainer :=
		widget.NewContainer(
			// 根容器的布局设置。
			widget.ContainerOpts.Layout(
				//GridLayout 网格布局模式，将小部件放置在网格中。
				widget.NewGridLayout(
					// 使用 Columns 参数来定义列的数量。
					widget.GridLayoutOpts.Columns(2),
					// 使用 ColumnStretch 和 RowStretch 参数来分别定义列和行的拉伸因子。
					// 只支持布尔值，true表示拉伸，false表示不拉伸。
					widget.GridLayoutOpts.Stretch([]bool{true, true, true, true}, []bool{true, false, true, false}),
					// Padding 定义了网格块的外间距大小。
					widget.GridLayoutOpts.Padding(widget.Insets{
						Top:    30,
						Left:   20,
						Bottom: 20,
						Right:  20,
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
		Size: 24,
	})

	for i := range 40 {
		// 文本颜色。
		rgba := color.RGBA{R: uint8((250 - 3*i) % 255), G: uint8((150 + 7*i) % 255), B: uint8((70 + i*i) % 255), A: 0xff}
		// 创建一个文本小部件，上面写“World_i”+现在时间
		nowTime := time.Now().Format("2006-01-02 15:04:05")
		helloWorldLabel := widget.NewText(
			widget.TextOpts.Text(strconv.Itoa(i)+" "+nowTime, fontFace, rgba),
		)
		// 要显示文本小部件，将其添加到根容器中。
		rootContainer.AddChild(helloWorldLabel)
	}

	game := game{
		ui: eui,
	}
	err = ebiten.RunGameWithOptions(&game, &readerConfig.RunOptions)
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
