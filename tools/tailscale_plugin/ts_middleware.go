//go:build !wails

package tailscale_plugin

import (
	"time"

	"github.com/labstack/echo/v4"
	"tailscale.com/client/local"
)

// TailscaleAuthMiddleware 创建Tailscale身份验证中间件
func TailscaleAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			client, _ := c.Request().Context().Value(tailscaleClientKey{}).(*local.Client)
			if client == nil {
				return next(c)
			}
			who, err := client.WhoIs(c.Request().Context(), c.Request().RemoteAddr)
			// 如果可以获取访问者的身份信息
			if err == nil && who.Node != nil && who.UserProfile != nil {
				// 将最后访问者的Tailscale信息存储到上下文(避免重复查询)
				c.Set("tailscale_user", who.UserProfile.LoginName)
				c.Set("tailscale_node", firstLabel(who.Node.ComputedName))
				c.Set("tailscale_remote_addr", c.Request().RemoteAddr)
				// 将访问者的Tailscale信息存储到 nowStatus
				clientInfo := &TailsClientInfo{
					LoginUserName:     who.UserProfile.LoginName,
					NodeComputedName:  firstLabel(who.Node.ComputedName),
					RequestRemoteAddr: c.Request().RemoteAddr,
					AccessTime:        time.Now(),
				}
				nowStatus.AddClientInfo(clientInfo)
			}

			// 继续处理请求
			return next(c)
		}
	}
}
