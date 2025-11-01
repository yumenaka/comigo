package sse_hub

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// 服务器端广播接口（演示用）：POST /push { "message": "...", "event": "可选", "id": "可选" }
type pushPayload struct {
	Message string `json:"message" form:"message"`
	Event   string `json:"event"   form:"event"`
	ID      string `json:"id"      form:"id"`
}

// PushHandler 服务器端广播接口：POST /api/push { "message": "...", "event": "可选", "id": "可选" }
func PushHandler(c echo.Context) error {
	var p pushPayload
	if err := c.Bind(&p); err != nil {
		return err
	}
	fmt.Printf("Event %s ID %s DATA %s", p.Event, p.ID, p.Message)
	MessageHub.Broadcast(Event{
		Name: p.Event,
		ID:   p.ID,
		Data: p.Message,
	})
	return c.NoContent(http.StatusAccepted)
}
