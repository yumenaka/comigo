//go:build linux

// Package autostart manages the current Comigo program's login startup.
package autostart

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/kardianos/service"
	"github.com/yumenaka/comigo/config"
)

var mutation sync.Mutex

// State 区分尚未启用和无法查询；不把系统状态写进 TOML。
type State struct {
	Enabled   bool   `json:"enabled"`
	Supported bool   `json:"supported"`
	Scope     string `json:"scope"`
}

// unitName 各启动壳独立注册，避免桌面、托盘和 CLI 相互覆盖。
func unitName() string {
	return "comigo-" + strings.TrimSuffix(config.PlatformConfigFilename(), ".toml")
}

func unitPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "systemd", "user", unitName()+".service"), nil
}

// Status 读取实际启用状态；注册文件缺失代表默认关闭。
func Status() (State, error) {
	state := State{Supported: true, Scope: "user"}
	if _, err := os.Stat("/run/systemd/system"); err != nil {
		state.Supported = false
		return state, nil
	}
	path, err := unitPath()
	if err != nil {
		return state, err
	}
	if _, err = os.Lstat(path); errors.Is(err, os.ErrNotExist) {
		return state, nil
	} else if err != nil {
		return state, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, "/usr/bin/systemctl", "--user", "is-enabled", unitName()+".service")
	command.Env = append(os.Environ(), "LC_ALL=C")
	output, err := command.Output()
	switch strings.TrimSpace(string(output)) {
	case "enabled", "enabled-runtime":
		state.Enabled = true
		return state, nil
	case "disabled", "not-found", "masked", "masked-runtime":
		return state, nil
	}
	if err != nil {
		return state, fmt.Errorf("query startup: %w", err)
	}
	return state, fmt.Errorf("unexpected startup state")
}

// serviceConfig 使用明确配置文件启动当前二进制，服务启动禁止默认书库与 TUI。
func serviceConfig() (*service.Config, error) {
	executable, err := os.Executable()
	if err != nil {
		return nil, err
	}
	configuration, err := filepath.Abs(config.GetCfg().ConfigFile)
	if err != nil {
		return nil, err
	}
	arguments := []string{"--no-default-library", "--config", configuration}
	target, dependencies := "default.target", ""
	if config.PlatformConfigFilename() == "config.toml" {
		arguments = append(arguments, "--no-tui", "--open-browser=false")
	} else {
		target = "graphical-session.target"
		dependencies = "After=graphical-session.target\nPartOf=graphical-session.target"
	}
	// 库负责双引号；补齐反斜线、控制字符和 systemd 百分号转义。
	// 路径仅作为模板数据传入，不能让文件名中的 {{...}} 成为模板指令。
	escape := strings.NewReplacer("\\", "\\\\", "\n", "\\n", "\r", "\\r", "\t", "\\t", "\a", "\\a", "\b", "\\b", "\f", "\\f", "\v", "\\v", "%", "%%")
	for i := range arguments {
		arguments[i] = escape.Replace(arguments[i])
	}
	unit := strings.NewReplacer("@TARGET@", target, "@DEPENDENCIES@", dependencies).Replace(userUnit)
	return &service.Config{Name: unitName(), DisplayName: "Comigo", Description: "Comigo reader", Executable: escape.Replace(executable), Arguments: arguments,
		Option: service.KeyValue{"UserService": true, "SystemdScript": unit}}, nil
}

// Set 只注册／撤销登录启动，不启动第二个进程，也不停止当前阅读服务。
func Set(enabled bool) error {
	mutation.Lock()
	defer mutation.Unlock()
	state, err := Status()
	if err != nil {
		return err
	}
	if !state.Supported {
		return errors.New("user services require systemd")
	}
	path, err := unitPath()
	if err != nil {
		return err
	}
	if info, err := os.Lstat(path); err == nil {
		if !info.Mode().IsRegular() || info.Mode().Perm()&0022 != 0 || info.Size() > 64*1024 {
			return errors.New("unsafe existing service file")
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !strings.HasPrefix(string(content), "# Managed by Comigo\n") {
			return errors.New("service file belongs to another installation")
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if state.Enabled == enabled {
		return nil
	}
	if enabled {
		if config.GetCfg().TemporaryReaderMode {
			return errors.New("save a persistent library configuration before enabling startup")
		}
		if err := config.UpdateConfigFile(); err != nil {
			return err
		}
	}
	cfg, err := serviceConfig()
	if err != nil {
		return err
	}
	manager, err := service.New(nil, cfg)
	if err != nil {
		return err
	}
	if !enabled {
		return manager.Uninstall()
	}
	// 已有但未启用的本程序单元先经库卸载，再按当前配置注册。
	if _, err := os.Lstat(path); err == nil {
		if err := manager.Uninstall(); err != nil {
			return err
		}
	}
	if err := manager.Install(); err != nil {
		return errors.Join(err, manager.Uninstall())
	}
	state, err = Status()
	if err != nil {
		return err
	}
	if !state.Enabled {
		return errors.New("startup was not enabled")
	}
	return nil
}

const userUnit = `# Managed by Comigo
[Unit]
Description={{Description}}
@DEPENDENCIES@

[Service]
Type=simple
ExecStart=:{{Path | cmd}}{{range Arguments}} {{. | cmd}}{{end}}
Restart=on-failure
RestartSec=10

[Install]
WantedBy=@TARGET@
`
