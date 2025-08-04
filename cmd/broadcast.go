package cmd

import (
	"sync"

	"github.com/spf13/viper"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers"
	"github.com/yumenaka/comigo/routers/upload_api"
	"github.com/yumenaka/comigo/templ/pages/settings"
	"github.com/yumenaka/comigo/util/logger"
	"github.com/yumenaka/comigo/util/scan"
)

// ---------------------------------------------------------------------------
// SystemBroadcast 相关逻辑
// ---------------------------------------------------------------------------

// broadcastRegistry 保存所有注册的系统广播通道，key 为自定义名称。
// 通过 map 的方式可以在不同包中动态注册并复用所需的广播通道。
var (
	broadcastRegistry = make(map[string]chan string)
	broadcastMutex    sync.RWMutex

	// SystemBroadcast 等于 RegisterBroadcast("system") 返回的通道。
	SystemBroadcast chan string
)

// RegisterBroadcast 创建或返回一个已存在的广播通道。
// 不同模块可以调用此方法注册/获取自己的通道，实现解耦。
func RegisterBroadcast(name string) chan string {
	broadcastMutex.Lock()
	defer broadcastMutex.Unlock()

	if ch, ok := broadcastRegistry[name]; ok {
		return ch
	}
	ch := make(chan string)
	broadcastRegistry[name] = ch
	return ch
}

// GetBroadcast 返回指定名称的广播通道（只读锁）。
func GetBroadcast(name string) (chan string, bool) {
	broadcastMutex.RLock()
	defer broadcastMutex.RUnlock()
	ch, ok := broadcastRegistry[name]
	return ch, ok
}

func init() {
	// 初始化默认的广播通道
	SystemBroadcast = RegisterBroadcast("uploadAPIBroadcast")
	// 在上传模块注册重新扫描通道
	upload_api.RescanBroadcast = &SystemBroadcast
	settings.RestartWebServerBroadcast = &SystemBroadcast
	// 在默认频道上开始收听传入的聊天消息
	go waitSystemMessages()
}

// 一个简单循环，从"SystemBroadcast"中连续读取数据，然后通过各自的 WebSocket 连接将消息传播到客户端。
func waitSystemMessages() {
	for {
		// 使用消息通道做各种事
		msg := <-SystemBroadcast // 广播频道
		// Send it out to every client that is currently connected
		switch msg {
		case "rescan_upload_path":
			logger.Infof("重新扫描上传文件夹：%s", msg)
			ReScanUploadPath()
			// 保存扫描结果到数据库
			if config.GetEnableDatabase() {
				err := scan.SaveResultsToDatabase(viper.ConfigFileUsed(), config.GetClearDatabaseWhenExit())
				if err != nil {
					return
				}
			}
		case "restart_web_server":
			logger.Infof("Config changed, restarting web server...\n", msg)
			routers.RestartWebServer()
		case "rescan_path_sample":
			logger.Infof("收到重新扫描消息：%s", msg)
			ReScanPath(msg, false)
		default:
			continue
		}
	}
}

// ReScanUploadPath 重新扫描上传目录,因为需要设置下载路径，gin 初始化后才能执行
func ReScanUploadPath() {
	// 没启用上传，则不扫描
	if !config.GetEnableUpload() {
		return
	}
	ReScanPath(config.GetUploadPath(), true)
}

func ReScanPath(path string, reScanFile bool) {
	// 扫描上传目录的文件
	option := scan.NewOption(config.GetCfg())
	books, err := scan.InitStore(path, option)
	if err != nil {
		logger.Infof(locale.GetString("scan_error")+"path:%s  %s", path, err)
		return
	}
	scan.AddBooksToStore(books, path, config.GetMinImageNum())
	model.MainStores.ResetBookGroupData()
}
