package state

import (
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools"
)

var (
	Version        string
	ServerStatus   *tools.ServerStatus
	ServerConfig   *config.Config
	StoreBookInfos []model.StoreBookInfo
	ChildBookInfos []model.BookInfo
)

// GetNowBookNum 获取当前显示书籍数量
func GetNowBookNum() int {
	if StoreBookInfos == nil {
		return 0
	}
	return len(StoreBookInfos)
}

// 初始化参数
func init() {
	Version = config.GetVersion()
	StoreBookInfos = nil
	//ServerStatus = tools.GetServerInfo(config.GetHost(), config.GetVersion(), uint16(config.GetCfg().Port), config.GetEnableUpload(), 0, "")
	ServerConfig = config.GetCfg()
}
