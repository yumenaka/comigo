package routers

import (
	"math/rand"
	"time"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/locale"
	"github.com/yumenaka/comigo/util/logger"
)

// SetWebServerPort 设置服务端口
func SetWebServerPort() {
	// 检测端口是否可用
	if !util.CheckPort(config.Config.Port) {
		// 端口被占用
		logger.Infof(locale.GetString("port_busy"), config.Config.Port)
		// 获取一个空闲的系统端口号
		port, err := util.GetFreePort()
		if err != nil {
			logger.Infof("Failed to get a free port: %v", err)
			// 如果无法获取空闲端口，则随机选择一个端口
			rand.New(rand.NewSource(time.Now().UnixNano()))
			if config.Config.Port+2000 > 65535 {
				config.Config.Port = rand.Intn(1024) + config.Config.Port
			} else {
				config.Config.Port = rand.Intn(20000) + 30000
			}
		} else {
			config.Config.Port = port
		}
		logger.Infof("Using port: %d", config.Config.Port)
	}
}
