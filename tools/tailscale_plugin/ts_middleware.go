package tailscale_plugin

import (
	"time"

	"github.com/labstack/echo/v4"
)

// TailscaleAuthMiddleware 创建Tailscale身份验证中间件
func TailscaleAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 中间件仅在检测到Tailscale标记时执行 WhoIs 查询（普通请求触发 Tailscale WhoIs 查询会导致2秒左右的延迟）
			if v, ok := c.Request().Context().Value(ctxKeyIsTailscale).(bool); !ok || !v {
				// fmt.Println("Not a Tailscale request, skipping WhoIs lookup.")
				return next(c)
			}
			// 避免重复查询
			if c.Get("tailscale_user") != nil && c.Get("tailscale_node") != "" {
				if nowStatus.CheckClientInfoExists(
					c.Get("tailscale_user").(string),
					c.Get("tailscale_node").(string),
					c.Get("tailscale_remote_addr").(string),
				) {
					// 已经存在此user, 说明访问者信息已经保存过了, 直接返回
					// fmt.Println("Tailscale request detected, but user info already exists, skipping WhoIs lookup.")
					return next(c)
				}
			}
			// fmt.Println("Tailscale request detected, performing WhoIs lookup.")
			// 检测到Tailscale标记, 获取请求者的Tailscale身份信息
			who, err := localClient.WhoIs(c.Request().Context(), c.Request().RemoteAddr)
			// 如果可以获取访问者的身份信息
			if err == nil {
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
