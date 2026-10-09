package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
)

// processTestCommand 隔离 Cobra、Viper 和全局配置，避免测试之间遗留已解析的 flag。
func processTestCommand(t *testing.T) *cobra.Command {
	t.Helper()
	oldRoot, oldViper, oldConfig, oldArgs, oldDefaults := RootCmd, runtimeViper, config.CopyCfg(), Args, reloadDefaults
	t.Cleanup(func() {
		RootCmd, runtimeViper, Args, reloadDefaults = oldRoot, oldViper, oldArgs, oldDefaults
		*config.GetCfg() = oldConfig
	})
	config.ResetConfigForRuntime()
	runtimeViper = viper.New()
	RootCmd = &cobra.Command{Use: "comi", RunE: oldRoot.RunE, PersistentPreRunE: oldRoot.PersistentPreRunE,
		Args: cobra.ArbitraryArgs, Version: oldRoot.Version, SilenceUsage: true, SilenceErrors: true}
	RootCmd.SetVersionTemplate(oldRoot.VersionTemplate())
	InitFlags()
	addProcessCommands()
	return RootCmd
}

// 托盘入口不能启动没有管理通道的后台子进程；用子进程隔离启动壳标记。
func TestTrayRejectsCLIProcessCommands(t *testing.T) {
	if os.Getenv("COMI_TRAY_COMMAND_TEST") == "1" {
		config.UseTrayConfigProfile()
		for _, name := range []string{"start", "go", "stop", "reload"} {
			processTestCommand(t)
			handled, err := RunProcessCommand([]string{name}, &bytes.Buffer{})
			if !handled || err == nil || err.Error() != locale.GetString("cli_process_commands_only") {
				t.Fatalf("%s 未拒绝 CLI 专用操作：handled=%v err=%v", name, handled, err)
			}
		}
		return
	}
	child := exec.Command(os.Args[0], "-test.run=^TestTrayRejectsCLIProcessCommands$")
	child.Env = append(os.Environ(), "COMI_TRAY_COMMAND_TEST=1")
	if output, err := child.CombinedOutput(); err != nil {
		t.Fatalf("托盘命令检查失败：%v\n%s", err, output)
	}
}

// 未登记管理通道的其他配置进程也必须显示。
func TestProcessStatusAcrossDirectories(t *testing.T) {
	if os.Getenv("COMI_STATUS_TEST_CHILD") == "1" {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		defer listener.Close()
		fmt.Println(listener.Addr().(*net.TCPAddr).Port)
		_, _ = bufio.NewReader(os.Stdin).ReadByte()
		return
	}
	root := processTestCommand(t)
	t.Setenv("COMIGO_CONFIG_DIR", t.TempDir())
	// 查询所有进程时，启动端口和另一配置目录不能隐藏当前子进程。
	if err := root.PersistentFlags().Set("port", "1"); err != nil {
		t.Fatal(err)
	}
	if err := root.PersistentFlags().Set("config", filepath.Join(t.TempDir(), "broken.toml")); err != nil {
		t.Fatal(err)
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	name := "comi"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	path := filepath.Join(t.TempDir(), name)
	if err := os.Link(executable, path); err != nil {
		t.Skipf("系统不支持临时可执行文件硬链接：%v", err)
	}
	child := exec.Command(path, "-test.run=^TestProcessStatusAcrossDirectories$")
	child.Env = append(os.Environ(), "COMI_STATUS_TEST_CHILD=1")
	stdout, err := child.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	stdin, err := child.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := child.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = stdin.Close(); _ = child.Process.Kill(); _ = child.Wait() })
	portLine := make(chan string, 1)
	go func() { line, _ := bufio.NewReader(stdout).ReadString('\n'); portLine <- strings.TrimSpace(line) }()
	var port string
	select {
	case port = <-portLine:
		if _, err := strconv.Atoi(port); err != nil {
			t.Fatalf("子进程未返回监听端口：%q", port)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("等待子进程监听超时")
	}
	var out bytes.Buffer
	if err := showProcessStatus(&out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), fmt.Sprintf("PID %d", child.Process.Pid)) || !strings.Contains(out.String(), port) {
		t.Fatalf("未找到其他目录的进程及端口：%s", out.String())
	}
}

// 三种语言的帮助和状态输出必须使用译文，不能泄漏 locale 键或凭据。
func TestProcessStatusLocalization(t *testing.T) {
	t.Cleanup(func() { locale.InitLanguageFromConfig("auto") })
	for _, lang := range []string{"zh", "en", "ja"} {
		t.Run(lang, func(t *testing.T) {
			root := processTestCommand(t)
			locale.InitLanguageFromConfig(lang)
			refreshCLIText()
			var out bytes.Buffer
			root.SetOut(&out)
			root.SetArgs([]string{"status", "--help"})
			if err := root.Execute(); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(out.String(), locale.GetString("cli_status")) || !strings.Contains(out.String(), locale.GetString("cli_status_details")) {
				t.Fatal(out.String())
			}
			for _, key := range []string{"cli_status_running", "cli_status_process", "cli_status_no_listener", "cli_status_ports_unavailable", "cli_status_scan_failed"} {
				if locale.GetString(key) == key {
					t.Fatalf("%s 缺少翻译 %s", lang, key)
				}
			}
			message := fmt.Sprintf(locale.GetString("cli_status_process"), 42, "1234")
			if !strings.Contains(message, "42") || !strings.Contains(message, "1234") || strings.Contains(message, "%!") {
				t.Fatal(message)
			}
		})
	}
}

// 验证各个子命令、原来的路径参数和 --version；管理命令不得加载损坏的配置。
func TestProcessCommands(t *testing.T) {
	root := processTestCommand(t)
	// 别名解析到同一个命令，参数规则和执行逻辑完全复用。
	for alias, name := range map[string]string{"ls": "status", "ps": "status", "go": "start"} {
		command, remaining, err := root.Find([]string{alias})
		original, _, _ := root.Find([]string{name})
		if err != nil || command != original || len(remaining) != 0 {
			t.Fatalf("%s 未解析为 %s: %v", alias, name, err)
		}
	}
	for _, name := range []string{"run", "start", "status", "stop", "reload", "upgrade", "version"} {
		command, _, err := root.Find([]string{name})
		if err != nil || command.Name() != name || command.HasSubCommands() {
			t.Fatalf("%s: command=%v err=%v", name, command, err)
		}
		if name != "run" && name != "start" && command.Args(command, []string{"extra"}) == nil {
			t.Fatalf("%s accepted extra argument", name)
		}
	}
	for _, args := range [][]string{nil, {"./run"}, {"books", "https://example.com/library/"}} {
		command, remaining, err := root.Find(args)
		if err != nil || command != root || !reflect.DeepEqual(args, remaining) {
			t.Fatalf("paths changed: %q -> %v, %q, %v", args, command, remaining, err)
		}
	}
	file := filepath.Join(t.TempDir(), "broken.toml")
	if err := os.WriteFile(file, []byte("broken = ["), 0o600); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	root.SetOut(&output)
	root.SetArgs([]string{"version", "--config", file})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	version := output.String()
	output.Reset()
	root.SetArgs([]string{"--version"})
	if err := root.Execute(); err != nil || output.String() != version {
		t.Fatalf("version differs: %q != %q (%v)", output.String(), version, err)
	}
}

// 验证 start 重建参数时不把配置文件名或书库名误认作命令，并保留空格及 -- 后的路径。
func TestBackgroundArgs(t *testing.T) {
	root := processTestCommand(t)
	command, _, _ := root.Find([]string{"start"})
	if err := command.ParseFlags([]string{"--config", "start", "--port", "2345", "--no-tui=false", "--", "books with spaces", "--version"}); err != nil {
		t.Fatal(err)
	}
	want := []string{"run", "--config=start", "--port=2345", "--no-tui", "--open-browser=false", "--", "books with spaces", "--version"}
	if got := backgroundArgs(command, command.Flags().Args()); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

// 验证重读 TOML 能删除旧值、保留启动 flag，且解析错误不污染当前配置。
func TestReadReloadConfig(t *testing.T) {
	root := processTestCommand(t)
	file := filepath.Join(t.TempDir(), "config.toml")
	write := func(content string) {
		t.Helper()
		if err := os.WriteFile(file, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	write("Port = 2345\nHost = 'first.example'\nLanguage = 'en'\n")
	root.SetArgs([]string{"run", "--config", file, "--max-depth", "7"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	write("Port = 3456\nMaxScanDepth = 2\n")
	candidate, err := readReloadConfig(file)
	if err != nil || candidate.Port != 3456 || candidate.Host != "" || candidate.MaxScanDepth != 7 {
		t.Fatalf("reload: port=%d host=%q depth=%d err=%v", candidate.Port, candidate.Host, candidate.MaxScanDepth, err)
	}
	Args = []string{t.TempDir()}
	write("StoreUrls = [" + strconv.Quote(Args[0]) + "]")
	if candidate, err := readReloadConfig(file); err != nil || len(candidate.StoreUrls) != 1 {
		t.Fatalf("duplicate startup path: %v, %v", candidate.StoreUrls, err)
	}
	old := config.CopyCfg()
	for _, content := range []string{"Port = [", "Port = -1", "EnableTLS=true\nCertFile = 'missing.pem'"} {
		write(content)
		if _, err := readReloadConfig(file); err == nil {
			t.Fatalf("accepted %q", content)
		}
		if !reflect.DeepEqual(old, config.CopyCfg()) {
			t.Fatal("failed reload changed running config")
		}
	}
}

// 验证控制通道需要令牌和 POST，并且在初始化完成前不接受管理操作。
func TestProcessControlAuthentication(t *testing.T) {
	control := &processControl{}
	handler := control.handler("secret")
	for _, test := range []struct {
		method, token string
		ready         bool
		status        int
	}{
		{http.MethodPost, "", true, http.StatusForbidden},
		{http.MethodPost, "wrong", true, http.StatusForbidden},
		{http.MethodGet, "secret", true, http.StatusForbidden},
		{http.MethodPost, "secret", false, http.StatusServiceUnavailable},
		{http.MethodPost, "secret", true, http.StatusNoContent},
	} {
		control.ready.Store(test.ready)
		request := httptest.NewRequest(test.method, "/ready", nil)
		request.Header.Set("Authorization", "Bearer "+test.token)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != test.status {
			t.Fatalf("got %d, want %d", response.Code, test.status)
		}
	}
	if err := processRequest(processState{Address: "example.com:80"}, "ready", ""); err == nil {
		t.Fatal("accepted non-loopback control address")
	}
}

// 验证同一进程不重复初始化管理通道，并且关闭后可重新创建。
func TestProcessControlLifecycle(t *testing.T) {
	processTestCommand(t)
	config.GetCfg().ConfigFile = filepath.Join(t.TempDir(), "config.toml")
	if err := StartProcessControl(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(CloseProcessControl)
	if err := StartProcessControl(); err == nil {
		t.Fatal("duplicate instance accepted")
	}
	ProcessReady()
	state, err := readProcessState(cliControl.path)
	if err != nil {
		t.Fatal(err)
	}
	if err := processRequest(state, "ready", ""); err != nil {
		t.Fatal(err)
	}
	CloseProcessControl()
	if err := StartProcessControl(); err != nil {
		t.Fatal(err)
	}
}

// 查询命令必须在损坏配置下直接完成，不能落入服务初始化。
func TestCLIQueriesReturnBeforeStartup(t *testing.T) {
	for _, args := range [][]string{{"help"}, {"help", "run"}, {"status"}, {"ps"}, {"completion", "bash"}, {"__complete", "run", "--"}, {"--lang=zh", "--help"}} {
		t.Run(args[0], func(t *testing.T) {
			processTestCommand(t)
			file := filepath.Join(t.TempDir(), "broken.toml")
			if err := os.WriteFile(file, []byte("broken=["), 0600); err != nil {
				t.Fatal(err)
			}
			var output bytes.Buffer
			handled, err := RunProcessCommand(append(args, "--config", file), &output)
			if !handled || err != nil || output.Len() == 0 {
				t.Fatalf("handled=%v err=%v output=%q", handled, err, output.String())
			}
			if args[0] == "--lang=zh" && !strings.Contains(output.String(), "登录 Cookie/JWT") {
				t.Fatalf("help not localized: %s", output.String())
			}
		})
	}
}

// 重载多实例时，只有唯一匹配或显式启动端口才能确定目标。
func TestSelectMultipleProcesses(t *testing.T) {
	root := processTestCommand(t)
	file := filepath.Join(t.TempDir(), "config.toml")
	config.GetCfg().ConfigFile = file
	base, err := processStatePath()
	if err != nil {
		t.Fatal(err)
	}
	for i, port := range []int{2345, 3456} {
		control := &processControl{}
		control.ready.Store(true)
		server := httptest.NewServer(control.handler("fixture"))
		t.Cleanup(server.Close)
		state := processState{PID: i + 1, Port: port, ConfigFile: file, Address: strings.TrimPrefix(server.URL, "http://"), Token: "fixture"}
		data, _ := json.Marshal(state)
		if err := os.WriteFile(processInstancePath(base, i+1), data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	if _, _, err := findProcessState(base); err == nil {
		t.Fatal("ambiguous instances accepted")
	}
	if err := root.PersistentFlags().Set("port", "3456"); err != nil {
		t.Fatal(err)
	}
	_, state, err := findProcessState(base)
	if err != nil || state.Port != 3456 {
		t.Fatalf("state=%v err=%v", state, err)
	}
}

// 查询和管理命令必须排除，但同名配置值和书库路径不能误排除。
func TestProcessInvocation(t *testing.T) {
	root := processTestCommand(t)
	root.InitDefaultHelpCmd()
	root.InitDefaultCompletionCmd()
	for _, test := range []struct {
		args      []string
		file      string
		transient bool
	}{
		{[]string{"status"}, "", true},
		{[]string{"ls"}, "", true},
		{[]string{"--config", "books.toml", "ps"}, "", true},
		{[]string{"stop"}, "", true},
		{[]string{"help", "run"}, "", true},
		{[]string{"run", "--help"}, "", true},
		{[]string{"-v"}, "", true},
		{[]string{"--upgrade"}, "", true},
		{[]string{"--upgrade=1"}, "", true},
		{[]string{"run", "-c", "status", "--", "--version"}, "status", false},
		{[]string{"--config=books with spaces.toml", "--", "status"}, "books with spaces.toml", false},
		{[]string{"--config", "--version", "run"}, "--version", false},
		{[]string{"--version=false", "./ps"}, "", false},
		{nil, "", false},
	} {
		file, transient := processInvocation(test.args)
		if file != test.file || transient != test.transient {
			t.Errorf("%q: got (%q, %v), want (%q, %v)", test.args, file, transient, test.file, test.transient)
		}
	}
	if root.PersistentFlags().Changed("config") || config.GetCfg().ConfigFile != "" {
		t.Fatal("发现其他进程时改写了当前配置")
	}
}

// 状态来自实时管理通道：排除自身及残留文件，重载后使用实际端口。
func TestProcessStatusUsesLivePort(t *testing.T) {
	processTestCommand(t)
	dir := t.TempDir()
	t.Setenv("COMIGO_CONFIG_DIR", dir)
	config.GetCfg().Port = 23456
	pid := os.Getpid()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(processState{PID: pid, Port: 23456})
	}))
	defer server.Close()
	base, err := processStatePath()
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	for _, instancePID := range []int{os.Getpid(), 100000} {
		pid = instancePID
		state := processState{PID: pid, Port: 1234, Address: strings.TrimPrefix(server.URL, "http://"), Token: "fixture"}
		data, _ := json.Marshal(state)
		if err := os.WriteFile(processInstancePath(base, pid), data, 0600); err != nil {
			t.Fatal(err)
		}
		out.Reset()
		if err := showProcessStatus(&out); err != nil {
			t.Fatal(err)
		}
		shown := strings.Contains(out.String(), fmt.Sprintf("PID %d", pid))
		if shown != (pid != os.Getpid()) || (shown && !strings.Contains(out.String(), "23456")) || strings.Contains(out.String(), "fixture") {
			t.Fatal(out.String())
		}
	}
	server.Close()
	out.Reset()
	if err := showProcessStatus(&out); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out.String(), "23456") {
		t.Fatal("残留文件被当成存活实例", out.String())
	}
}
