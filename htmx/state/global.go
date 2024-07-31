package state

import (
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/entity"
)

type GlobalState struct {
	Version         string
	SingleUserMode  bool
	NowBookID       string
	OnlineUserCount int
	BooksList       *entity.BookInfoList
}

// GetAllBookNum 获取所有书籍数量
func (g *GlobalState) GetAllBookNum() int {
	if g.BooksList == nil {
		return 0
	}
	return len(g.BooksList.BookInfos)
}

var Global GlobalState

func init() {
	Global.Version = config.Version
	Global.SingleUserMode = false
	Global.NowBookID = ""
	Global.OnlineUserCount = 0
	Global.BooksList = nil
}
