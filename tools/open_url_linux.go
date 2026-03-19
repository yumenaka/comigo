//go:build linux

package tools

import (
	"fmt"
	"os/exec"

	"github.com/yumenaka/comigo/tools/logger"
)

// openURL 在 Linux 上打开 URL。
func openURL(uri string) error {
	cmd := exec.Command("xdg-open", uri)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("open url: %w", err)
	}
	logger.Infof("Opening URL: %s", uri)
	return nil
}
