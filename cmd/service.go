package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/yumenaka/comigo/tools/autostart"
)

// serviceCommand 在加载书库和启动界面前处理服务设置，三个程序入口共用。
func serviceCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "service status | autostart true|false",
		Short: "Manage Comigo login startup without opening a library",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(command *cobra.Command, args []string) error {
			switch args[0] {
			case "status":
				if len(args) != 1 {
					return fmt.Errorf("service status takes no arguments")
				}
			case "autostart":
				if len(args) != 2 || (args[1] != "true" && args[1] != "false") {
					return fmt.Errorf("usage: comi service autostart true|false")
				}
				value, _ := strconv.ParseBool(args[1])
				if err := autostart.Set(value); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unknown service command %q", args[0])
			}
			state, err := autostart.Status()
			if err != nil {
				return err
			}
			return json.NewEncoder(command.OutOrStdout()).Encode(state)
		},
	}
}
