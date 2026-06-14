package data_api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/routers/apiresp"
	"github.com/yumenaka/comigo/tools/tailscale_plugin"
)

// GetTailscaleStatus 处理Tailscale网络的身份信息查询请求
func GetTailscaleStatus(c echo.Context) error {
	tailscaleStatus, err := tailscale_plugin.GetTailscaleStatus(c.Request().Context())
	if err != nil {
		return apiresp.Error(c, http.StatusInternalServerError, "tailscale_status_failed", err.Error(), nil)
	}
	return c.JSON(http.StatusOK, tailscaleStatus)
}

// GetTailscaleStatusSSE 通过 SSE 推送 Tailscale 状态，替代设置页里的轮询。
func GetTailscaleStatusSSE(c echo.Context) error {
	flusher, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "streaming unsupported")
	}

	res := c.Response()
	res.Header().Set(echo.HeaderContentType, "text/event-stream; charset=utf-8")
	res.Header().Set(echo.HeaderCacheControl, "no-cache")
	res.Header().Set("Connection", "keep-alive")
	res.Header().Set("X-Accel-Buffering", "no")
	res.WriteHeader(http.StatusOK)

	writeStatus := func() error {
		tailscaleStatus, err := tailscale_plugin.GetTailscaleStatus(c.Request().Context())
		if err != nil {
			payload, _ := json.Marshal(map[string]string{"error": err.Error()})
			return writeTailscaleSSE(res, flusher, "tailscale_status_error", payload)
		}
		payload, err := json.Marshal(tailscaleStatus)
		if err != nil {
			return err
		}
		return writeTailscaleSSE(res, flusher, "tailscale_status", payload)
	}

	if _, err := res.Write([]byte("retry: 5000\n\n")); err != nil {
		return err
	}
	if err := writeStatus(); err != nil {
		return err
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	heartbeat := time.NewTicker(15 * time.Second)
	defer heartbeat.Stop()
	for {
		select {
		case <-c.Request().Context().Done():
			return nil
		case <-ticker.C:
			if err := writeStatus(); err != nil {
				return err
			}
		case <-heartbeat.C:
			if _, err := res.Write([]byte(": ping\n\n")); err != nil {
				return err
			}
			flusher.Flush()
		}
	}
}

// writeTailscaleSSE 写单个 SSE 事件，逐行 data 可避免 JSON 中换行破坏事件格式。
func writeTailscaleSSE(res *echo.Response, flusher http.Flusher, event string, payload []byte) error {
	if _, err := res.Write([]byte("event: " + event + "\n")); err != nil {
		return err
	}
	for _, line := range strings.Split(string(payload), "\n") {
		if _, err := res.Write([]byte("data: " + line + "\n")); err != nil {
			return err
		}
	}
	if _, err := res.Write([]byte("\n")); err != nil {
		return err
	}
	flusher.Flush()
	return nil
}
