package scan

import "github.com/yumenaka/comigo/model"

// 处理视频、音频等媒体文件
func handleMediaFiles(newBook *model.Book) {
	newBook.PageCount = 1
	newBook.InitComplete = true
}
