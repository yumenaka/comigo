//go:build !js

package model

import (
	"github.com/cheggaaa/pb/v3"
	"github.com/xxjwxc/gowp/workpool"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/util/logger"
)

// ScanAllImage 服务器端分析分辨率、漫画单双页，只适合已解压文件
func (b *Book) ScanAllImage() {
	logger.Infof(locale.GetString("check_image_start"))
	bar := pb.StartNew(b.GetPageCount())
	for i := range b.Pages.Images {
		analyzePageImages(&b.Pages.Images[i], b.FilePath)
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

	for i := range b.Pages.Images {
		i := i // 避免闭包问题
		wp.Do(func() error {
			analyzePageImages(&b.Pages.Images[i], b.FilePath)
			bar.Increment()
			return nil
		})
	}
	_ = wp.Wait()
	bar.Finish()
	logger.Infof(locale.GetString("check_image_completed"))
}
