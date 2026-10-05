package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	RootCmd = &cobra.Command{Use: "comi", Run: oldRoot.Run, PersistentPreRunE: oldRoot.PersistentPreRunE,
		Args: cobra.ArbitraryArgs, Version: oldRoot.Version, SilenceUsage: true, SilenceErrors: true}
	RootCmd.SetVersionTemplate(oldRoot.VersionTemplate())
	InitFlags()
	addProcessCommands()
	return RootCmd
}

// 验证六个子命令、原来的路径参数和 --version；管理命令不得加载损坏的配置。
func TestProcessCommands(t *testing.T) {
	root := processTestCommand(t)
	for _, name := range []string{"run", "start", "stop", "reload", "upgrade", "version"} {
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
	for _, content := range []string{"Port = [", "Port = -1", "CertFile = 'missing.pem'"} {
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

// 验证文件锁拒绝重复实例，并且关闭管理通道后可以再次启动。
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
