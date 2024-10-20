package state

import (
	"strings"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/entity"
	"github.com/yumenaka/comigo/util"
)

type GlobalState struct {
	Debug           bool
	Version         string
	SingleUserMode  bool
	StaticMode      bool
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
	Global.Debug = config.Config.Debug
	Global.Version = config.Version
	Global.SingleUserMode = false
	// 是否为静态模式 默认为false
	// 静态模式下，7777 端口的反代服务器无法正常加载静态JS资源，导致页面无法正常显示。
	// 反代出错的原因不明，暂时不管了。调试静态模式的时候看1234就好。
	Global.StaticMode = false
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
