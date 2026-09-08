package cmd

import (
	"bytes"
	"encoding/json"
	"testing"
)

// 桌面查询只输出机器协议，不触发服务初始化或书库扫描。
func TestDesktopInfo(t *testing.T) {
	var out bytes.Buffer
	handled, err := RunDesktop([]string{"desktop", "info"}, &out)
	if !handled || err != nil {
		t.Fatal(handled, err)
	}
	var info struct {
		Version  string `json:"version"`
		Protocol int    `json:"desktopProtocol"`
	}
	if err = json.Unmarshal(out.Bytes(), &info); err != nil || info.Version == "" || info.Protocol != 1 {
		t.Fatal(info, err)
	}
	if handled, _ := RunDesktop([]string{"book-folder"}, &out); handled {
		t.Fatal("ordinary CLI args intercepted")
	}
	for _, args := range [][]string{{"desktop"}, {"desktop", "unknown"}, {"desktop", "info", "extra"}} {
		if _, err := RunDesktop(args, &out); err == nil {
			t.Fatalf("accepted %v", args)
		}
	}
}
