package state

import (
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util"
)

type GlobalState struct {
	Version       string
	ShelfBookList *model.BookInfoList
	ServerStatus  *util.ServerStatus
}

var Global GlobalState

var ServerConfig *config.Config

// GetAllBookNum 获取所有书籍数量
func (g *GlobalState) GetAllBookNum() int {
	if g.ShelfBookList == nil {
		return 0
	}
	return len(g.ShelfBookList.BookInfos)
}

func IsLogin() bool {
	return config.GetUsername() != "" && config.GetPassword() != ""
}

func init() {
	Global.Version = config.GetVersion()
	Global.ShelfBookList = nil
	Global.ServerStatus = util.GetServerInfo(config.GetHost(), config.GetVersion(), config.GetPort(), config.GetEnableUpload(), 0)
	ServerConfig = config.GetCfg()
}
