package data_api

import (
	"testing"
)

// 验证同设备多标签、同用户多设备、匿名连接分别按约定去重。
func TestConnectionSnapshot(t *testing.T) {
	online.Lock()
	old := online.connections
	online.connections = map[string]Connection{
		"a": {Device: "browser-a", User: "admin"},
		"b": {Device: "browser-a", User: "admin"},
		"c": {Device: "browser-b", User: "admin"},
		"d": {Device: "guest"},
	}
	online.Unlock()
	defer func() { online.Lock(); online.connections = old; online.Unlock() }()
	connections, users, devices := connectionSnapshot()
	if len(connections) != 4 || users != 1 || devices != 3 {
		t.Fatalf("connections=%d users=%d devices=%d", len(connections), users, devices)
	}
}
