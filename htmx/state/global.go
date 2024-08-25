package state

import (
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/entity"
	"github.com/yumenaka/comigo/util"
	"strings"
)

type GlobalState struct {
	Version         string
	SingleUserMode  bool
	RequestBookID   string
	OnlineUserCount int
	TopBooks        *entity.BookInfoList
	ServerStatus    *util.ServerStatus
}

// GetAllBookNum 获取所有书籍数量
func (g *GlobalState) GetAllBookNum() int {
	if g.TopBooks == nil {
		return 0
	}
	return len(g.TopBooks.BookInfos)
}

var Global GlobalState

func init() {
	Global.Version = config.Version
	Global.SingleUserMode = false
	Global.RequestBookID = ""
	Global.OnlineUserCount = 0
	Global.TopBooks = nil
	Global.ServerStatus = util.GetServerInfo(config.Config.Host, config.Version, config.Config.Port, config.Config.EnableUpload, 0)
}

func GetCloverBackgroundImageUrl(book *entity.BookInfo) string {
	imageUrl := book.Cover.Url
	if strings.HasPrefix(book.Cover.Url, "/api") {
		imageUrl = book.Cover.Url + "&resize_width=256&resize_height=360&thumbnail_mode=true"
	}
	return imageUrl
}
