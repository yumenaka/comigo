package tailscale_plugin

import (
	"fmt"
	"strings"

	"github.com/yumenaka/comigo/tools/logger"
)

// tailscaleUserLogf 接管 tsnet.Server.UserLogf，避免第三方库绕过项目 logger 直接写入标准输出。
func tailscaleUserLogf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	if shouldSuppressTailscaleUserLog(message) {
		return
	}
	logger.Info(message)
}

// shouldSuppressTailscaleUserLog 过滤 tsnet 的内部状态噪音。
// AuthLoop 行只表示 tsnet 的认证循环已经结束，不是用户需要处理的事件；TUI 中显示它会干扰当前状态面板。
func shouldSuppressTailscaleUserLog(message string) bool {
	return strings.HasPrefix(message, "AuthLoop:")
}
