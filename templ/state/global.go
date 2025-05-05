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
	Global.Version = config.GetVersion()
	Global.ShelfBookList = nil
	Global.ServerStatus = util.GetServerInfo(config.GetHost(), config.GetVersion(), config.GetPort(), config.GetEnableUpload(), 0)
	ServerConfig = config.GetCfg()
}

// GetStaticFileMode 是否开启静态模式，开启Debug模式时，静态模式会被强制开启
// 需要避免 </script>或 </body> 提前截断script标签的问题
func (s *GlobalState) GetStaticFileMode() bool {
	// logger.Infof("GetStaticFileMode: %v", config.GetStaticFileMode())
	return config.GetStaticFileMode()
}
