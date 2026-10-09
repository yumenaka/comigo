package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
)

// 验证停止请求的实际目标、部分失败不阻断其余目标，以及三语结果包含 PID 和实际端口。
func TestStopProcesses(t *testing.T) {
	t.Cleanup(func() { locale.InitLanguageFromConfig("auto") })
	for _, lang := range []string{"zh", "en", "ja"} {
		for _, tc := range []struct {
			name, input       string
			count, fail, port int
			cancel            bool
			want              []int
		}{
			{name: "none", fail: -1},
			{name: "single", count: 1, fail: -1, want: []int{1}},
			{name: "choose", input: "9\n2\n", count: 2, fail: -1, want: []int{0, 1}},
			{name: "port filter", count: 2, fail: -1, port: 3001, want: []int{0, 1}},
			{name: "all", input: "0\n", count: 2, fail: -1, want: []int{1, 1}},
			{name: "partial failure", input: "0\n", count: 2, fail: 0, want: []int{1, 1}},
			{name: "cancel", count: 2, fail: -1, cancel: true, want: []int{0, 0}},
			{name: "end of input", count: 2, fail: -1, want: []int{0, 0}},
		} {
			t.Run(lang+"/"+tc.name, func(t *testing.T) {
				root := processTestCommand(t)
				locale.InitLanguageFromConfig(lang)
				file := filepath.Join(t.TempDir(), "config.toml")
				config.GetCfg().ConfigFile = file
				base, err := processStatePath()
				if err != nil {
					t.Fatal(err)
				}
				calls := make([]atomic.Int32, tc.count)
				for i := range tc.count {
					state := processState{PID: 100000 + i, Port: 3000 + i, ConfigFile: file, Token: "private-token"}
					path := processInstancePath(base, state.PID)
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.Method != http.MethodPost || r.Header.Get("Authorization") != "Bearer private-token" {
							t.Error("停止请求缺少认证")
							w.WriteHeader(403)
							return
						}
						if r.URL.Path == "/status" {
							current := state
							current.Port = 4000 + i
							_ = json.NewEncoder(w).Encode(current)
							return
						}
						if r.URL.Path != "/stop" {
							t.Error(r.URL.Path)
							w.WriteHeader(404)
							return
						}
						calls[i].Add(1)
						if i == tc.fail {
							http.Error(w, "denied", 403)
							return
						}
						if err := os.Remove(path); err != nil {
							t.Error(err)
						}
						w.WriteHeader(http.StatusNoContent)
					}))
					t.Cleanup(server.Close)
					state.Address = strings.TrimPrefix(server.URL, "http://")
					data, _ := json.Marshal(state)
					if err := os.WriteFile(path, data, 0600); err != nil {
						t.Fatal(err)
					}
				}
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				if tc.cancel {
					cancel()
				}
				root.SetContext(ctx)
				var out bytes.Buffer
				root.SetOut(&out)
				root.SetIn(strings.NewReader(tc.input))
				args := []string{"stop", "--config", file}
				if tc.port != 0 {
					args = append(args, "--port", strconv.Itoa(tc.port))
				}
				root.SetArgs(args)
				err = root.Execute()
				if (err != nil) != (tc.fail >= 0) {
					t.Fatalf("err=%v output=%s", err, out.String())
				}
				result := out.String()
				if err != nil {
					result += err.Error()
				}
				if strings.Contains(result, "private-token") || strings.Contains(result, "%!") || strings.Contains(result, "cli_stop") {
					t.Fatal(result)
				}
				if tc.count == 0 && !strings.Contains(result, locale.GetString("cli_stop_none")) {
					t.Fatal(result)
				}
				if tc.count > 1 && tc.port == 0 && tc.input == "" && !strings.Contains(result, locale.GetString("cli_stop_cancelled")) {
					t.Fatal(result)
				}
				for i, want := range tc.want {
					if calls[i].Load() != int32(want) {
						t.Fatalf("目标 %d: calls=%d want=%d", i, calls[i].Load(), want)
					}
					if want == 0 {
						continue
					}
					key := "cli_stopped"
					var message string
					if i == tc.fail {
						message = fmt.Sprintf(locale.GetString("cli_stop_failed"), 100000+i, strconv.Itoa(4000+i), "403 Forbidden: denied")
					} else {
						message = fmt.Sprintf(locale.GetString(key), 100000+i, strconv.Itoa(4000+i))
					}
					if !strings.Contains(result, message) {
						t.Fatalf("未报告目标结果 %q: %s", message, result)
					}
				}
			})
		}
	}
}
