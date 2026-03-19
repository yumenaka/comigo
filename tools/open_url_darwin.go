//go:build darwin

package tools

import (
	"fmt"
	"os/exec"

	"github.com/yumenaka/comigo/tools/logger"
)

// openURL 在 macOS 上打开 URL。
func openURL(uri string) error {
	cmd := exec.Command("open", uri)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("open url: %w", err)
	}
	logger.Infof("Opening URL: %s", uri)
	return nil
}
