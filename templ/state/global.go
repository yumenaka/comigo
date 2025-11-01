package state

import (
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools"
)

var (
	Version      string
	ServerStatus *tools.ServerStatus
	ServerConfig *config.Config
	NowBookInfos *model.BookInfos
)

// GetNowBookNum 获取当前显示书籍数量
func GetNowBookNum() int {
	if NowBookInfos == nil {
		return 0
	}
	return len(*NowBookInfos)
}

// 初始化参数
func init() {
	Version = config.GetVersion()
	NowBookInfos = nil
	//ServerStatus = tools.GetServerInfo(config.GetHost(), config.GetVersion(), uint16(config.GetCfg().Port), config.GetEnableUpload(), 0, "")
	ServerConfig = config.GetCfg()
}
