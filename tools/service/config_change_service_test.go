package service

import (
	"testing"
	"time"

	"github.com/yumenaka/comigo/config"
)

// 验证配置变化只触发对应的服务动作，未改动字段复用原配置。
func TestBuildConfigChangeAction(t *testing.T) {
	old := config.Config{Port: 1234, StoreUrls: []string{"/a"}, TailscaleHostname: "comigo", TailscalePort: 443}
	for _, tc := range []struct {
		name   string
		change func(*config.Config)
		want   ConfigChangeAction
	}{
		{"不变", func(c *config.Config) {}, ConfigChangeAction{}},
		{"多个配置", func(c *config.Config) {
			c.Port = 5678
			c.StoreUrls = []string{"/a", "/b"}
			c.EnableTailscale = true
			c.AutoRescanIntervalMinutes = 10
		}, ConfigChangeAction{ReScanStores: true, ReStartWebServer: true, StartTailscale: true, UpdateAutoRescan: true}},
		{"基础路径", func(c *config.Config) { c.BasePath = "/proxy" }, ConfigChangeAction{ReStartWebServer: true}},
		{"启用 Tailscale", func(c *config.Config) { c.EnableTailscale = true }, ConfigChangeAction{StartTailscale: true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			next := old
			tc.change(&next)
			if got := BuildConfigChangeAction(old, &next); got != tc.want {
				t.Fatalf("action=%+v, want %+v", got, tc.want)
			}
		})
	}
}

// 等价基础路径无需重启。
func TestBuildConfigChangeActionIgnoresEquivalentBasePath(t *testing.T) {
	old := config.Config{BasePath: "/proxy/"}
	next := old
	next.BasePath = "/proxy"
	if action := BuildConfigChangeAction(old, &next); action.ReStartWebServer {
		t.Fatal("等价基础路径触发了重启")
	}
}

// 已启用 Tailscale 时，关闭或修改连接参数分别触发停止或重启。
func TestBuildConfigChangeActionUpdatesTailscale(t *testing.T) {
	old := config.Config{EnableTailscale: true, TailscaleAuthKey: "old", TailscaleHostname: "comigo", TailscalePort: 443}
	for _, tc := range []struct {
		name   string
		change func(*config.Config)
		want   ConfigChangeAction
	}{
		{"关闭", func(c *config.Config) { c.EnableTailscale = false }, ConfigChangeAction{StopTailscale: true}},
		{"密钥", func(c *config.Config) { c.TailscaleAuthKey = "new" }, ConfigChangeAction{ReStartTailscale: true}},
		{"主机名", func(c *config.Config) { c.TailscaleHostname = "reader" }, ConfigChangeAction{ReStartTailscale: true}},
		{"端口", func(c *config.Config) { c.TailscalePort = 8443 }, ConfigChangeAction{ReStartTailscale: true}},
		{"Funnel", func(c *config.Config) { c.FunnelTunnel = true }, ConfigChangeAction{ReStartTailscale: true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			next := old
			tc.change(&next)
			if got := BuildConfigChangeAction(old, &next); got != tc.want {
				t.Fatalf("action=%+v, want %+v", got, tc.want)
			}
		})
	}
}

// 验证配置请求只提交重启信号，不等待新监听就绪，避免与 HTTP Shutdown 互相等待。
func TestApplyConfigChangeDoesNotWaitForRestart(t *testing.T) {
	old := config.Config{Username: "old"}
	next := config.Config{Username: "new"}
	signal := make(chan string, 1)
	done := make(chan struct{})
	go func() { ApplyConfigChange(old, &next, signal); close(done) }()
	select {
	case <-done:
		if got := <-signal; got != "restart_web_server" {
			t.Fatalf("unexpected signal: %s", got)
		}
	case <-time.After(time.Second):
		t.Fatal("配置更新不应等待重启完成")
	}
}

// 验证 TLS 配置变化会重建监听，而非只改变界面显示的协议。
func TestTLSConfigChangeRestartsWebServer(t *testing.T) {
	for _, changed := range []config.Config{
		{EnableTLS: true}, {AutoTLSCertificate: true}, {CertFile: "cert.pem"}, {KeyFile: "key.pem"},
	} {
		if !BuildConfigChangeAction(config.Config{}, &changed).ReStartWebServer {
			t.Fatal("TLS change did not restart server")
		}
	}
}

// ZIP 编码改变需要重新扫描，才能让网页修改影响已有书籍。
func TestZIPEncodingChangeRescansLibrary(t *testing.T) {
	old := config.Config{ZipFileTextEncoding: "gbk"}
	next := old
	next.ZipFileTextEncoding = "shiftjis"
	if action := BuildConfigChangeAction(old, &next); !action.ReScanStores || action.ReStartWebServer {
		t.Fatalf("action=%+v", action)
	}
}
