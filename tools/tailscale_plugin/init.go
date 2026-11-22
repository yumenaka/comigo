//go:build !wails

package tailscale_plugin

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"tailscale.com/client/local"
	"tailscale.com/tsnet"
)

type TailscaleConfig struct {
	Hostname   string
	Port       uint16
	FunnelMode bool
	ConfigDir  string
	AuthKey    string
}

// listenTailscale 创建本实例的监听；Funnel 可能等待用户登录。
func listenTailscale(ctx context.Context, server *tsnet.Server, client *local.Client, c TailscaleConfig) (net.Listener, error) {
	addr := fmt.Sprintf(":%d", c.Port)
	if c.FunnelMode {
		// ListenFunnel 内部用 Background 等待登录，先让可取消的 Up 完成这一步。
		if _, err := server.Up(ctx); err != nil {
			return nil, err
		}
		return server.ListenFunnel("tcp", addr)
	}
	ln, err := server.Listen("tcp", addr)
	if err == nil && c.Port == 443 {
		ln = tls.NewListener(ln, &tls.Config{GetCertificate: client.GetCertificate})
	}
	return ln, err
}
