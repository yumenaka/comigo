package apiresp

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type ErrorBody struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

type SuccessBody struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

func Error(c echo.Context, status int, code string, message string, details interface{}) error {
	return c.JSON(status, ErrorBody{
		Code:      code,
		Message:   message,
		Details:   details,
		RequestID: requestID(c),
		Timestamp: time.Now().Unix(),
	})
}

func BadRequest(c echo.Context, code string, message string, details interface{}) error {
	return Error(c, http.StatusBadRequest, code, message, details)
}

func Forbidden(c echo.Context, code string, message string, details interface{}) error {
	return Error(c, http.StatusForbidden, code, message, details)
}

func Success(c echo.Context, code string, message string, data interface{}) error {
	return c.JSON(http.StatusOK, SuccessBody{
		Code:      code,
		Message:   message,
		Data:      data,
		RequestID: requestID(c),
		Timestamp: time.Now().Unix(),
	})
}

func requestID(c echo.Context) string {
	if rid := c.Response().Header().Get(echo.HeaderXRequestID); rid != "" {
		return rid
	}
	return c.Request().Header.Get(echo.HeaderXRequestID)
}
