package routers

import (
	"fmt"

	"github.com/yumenaka/comi/book"
	"github.com/yumenaka/comi/common"
)

//用来防止重复注册的URL表，key是bookID，值是StaticURL
var staticUrlMap = make(map[string]string)

func checkUrlRegistered(bookID string) bool {
	_, ok := staticUrlMap[bookID]
	return ok
}

//SetDownloadLink 重新设定压缩包下载链接
func SetDownloadLink() {
	if book.GetBooksNumber() >= 1 {
		allBook, err := book.GetAllBookInfoList("name")
		if err != nil {
			fmt.Println("设置文件下载失败")
		} else {
			for _, info := range allBook.BookInfos {
				//下载文件
				if info.Type != book.TypeBooksGroup && info.Type != book.TypeDir {
					//staticUrl := "/raw/" + info.BookID + "/" + url.QueryEscape(info.Name)
					staticUrl := "/raw/" + info.BookID + "/" + info.Name
					if checkUrlRegistered(info.BookID) {
						if common.Config.Debug {
							fmt.Println("static 注册过了：" + info)
						}
						continue
					} else {
						api.StaticFile(staticUrl, info.FilePath)
						staticUrlMap[info.BookID] = staticUrl
					}
				}
			}
		}
	}
}
