package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools/releases"
)

// RunDesktop 在服务初始化前处理管理命令；不会扫描书库或启动监听。
func RunDesktop(args []string, out io.Writer) (bool, error) {
	InitFlags()
	command, _, findErr := RootCmd.Find(args)
	if findErr == nil && command.Name() == "service" {
		cobra.OnInitialize(LoadConfigFile)
		RootCmd.SetArgs(args)
		RootCmd.SetOut(out)
		return true, RootCmd.Execute()
	}
	if len(args) == 0 || args[0] != "desktop" {
		return false, nil
	}
	if len(args) < 2 {
		return true, errors.New("usage: comi desktop info|check-update")
	}
	switch args[1] {
	case "info":
		if len(args) != 2 {
			return true, errors.New("desktop info takes no arguments")
		}
		return true, json.NewEncoder(out).Encode(struct {
			Version  string `json:"version"`
			Protocol int    `json:"desktopProtocol"`
		}{config.GetVersion(), 1})
	case "check-update":
		if len(args) != 2 {
			return true, errors.New("desktop check-update takes no arguments")
		}
		status, err := releases.Default.Check(context.Background(), config.GetVersion())
		if encodeErr := json.NewEncoder(out).Encode(status); encodeErr != nil {
			return true, encodeErr
		}
		return true, err
	default:
		return true, fmt.Errorf("unknown desktop command %q", args[1])
	}
}
