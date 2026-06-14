//go:build !js

package model

import (
	"bytes"
	"image"

	"github.com/bbrks/go-blurhash"
	"github.com/cheggaaa/pb/v3"
	"github.com/disintegration/imaging"
	"github.com/xxjwxc/gowp/workpool"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/file"
	"github.com/yumenaka/comigo/tools/logger"
)

// analyzeImage 获取某页漫画的分辨率与blurhash
func (i *PageInfo) analyzeImage(bookPath string) (err error) {
	var img image.Image
	imgData, err := file.GetSingleFile(bookPath, i.Name, "gbk")
	if err != nil {
		logger.Infof("%s", err)
	}
	buf := bytes.NewBuffer(imgData)
	img, err = imaging.Decode(buf)
	if err != nil {
		logger.Infof(locale.GetString("check_image_error")+" %v\n", err)
	} else {
		i.Width = img.Bounds().Dx()
		i.Height = img.Bounds().Dy()
		// 很耗费服务器资源，以后再研究。
		str, err := blurhash.Encode(1, 1, img)
		if err != nil {
			// Handle errors
			logger.Infof(locale.GetString("check_image_error")+" %v\n", err)
		}
		i.Blurhash = str
	}
	return err
}

// analyzePageImages 解析漫画的分辨率与类型
func analyzePageImages(p *PageInfo, bookPath string) {
	err := p.analyzeImage(bookPath)
	if err != nil {
		logger.Infof(locale.GetString("check_image_error") + err.Error())
		return
	}
	if p.Width == 0 && p.Height == 0 {
		p.ImgType = "Unknown"
		return
	}
	if p.Width > p.Height {
		p.ImgType = "DoublePage"
	} else {
		p.ImgType = "SinglePage"
	}
}

// ScanAllImage 服务器端分析分辨率、漫画单双页，只适合已解压文件
func (b *Book) ScanAllImage() {
	logger.Infof(locale.GetString("check_image_start"))
	bar := pb.StartNew(b.GetPageCount())
	for i := range b.PageInfos {
		analyzePageImages(&b.PageInfos[i], b.BookPath)
		bar.Increment()
	}
	bar.Finish()
	logger.Infof(locale.GetString("check_image_completed"))
}

// ScanAllImageGo 并发分析图片
func (b *Book) ScanAllImageGo() {
	logger.Infof(locale.GetString("check_image_start"))
	wp := workpool.New(10) // 设置最大线程数
	bar := pb.StartNew(b.GetPageCount())

	for i := range b.PageInfos {
		i := i // 避免闭包问题
		wp.Do(func() error {
			analyzePageImages(&b.PageInfos[i], b.BookPath)
			bar.Increment()
			return nil
		})
	}
	_ = wp.Wait()
	bar.Finish()
	logger.Infof(locale.GetString("check_image_completed"))
}
