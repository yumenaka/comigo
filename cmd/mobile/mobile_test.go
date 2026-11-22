package mobile

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers"
)

// unavailableHealthTransport 仅在测试中模拟服务未就绪，不连接外部服务器。
type unavailableHealthTransport struct{}

// RoundTrip 返回就绪失败，覆盖实际 Start 的回滚路径。
func (unavailableHealthTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusServiceUnavailable, Body: io.NopCloser(strings.NewReader("not ready")), Header: make(http.Header)}, nil
}

// TestStartupSchedulesOnlyAfterReady 验证就绪失败不留调度，重新启动成功后仍可正常扫描和停止。
func TestStartupSchedulesOnlyAfterReady(t *testing.T) {
	oldCfg, oldStore, oldScanner, oldInfo, oldTransport := config.CopyCfg(), model.IStore, config.GlobalLibraryScanner, lastStartInfo, http.DefaultTransport
	t.Setenv("COMIGO_CONFIG_DIR", os.Getenv("COMIGO_CONFIG_DIR"))
	t.Cleanup(func() {
		if config.GlobalLibraryScanner != nil {
			_ = config.GlobalLibraryScanner.Stop()
		}
		_ = routers.StopWebServer()
		config.ReplaceConfig(oldCfg)
		model.IStore = oldStore
		config.GlobalLibraryScanner = oldScanner
		lastStartInfo = oldInfo
		http.DefaultTransport = oldTransport
	})
	config.GlobalLibraryScanner = nil
	lastStartInfo = StartResult{}
	startCfg := StartConfig{ConfigDir: t.TempDir(), CacheDir: t.TempDir(), Port: 0, DisableLAN: true, Language: "en", AutoRescanIntervalMinutes: 1, WaitReadyTimeoutSeconds: 1}
	input, err := json.Marshal(startCfg)
	if err != nil {
		t.Fatal(err)
	}
	http.DefaultTransport = unavailableHealthTransport{}
	var result StartResult
	if err := json.Unmarshal([]byte(Start(string(input))), &result); err != nil {
		t.Fatal(err)
	}
	if result.Success || result.Running || !strings.Contains(result.Error, "not ready") {
		t.Fatalf("未触发就绪失败: %+v", result)
	}
	if config.GlobalLibraryScanner != nil && config.GlobalLibraryScanner.GetInterval() != 0 {
		t.Fatal("失败启动遗留自动扫描")
	}
	http.DefaultTransport = oldTransport
	if err := json.Unmarshal([]byte(Start(string(input))), &result); err != nil {
		t.Fatal(err)
	}
	if !result.Success || !result.Running || config.GlobalLibraryScanner.GetInterval() != 1 {
		t.Fatalf("重新启动未创建调度: %+v", result)
	}
	if err := json.Unmarshal([]byte(Stop()), &result); err != nil {
		t.Fatal(err)
	}
	if !result.Success || result.Running || config.GlobalLibraryScanner.GetInterval() != 0 {
		t.Fatalf("停止未清理调度: %+v", result)
	}
}
