package routers

import (
	"fmt"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/plugin"
	"math/rand"
	"time"
)

// 5、setFrpClient
func setFrpClient() {
	//frp服务
	if !common.Config.EnableFrpcServer {
		return
	}
	if common.Config.FrpConfig.RemotePort <= 0 || common.Config.FrpConfig.RemotePort > 65535 {
		common.Config.FrpConfig.RemotePort = common.Config.Port
	}
	if common.Config.FrpConfig.RandomRemotePort {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		common.Config.FrpConfig.RemotePort = 50000 + r.Intn(10000)
	}
	frpcError := plugin.StartFrpC(common.Config.CachePath)
	if frpcError != nil {
		fmt.Println(locale.GetString("frpc_server_error"), frpcError.Error())
	} else {
		fmt.Println(locale.GetString("frpc_server_start"))
	}

}
