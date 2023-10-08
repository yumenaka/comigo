package routers

import (
	"fmt"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/util"
	"math/rand"
	"strconv"
	"time"
)

// CheckWebPort 3、选择服务端口
func CheckWebPort() {
	//检测端口
	if !util.CheckPort(common.Config.Port) {
		//获取一个空闲可用的系统端口号
		port, err := util.GetFreePort()
		if err != nil {
			fmt.Println(err)
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			if common.Config.Port+2000 > 65535 {
				common.Config.Port = common.Config.Port + r.Intn(1024)
			} else {
				common.Config.Port = 30000 + r.Intn(20000)
			}
		} else {
			common.Config.Port = port
		}
		fmt.Println(locale.GetString("port_busy") + strconv.Itoa(common.Config.Port))
	}
}
