package routers

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"golang.org/x/crypto/acme"
)

// 已预载证书的自定义 TLS 服务仍须协商 HTTP/2，不能退化为只支持 HTTP/1.1。
func TestCustomTLSServesHTTP2(t *testing.T) {
	fixture := httptest.NewTLSServer(http.NotFoundHandler())
	certificates := fixture.TLS.Certificates
	fixture.Close()
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := &http.Server{Handler: http.NotFoundHandler(), TLSConfig: &tls.Config{Certificates: certificates}}
	t.Cleanup(func() { _ = server.Close() })
	go serveHTTPServer(listener, server, webServeCustomTLS)
	transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, ForceAttemptHTTP2: true}
	t.Cleanup(transport.CloseIdleConnections)
	client := &http.Client{Transport: transport}
	response, err := client.Get("https://" + listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.ProtoMajor != 2 {
		t.Fatalf("协商协议=%s，期望 HTTP/2", response.Proto)
	}
}

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

// TLS 开关必须决定协议，坏证书与自动 TLS 前提错误必须同步返回。
func TestTLSStartupRejectsInvalidConfiguration(t *testing.T) {
	old := config.CopyCfg()
	t.Cleanup(func() { *config.GetCfg() = old })
	for _, tc := range []struct {
		name             string
		tls, auto, local bool
		cert, key, host  string
		wantErr          bool
	}{
		{name: "disabled", cert: "missing.crt", key: "missing.key"},
		{name: "no certificate", tls: true, wantErr: true},
		{name: "half certificate", tls: true, cert: "missing.crt", wantErr: true},
		{name: "bad certificate", tls: true, cert: "missing.crt", key: "missing.key", wantErr: true},
		{name: "auto without host", auto: true, wantErr: true},
		{name: "auto local", auto: true, local: true, host: "reader.example.com", wantErr: true},
		{name: "conflicting modes", tls: true, auto: true, host: "reader.example.com", wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg := config.GetCfg()
			cfg.EnableTLS, cfg.AutoTLSCertificate, cfg.DisableLAN = tc.tls, tc.auto, tc.local
			cfg.CertFile, cfg.KeyFile, cfg.Host = tc.cert, tc.key, tc.host
			server, mode, err := buildHTTPServer(echo.New())
			if (err != nil) != tc.wantErr {
				t.Fatalf("mode=%v err=%v", mode, err)
			}
			if !tc.wantErr && (mode != webServeHTTP || server.TLSConfig != nil) {
				t.Fatal("disabled TLS served HTTPS")
			}
		})
	}
}
