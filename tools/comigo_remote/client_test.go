package comigo_remote

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yumenaka/comigo/config"
)

// 验证远端版本读取优先使用服务信息中的版本字段。
func TestGetServerVersionUsesVersionField(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/server" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte(`{"Version":"v1.2.36"}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, 5)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	got, err := client.GetServerVersion()
	if err != nil {
		t.Fatalf("GetServerVersion: %v", err)
	}
	if got != "v1.2.36" {
		t.Fatalf("version = %q, want v1.2.36", got)
	}
}

// TestGetTopShelfRequiresSupportedVersion 旧版及无效版本在请求书架前被拒绝，不兼容旧版数据。
func TestGetTopShelfRequiresSupportedVersion(t *testing.T) {
	for _, version := range []string{"", "invalid", "v1.3.9", "v1.4.0-rc.1", config.GetMinSupportedVersion(), "v99.0.0"} {
		t.Run(version, func(t *testing.T) {
			shelfRequested := false
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/api/server":
					_ = json.NewEncoder(w).Encode(map[string]string{"Version": version})
				case "/api/top-shelf":
					shelfRequested = true
					_, _ = w.Write([]byte(`[]`))
				default:
					http.NotFound(w, r)
				}
			}))
			defer server.Close()
			client, err := NewClient(server.URL, 5)
			if err != nil {
				t.Fatal(err)
			}
			_, err = client.GetTopShelf("default")
			if config.SupportsBookVersion(version) {
				if err != nil || !shelfRequested {
					t.Fatalf("supported remote rejected: %v", err)
				}
			} else {
				var versionErr *UnsupportedVersionError
				if !errors.As(err, &versionErr) || shelfRequested {
					t.Fatalf("unsupported remote accessed: requested=%v err=%v", shelfRequested, err)
				}
			}
		})
	}
}
