package cmd

import (
	"sync"
	"time"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers"
	"github.com/yumenaka/comigo/routers/upload_api"
	"github.com/yumenaka/comigo/templ/pages/settings"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/scan"
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
		//// 重新扫描上传目录
		//case "rescan_upload_path":
		//	logger.Infof("Rescan the upload folder：%s", msg)
		//	if config.GetEnableUpload() {
		//		ReScanPath(config.GetUploadPath(), true)
		//	}
		//	// 保存扫描结果到数据库
		//	if config.GetEnableDatabase() {
		//		err := scan.SaveResultsToDatabase(viper.ConfigFileUsed(), config.GetClearDatabaseWhenExit())
		//		if err != nil {
		//			return
		//		}
		//	}
		// 重启网页服务器
		case "restart_web_server":
			logger.Infof("Config changed, restarting web server...\n", msg)
			routers.RestartWebServer()
			// 阻塞等待端口就绪，确保服务可用
			tools.WaitUntilServerReady("localhost", config.GetPort(), 15*time.Second)
			// 在命令行显示QRCode
			ShowQRCode()
			// 重启网页服务器
		case "start_tailscale":
			logger.Infof("Config changed, starting tailscale...\n", msg)
			routers.StartTailscale()
		case "stop_tailscale":
			logger.Infof("Config changed, stopping tailscale...\n", msg)
			routers.StopTailscale()
		case "restart_tailscale":
			logger.Infof("Config changed, restart tailscale...\n", msg)
			routers.StopTailscale()
			routers.StartTailscale()
		// 重新扫描指定目录
		case "rescan_path_sample":
			logger.Infof("收到重新扫描消息：%s", msg)
			ReScanPath(msg, false)
		default:
			continue
		}
	}
}

// ReScanPath  重新扫描目录,因为需要设置下载路径，gin 初始化后才能执行
func ReScanPath(storeUrl string, reScanFile bool) {
	// 扫描上传目录的文件
	option := scan.NewOption(config.GetCfg())
	books, err := scan.InitStore(storeUrl, option)
	if err != nil {
		logger.Infof(locale.GetString("scan_error")+"path:%s  %s", storeUrl, err)
		return
	}
	scan.AddBooksToStore(storeUrl, books, config.GetMinImageNum())
}
