package data_api

import (
	"crypto/rand"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/routers/login"
)

// Connection 表示一个实时连接；设备按浏览器 Cookie 去重，用户按登录账号去重。
type Connection struct {
	ID          string    `json:"id"`
	Device      string    `json:"device"`
	User        string    `json:"user"`
	IP          string    `json:"ip"`
	UserAgent   string    `json:"userAgent"`
	ConnectedAt time.Time `json:"connectedAt"`
}

var online = struct {
	sync.Mutex
	connections map[string]Connection
}{connections: make(map[string]Connection)}

// DeviceCookie 在页面响应时建立设备标识，供同一浏览器的多个标签页共享。
func DeviceCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("comigo_device")
		if err != nil || cookie.Value == "" {
			cookie = &http.Cookie{Name: "comigo_device", Value: rand.Text(), Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: c.IsTLS(), MaxAge: 365 * 24 * 60 * 60}
			c.SetCookie(cookie)
		}
		c.Set("device", cookie.Value)
		return next(c)
	}
}

// TrackConnection 只登记通过认证的 SSE/WebSocket 连接，断开时立即移除。
func TrackConnection(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		connection := Connection{ID: rand.Text(), Device: c.Get("device").(string), IP: c.RealIP(), UserAgent: c.Request().UserAgent(), ConnectedAt: time.Now()}
		if token, ok := c.Get("user").(*jwt.Token); ok {
			connection.User = token.Claims.(*login.JwtCustomClaims).Username
		}
		online.Lock()
		online.connections[connection.ID] = connection
		online.Unlock()
		defer func() { online.Lock(); delete(online.connections, connection.ID); online.Unlock() }()
		return next(c)
	}
}

// connectionSnapshot 返回快照及去重统计；访客只有设备数，不冒充已登录用户。
func connectionSnapshot() ([]Connection, int, int) {
	online.Lock()
	defer online.Unlock()
	connections := make([]Connection, 0, len(online.connections))
	users, devices := map[string]bool{}, map[string]bool{}
	for _, connection := range online.connections {
		connections = append(connections, connection)
		devices[connection.Device] = true
		if connection.User != "" {
			users[connection.User] = true
		}
	}
	sort.Slice(connections, func(i, j int) bool { return connections[i].ConnectedAt.Before(connections[j].ConnectedAt) })
	return connections, len(users), len(devices)
}

// GetConnections 返回当前在线连接和登录用户、浏览器设备数量。
func GetConnections(c echo.Context) error {
	connections, users, devices := connectionSnapshot()
	return c.JSON(http.StatusOK, echo.Map{"connections": connections, "onlineUsers": users, "onlineDevices": devices})
}
