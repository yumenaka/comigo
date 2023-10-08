package routers

import (
	"fmt"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/util"
	"math/rand"
	"strconv"
	"time"
)

// CheckWebPort 3、选择服务端口
func CheckWebPort() {
	//检测端口
	if !util.CheckPort(config.Config.Port) {
		//获取一个空闲可用的系统端口号
		port, err := util.GetFreePort()
		if err != nil {
			fmt.Println(err)
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			if config.Config.Port+2000 > 65535 {
				config.Config.Port = config.Config.Port + r.Intn(1024)
			} else {
				config.Config.Port = 30000 + r.Intn(20000)
			}
		} else {
			config.Config.Port = port
		}
		fmt.Println(locale.GetString("port_busy") + strconv.Itoa(config.Config.Port))
	}
}
