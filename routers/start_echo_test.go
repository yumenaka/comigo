package routers

import (
	"slices"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"golang.org/x/crypto/acme"
)

// 验证自动证书模式同时保留常规 HTTPS 协议，避免浏览器握手失败。
func TestAutoTLSNextProtosIncludesHTTPSProtocols(t *testing.T) {
	t.Helper()

	got := autoTLSNextProtos()
	want := []string{acme.ALPNProto, "h2", "http/1.1"}

	for _, proto := range want {
		if !slices.Contains(got, proto) {
			t.Fatalf("autoTLSNextProtos() 缺少协议 %q，当前值: %v", proto, got)
		}
	}
}

// TestBuildHTTPServerSetsConnectionTimeouts 验证 LAN 服务不会无限等待慢速请求头。
func TestBuildHTTPServerSetsConnectionTimeouts(t *testing.T) {
	cfg := config.GetCfg()
	oldAutoTLS, oldCert, oldKey := cfg.AutoTLSCertificate, cfg.CertFile, cfg.KeyFile
	cfg.AutoTLSCertificate, cfg.CertFile, cfg.KeyFile = false, "", ""
	defer func() {
		cfg.AutoTLSCertificate, cfg.CertFile, cfg.KeyFile = oldAutoTLS, oldCert, oldKey
	}()
	server, _, err := buildHTTPServer(echo.New())
	if err != nil {
		t.Fatal(err)
	}
	if server.ReadHeaderTimeout != readHeaderTimeout || server.IdleTimeout != idleTimeout {
		t.Fatalf("server timeouts = %v/%v", server.ReadHeaderTimeout, server.IdleTimeout)
	}
}

// 对外服务切换必须改变实际监听地址。
func TestHTTPServerListenScope(t *testing.T) {
	old := config.CopyCfg()
	t.Cleanup(func() { *config.GetCfg() = old })
	config.GetCfg().AutoTLSCertificate = false
	config.GetCfg().CertFile = ""
	config.GetCfg().KeyFile = ""
	config.GetCfg().Port = 1234
	for _, tc := range []struct {
		local   bool
		address string
	}{{true, "127.0.0.1:1234"}, {false, "0.0.0.0:1234"}} {
		config.GetCfg().DisableLAN = tc.local
		server, _, err := buildHTTPServer(echo.New())
		if err != nil || server.Addr != tc.address {
			t.Fatalf("监听地址: %v, %v", server, err)
		}
	}
}
