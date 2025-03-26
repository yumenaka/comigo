package scan

import "github.com/yumenaka/comigo/model"

// 处理视频、音频等媒体文件
func handleMediaFiles(newBook *model.Book, imageName, imageUrl string) {
	newBook.PageCount = 1
	newBook.InitComplete = true
	newBook.SetCover(model.MediaFileInfo{Name: imageName, Url: imageUrl})
}
