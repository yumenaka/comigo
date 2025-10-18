package tools

import (
	"net"
	"strconv"
	"time"

	"github.com/yumenaka/comigo/tools/logger"
)

// WaitUntilServerReady 循环尝试与端口建立 TCP 连接，直到成功或超时
func WaitUntilServerReady(host string, port uint16, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	addr := host + ":" + strconv.Itoa(int(port))
	for {
		if time.Now().After(deadline) {
			logger.Infof("Server not ready within %v, continue anyway", timeout)
			return
		}
		conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
		if err == nil {
			_ = conn.Close()
			return
		}
		time.Sleep(300 * time.Millisecond)
	}
}
