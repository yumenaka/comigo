package components

import (
	"bytes"
	"fmt"
	"image/color"

	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/skip2/go-qrcode"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/util/logger"
)

// QRCodeWindow QRCode弹出窗口
func QRCodeWindow() *widget.Window {
	// QrCode窗口Container
	QrCodeWindowContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(color.NRGBA{R: 25, G: 25, B: 25, A: 255})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	//QRCode图片
	qrcodeStr := config.GetQrcodeURL()
	qrcodeData, err := qrcode.Encode(qrcodeStr, qrcode.Medium, 256)
	if err != nil {
		logger.Infof("%s", err)
	}
	// 从图像数据中创建一个新的 Reader
	qrcodeReader := bytes.NewReader(qrcodeData)
	// 从图像数据中创建一个新的图像
	qrcodeImage, _, err := ebitenutil.NewImageFromReader(qrcodeReader)
	if err != nil {
		println(fmt.Sprintf("NewImageFromReader error:%s", err))
	}
	qrcodeGraphic := widget.NewGraphic(
		// 设置图像
		widget.GraphicOpts.Image(qrcodeImage),
		// 设置小部件的通用选项
		widget.GraphicOpts.WidgetOpts(
			//要配置单个小部件与其兄弟小部件有不同的布局，可以在小部件上设置一个可选的“布局数据”。
			//布局数据的类型取决于所使用的布局实现。例如，AnchorLayout 需要使用 AnchorLayoutData。
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				//设置图片的水平锚定位置
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				//设置图片的垂直锚定位置
				VerticalPosition: widget.AnchorLayoutPositionCenter,
			}),
		),
	)
	// 将QrCode图片添加到窗口容器
	QrCodeWindowContainer.AddChild(qrcodeGraphic)
	// 为二维码窗口标题加载字体
	titleFace, _ := loadFont(16)
	// 为QRCode窗口创建标题栏
	titleContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(color.NRGBA{R: 150, G: 150, B: 150, A: 255})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	// QRCode窗口标题
	titleContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("Scan QRCode to read", titleFace, color.NRGBA{R: 255, G: 255, B: 255, A: 255}),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
	))

	// 创建弹出窗口对象。窗口对象不依附于任何容器。它的位置和
	// 大小需通过在窗口上使用 SetLocation 方法手动设置，并使用 ui.AddWindow() 方法将其添加到 UI 中
	// 设置下面的按钮回调，以查看窗口是如何被添加到 UI 中的。
	QRCodeWindow := widget.NewWindow(
		//设置窗口的主要内容
		widget.WindowOpts.Contents(QrCodeWindowContainer),
		//为窗口设置标题栏（可选）
		widget.WindowOpts.TitleBar(titleContainer, 32),
		//将窗口置于所有其他内容之上，并阻止其他地方的输入
		widget.WindowOpts.Modal(),
		//设置如何关闭窗口。CLICK_OUT将在点击窗口对象以外的任何地方时关闭窗口
		widget.WindowOpts.CloseMode(widget.CLICK_OUT),
		//表示窗口可拖动。这必须有一个标题栏才能工作。
		widget.WindowOpts.Draggable(),
		//设置窗口可调整大小
		widget.WindowOpts.Resizeable(),
		//设置窗口的最小大小
		widget.WindowOpts.MinSize(266, 296),
		//设置窗口的最大大小
		widget.WindowOpts.MaxSize(266, 296),
		//设置移动完成时触发的回调
		widget.WindowOpts.MoveHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Window Moved")
		}),
		//设置在调整大小完成后触发的回调
		widget.WindowOpts.ResizeHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Window Resized")
		}),
	)
	return QRCodeWindow
}
