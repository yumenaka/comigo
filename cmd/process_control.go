package cmd

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofrs/flock"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
)

// processState 只保存本机管理地址和随机凭据；不通过 PID 发送系统信号，避免误杀复用 PID 的进程。
type processState struct {
	PID     int
	Address string
	Token   string
}

// processControl 的生命周期与 CLI 相同；Web 服务换端口或重启不会中断这个独立通道。
type processControl struct {
	server *http.Server
	lock   *flock.Flock
	path   string
	ready  atomic.Bool
	mu     sync.Mutex // 串行处理 reload 与 stop，避免两者同时改动服务。
}

var cliControl *processControl
var processStop = make(chan struct{})
var processStopOnce sync.Once

// processStatePath 用配置目录定位实例，因此更改 TOML 中的端口后仍能 reload/stop。
// 自定义配置的管理命令应带相同的 --config，或使用相同的 COMIGO_CONFIG_DIR。
func processStatePath() (string, error) {
	if config.GetCfg().ConfigFile == "" {
		_, file := config.FindConfigFile()
		if file != "" {
			return filepath.Join(filepath.Dir(file), ".comi-process.json"), nil
		}
	}
	dir, err := config.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ".comi-process.json"), nil
}

// StartProcessControl 在启动 Web 服务之前占用实例锁；异常退出后的锁由操作系统自动释放。
// 仅 CLI 入口调用，嵌入式、Wails 和托盘继续由各自宿主管理生命周期。
func StartProcessControl() error {
	statePath, err := processStatePath()
	if err != nil {
		return err
	}
	control := &processControl{path: statePath, lock: flock.New(statePath + ".lock")}
	locked, err := control.lock.TryLock()
	if err != nil {
		return err
	}
	if !locked {
		return fmt.Errorf("%s", locale.GetString("cli_already_running"))
	}
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		_ = control.lock.Unlock()
		return err
	}
	state := processState{PID: os.Getpid(), Address: listener.Addr().String(), Token: rand.Text()}
	control.server = &http.Server{
		Handler: control.handler(state.Token), ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout: 5 * time.Second,
	}
	// 先删除旧文件再用 0600 创建，避免继承旧文件较宽的权限。
	if err = os.Remove(statePath); err != nil && !os.IsNotExist(err) {
		_ = listener.Close()
		_ = control.lock.Unlock()
		return err
	}
	data, _ := json.Marshal(state)
	if err := os.WriteFile(statePath, data, 0o600); err != nil {
		_ = listener.Close()
		_ = control.lock.Unlock()
		return err
	}
	cliControl = control
	go control.server.Serve(listener)
	return nil
}

// ProcessReady 只在 Web 监听和初次扫描完成后发布就绪，供 start 确认自己启动的进程。
func ProcessReady() {
	cliControl.ready.Store(true)
}

// CloseProcessControl 清理状态后再释放锁；避免旧实例退出时删除新实例的状态文件。
func CloseProcessControl() {
	if cliControl != nil {
		_ = cliControl.server.Close()
		_ = os.Remove(cliControl.path)
		_ = cliControl.lock.Unlock()
		cliControl = nil
	}
}

// handler 的所有请求都必须持有文件内的随机令牌，网页与其他本机用户不能直接控制进程。
// 控制请求不接受任意配置内容，只允许当前用户指定本机 TOML 文件。
func (control *processControl) handler(token string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost || subtle.ConstantTimeCompare([]byte(request.Header.Get("Authorization")), []byte("Bearer "+token)) != 1 {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		control.mu.Lock()
		defer control.mu.Unlock()
		if !control.ready.Load() {
			http.Error(w, "starting", http.StatusServiceUnavailable)
			return
		}
		switch request.URL.Path {
		case "/ready":
		case "/stop":
			// 先把成功响应发到客户端，再唤醒退出流程，避免清理连接导致客户端收到 EOF。
			w.WriteHeader(http.StatusNoContent)
			w.(http.Flusher).Flush()
			control.ready.Store(false)
			processStopOnce.Do(func() { close(processStop) })
			return
		case "/reload":
			var file string
			if err := json.NewDecoder(http.MaxBytesReader(w, request.Body, 64*1024)).Decode(&file); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err := reloadConfig(file); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		default:
			http.NotFound(w, request)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

// readProcessState 不输出凭据；残留文件不能证明服务存活，调用方还要实际连接验证。
func readProcessState(path string) (processState, error) {
	var state processState
	data, err := os.ReadFile(path)
	if err == nil {
		err = json.Unmarshal(data, &state)
	}
	return state, err
}

// processRequest 禁用代理和重定向，确保管理令牌只发到本机的既定地址。
func processRequest(state processState, action, file string) error {
	host, _, err := net.SplitHostPort(state.Address)
	if err != nil || host != "127.0.0.1" {
		return fmt.Errorf("invalid local control address")
	}
	body, _ := json.Marshal(file)
	request, err := http.NewRequest(http.MethodPost, "http://"+state.Address+"/"+action, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+state.Token)
	transport := &http.Transport{DialContext: (&net.Dialer{Timeout: time.Second}).DialContext}
	defer transport.CloseIdleConnections()
	client := &http.Client{Transport: transport, Timeout: 2 * time.Minute,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusNoContent {
		message, _ := io.ReadAll(io.LimitReader(response.Body, 64*1024))
		return fmt.Errorf("%s: %s", response.Status, strings.TrimSpace(string(message)))
	}
	return nil
}

// controlProcess 对当前配置目录中的实例执行操作；stop 等待清理结束后才返回成功。
func controlProcess(action string) error {
	statePath, err := processStatePath()
	if err != nil {
		return err
	}
	state, err := readProcessState(statePath)
	if err != nil {
		return fmt.Errorf("%s: %w", locale.GetString("cli_not_running"), err)
	}
	file := config.GetCfg().ConfigFile
	if file != "" {
		file, err = filepath.Abs(file)
		if err != nil {
			return err
		}
	}
	if err := processRequest(state, action, file); err != nil {
		return err
	}
	if action == "stop" {
		for deadline := time.Now().Add(15 * time.Second); time.Now().Before(deadline); time.Sleep(50 * time.Millisecond) {
			current, err := readProcessState(statePath)
			if os.IsNotExist(err) || (err == nil && current.Token != state.Token) {
				return nil
			}
		}
		return fmt.Errorf("%s", locale.GetString("cli_stop_timeout"))
	}
	return nil
}
