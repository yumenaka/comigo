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

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
)

// processState 只保存本机管理地址和随机凭据；不通过 PID 发送系统信号，避免误杀复用 PID 的进程。
type processState struct {
	PID        int
	Address    string
	Token      string
	Port       int // 实际启动端口，重载后保持稳定，用于实例选择。
	ConfigFile string
	URL        string `json:",omitempty"` // 仅状态查询返回当前阅读地址。
}

// processControl 的生命周期与 CLI 相同；Web 服务换端口或重启不会中断这个独立通道。
type processControl struct {
	server *http.Server
	path   string
	ready  atomic.Bool
	mu     sync.Mutex // 串行处理 reload 与 stop，避免两者同时改动服务。
}

var cliControl *processControl
var processStop = make(chan struct{})
var processStopOnce sync.Once

// processStatePath 定位当前配置的管理记录；跨目录发现由 localProcesses 补充。
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

// StartProcessControl 为每个进程建立独立管理通道，不限制其他实例。
// 仅 CLI 入口调用，嵌入式、Wails 和托盘继续由各自宿主管理生命周期。
func StartProcessControl() error {
	statePath, err := processStatePath()
	if err != nil {
		return err
	}
	if cliControl != nil {
		return fmt.Errorf("%s", locale.GetString("cli_control_initialized"))
	}
	statePath = processInstancePath(statePath, os.Getpid())
	control := &processControl{path: statePath}
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return err
	}
	file := config.GetCfg().ConfigFile
	if file != "" {
		file, _ = filepath.Abs(file)
	}
	state := processState{PID: os.Getpid(), Address: listener.Addr().String(), Token: rand.Text(), Port: config.GetCfg().Port, ConfigFile: file}
	control.server = &http.Server{
		Handler: control.handler(state.Token), ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout: 5 * time.Second,
	}
	// 先删除旧文件再用 0600 创建，避免继承旧文件较宽的权限。
	if err = os.Remove(statePath); err != nil && !os.IsNotExist(err) {
		_ = listener.Close()
		return err
	}
	data, _ := json.Marshal(state)
	if err := os.WriteFile(statePath, data, 0o600); err != nil {
		_ = listener.Close()
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

// CloseProcessControl 只删除当前进程的状态文件，不影响其他实例。
func CloseProcessControl() {
	if cliControl != nil {
		_ = cliControl.server.Close()
		_ = os.Remove(cliControl.path)
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
		case "/status":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(processState{PID: os.Getpid(), Port: config.GetCfg().Port, URL: config.GetLocalBrowserURL()})
			return
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
func processRequest(state processState, action, file string, result ...any) error {
	host, _, err := net.SplitHostPort(state.Address)
	if err != nil || host != "127.0.0.1" {
		return fmt.Errorf("%s", locale.GetString("cli_control_invalid_address"))
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
	if action == "status" || action == "ready" {
		client.Timeout = 3 * time.Second
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if action == "status" && response.StatusCode == http.StatusOK && len(result) == 1 {
		return json.NewDecoder(io.LimitReader(response.Body, 64*1024)).Decode(result[0])
	}
	if response.StatusCode != http.StatusNoContent {
		message, _ := io.ReadAll(io.LimitReader(response.Body, 64*1024))
		return fmt.Errorf("%s: %s", response.Status, strings.TrimSpace(string(message)))
	}
	return nil
}

// controlProcess 重载当前配置目录中唯一匹配的实例。
func controlProcess(action string) error {
	statePath, err := processStatePath()
	if err != nil {
		return err
	}
	_, state, err := findProcessState(statePath)
	if err != nil {
		return err
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
	return nil
}

// processInstancePath 用 PID 区分同一配置目录内的多个服务。
func processInstancePath(base string, pid int) string {
	return strings.TrimSuffix(base, ".json") + fmt.Sprintf("-%d.json", pid)
}

// findProcessState 只选择仍可连接的实例；存在多个目标时要求指定启动端口。
func findProcessState(base string) (string, processState, error) {
	paths, err := filepath.Glob(strings.TrimSuffix(base, ".json") + "-*.json")
	if err != nil {
		return "", processState{}, err
	}
	file := config.GetCfg().ConfigFile
	if file != "" {
		file, _ = filepath.Abs(file)
	}
	var chosen processState
	chosenPath := ""
	for _, path := range paths {
		state, err := readProcessState(path)
		if err != nil || (file != "" && file != state.ConfigFile) {
			continue
		}
		if RootCmd.PersistentFlags().Changed("port") && state.Port != config.GetCfg().Port {
			continue
		}
		if processRequest(state, "ready", "") != nil {
			continue
		}
		if chosenPath != "" {
			return "", processState{}, fmt.Errorf("%s", locale.GetString("cli_multiple_processes"))
		}
		chosenPath, chosen = path, state
	}
	if chosenPath == "" {
		return "", processState{}, fmt.Errorf("%s", locale.GetString("cli_not_running"))
	}
	return chosenPath, chosen, nil
}
