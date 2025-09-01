package tailscale

import (
	"github.com/yumenaka/comigo/util/logger"
	"tailscale.com/tsnet"
)

// Run https://tailscale.com/kb/1521/hello-tsnet
func RunTailscale(hostname string, port uint16) {
	srv := new(tsnet.Server)
	srv.Hostname = hostname
	srv.Port = port
	if err := srv.Start(); err != nil {
		logger.Fatalf("can't start tsnet server: %v", err)
	}
	defer srv.Close()
	// 使用tsnet.server对象与tailscale互动...
}
