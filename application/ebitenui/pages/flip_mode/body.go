package flip_mode

import (
	"bytes"
	"fmt"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yumenaka/comigo/entity"
	"github.com/yumenaka/comigo/util"
	fileutil "github.com/yumenaka/comigo/util/file"
	"image/color"
	"log"
)

var nowPage int = 1

func BodyContainer() widget.PreferredSizeLocateableWidget {
	bodyContainer := widget.NewContainer(
		// 设置容器的背景图像。#E0D9CD
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{R: 0xE0, G: 0xD9, B: 0xCD, A: 0xff})),
		// 设置容器的通用选项
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(50, 50),
		),
		// 该容器将使用锚点布局来布局其单个子部件
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			//“如何布置孩子们的方向”
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			//Set how much padding before displaying content
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
			//Set how far apart to space the children
			widget.RowLayoutOpts.Spacing(15),
		)),
	)

	randomBook, err := entity.GetRandomBook()
	if err != nil {
		log.Print(err)
	}
	for k, i := range randomBook.Pages.Images {
		if k > nowPage {
			break
		}
		println(fmt.Sprintf("GetFile: %s", i.NameInArchive))
		option := fileutil.GetPictureDataOption{
			PictureName:      i.NameInArchive,
			BookIsPDF:        randomBook.Type == entity.TypePDF,
			BookIsDir:        randomBook.Type == entity.TypeDir,
			BookIsNonUTF8Zip: randomBook.NonUTF8Zip,
			BookFilePath:     randomBook.FilePath,
			ResizeMaxWidth:   800,
			ResizeMaxHeight:  500,
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
		)
		// 获取首选大小
		w, h := g.PreferredSize()
		fmt.Println(w, h)
		bodyContainer.AddChild(g)
	}

	return bodyContainer
}
