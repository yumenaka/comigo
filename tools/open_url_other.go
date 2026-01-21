//go:build !windows && !darwin && !linux

package tools

import "fmt"

// openURL 在未知平台上的兜底实现。
func openURL(uri string) error {
	return fmt.Errorf("unsupported platform: cannot open url %q", uri)
}
