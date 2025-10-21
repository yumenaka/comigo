package routers

import (
	"math/rand"
	"time"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
)

// SetHttpPort 设置服务端口
func SetHttpPort() {
	// 检测端口是否可用
	if !tools.CheckPort(uint16(config.GetCfg().Port)) {
		// 端口被占用
		logger.Infof(locale.GetString("port_busy"), config.GetCfg().Port)
		// 获取一个空闲的系统端口号
		port, err := tools.GetFreePort()
		if err != nil {
			logger.Infof("Failed to get a free port: %v", err)
			// 如果无法获取空闲端口，则随机选择一个端口
			rand.New(rand.NewSource(time.Now().UnixNano()))
			if config.GetCfg().Port+2000 > 65535 {
				config.SetPort(rand.Intn(1024) + config.GetCfg().Port)
			} else {
				config.SetPort(rand.Intn(20000) + 30000)
			}
		} else {
			config.SetPort(port)
		}
		logger.Infof("Using port: %d", config.GetCfg().Port)
	}
}
