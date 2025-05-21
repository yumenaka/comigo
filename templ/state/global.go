package state

import (
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util"
)

// 感觉这个抽象有点多余？
// type GlobalState struct {
// 	Version      string
// 	NowBookList  *model.BookInfoList
// 	ServerStatus *util.ServerStatus
// }
// var Global GlobalState

var (
	Version      string
	ServerStatus *util.ServerStatus
	ServerConfig *config.Config
	NowBookList  *model.BookInfoList
)

// GetNowBookNum 获取当前显示书籍数量
func GetNowBookNum() int {
	if NowBookList == nil {
		return 0
	}
	return len(NowBookList.BookInfos)
}

// IsLogin 判断是否登录
func IsLogin() bool {
	return config.GetUsername() != "" && config.GetPassword() != ""
}

// 初始化参数
func init() {
	Version = config.GetVersion()
	NowBookList = nil
	ServerStatus = util.GetServerInfo(config.GetHost(), config.GetVersion(), config.GetPort(), config.GetEnableUpload(), 0)
	ServerConfig = config.GetCfg()
}
