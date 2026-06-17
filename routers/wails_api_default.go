//go:build !wails

package routers

import "github.com/labstack/echo/v4"

func bindWailsAPI(group *echo.Group) {}
