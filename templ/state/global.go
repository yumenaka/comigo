package state

import (
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util"
)

type GlobalState struct {
	Debug           bool
	Version         string
	StaticFileMode  bool
	OnlineUserCount int
	ShelfBookList   *model.BookInfoList
	ServerStatus    *util.ServerStatus
}

var ServerConfig *config.Config

// GetAllBookNum 获取所有书籍数量
func (g *GlobalState) GetAllBookNum() int {
	if g.ShelfBookList == nil {
		return 0
	}
	return len(g.ShelfBookList.BookInfos)
}

var Global GlobalState

func init() {
	Global.Debug = config.GetDebug()
	Global.Version = config.GetVersion()
	// 是否开启静态模式，开启Debug模式时，静态模式会被强制开启
	// 需要避免 </script>或 </body> 提前截断script标签的问题
	Global.StaticFileMode = config.GetStaticFileMode()
	Global.OnlineUserCount = 0
	Global.ShelfBookList = nil
	Global.ServerStatus = util.GetServerInfo(config.GetHost(), config.GetVersion(), config.GetPort(), config.GetEnableUpload(), 0)
	ServerConfig = config.GetCfg()
}

func (s *GlobalState) GetStaticFileMode() bool {
	return config.GetStaticFileMode()
}
