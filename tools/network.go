package tools

import (
	"net"
	"regexp"
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

var domainRegex = regexp.MustCompile(`^(?i)[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)*$`)

func IsValidDomain(host string) bool {
	if len(host) == 0 || len(host) > 253 {
		return false
	}
	return domainRegex.MatchString(host)
}
