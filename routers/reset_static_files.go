package routers

import (
	"fmt"
	"github.com/yumenaka/comi/book"
)

//ResetStaticFiles 重新设定压缩包下载链接
func ResetStaticFiles() {
	if book.GetBooksNumber() >= 1 {
		allBook, err := book.GetAllBookInfoList("name")
		if err != nil {
			fmt.Println("设置文件下载失败")
		} else {
			for _, info := range allBook.BookInfos {
				//下载文件
				if info.Type != book.TypeBooksGroup && info.Type != book.TypeDir {
					api.StaticFile("/raw/"+info.BookID+"/"+info.Name, info.FilePath)
				}
			}
		}
	}
}
