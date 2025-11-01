package sse_hub

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func SSEHandler(c echo.Context) error {
	// 生产环境：可以检查 c.Request().Header.Get("Last-Event-ID") 实现断点续传
	id := fmt.Sprintf("%d", time.Now().UnixNano())
	ch := make(chan Event, 16)
	MessageHub.Add(id, ch)
	defer MessageHub.Remove(id)

	res := c.Response()
	res.Header().Set(echo.HeaderContentType, "text/event-stream; charset=utf-8")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Connection", "keep-alive")
	res.Header().Set("X-Accel-Buffering", "no") // Nginx 关闭缓冲（很重要）
	// 如果用了 gzip 中间件，请为该路由跳过 gzip（SSE 不宜压缩）。

	flusher, ok := res.Writer.(http.Flusher)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "streaming unsupported")
	}

	// 设置客户端自动重连间隔（毫秒）
	_, _ = res.Write([]byte("retry: 10000\n\n"))
	// 发送一条注释防止某些代理缓冲
	_, _ = res.Write([]byte(": connected\n\n"))
	flusher.Flush()

	keep := time.NewTicker(15 * time.Second) // 心跳，避免中间代理断开
	defer keep.Stop()

	for {
		select {
		case <-c.Request().Context().Done():
			return nil
		case <-keep.C:
			_, _ = res.Write([]byte(": ping\n\n"))
			flusher.Flush()
		case ev, ok := <-ch:
			if !ok {
				return nil
			}
			if ev.ID != "" {
				_, _ = res.Write([]byte("id: " + ev.ID + "\n"))
			}
			if ev.Name != "" {
				_, _ = res.Write([]byte("event: " + ev.Name + "\n"))
			}
			for _, line := range splitLines(ev.Data) {
				_, _ = res.Write([]byte("data: " + line + "\n"))
			}
			_, _ = res.Write([]byte("\n"))
			flusher.Flush()
		}
	}
}
