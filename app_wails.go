//go:build wails && !js

package main

import "github.com/yumenaka/comigo/routers"

// DeleteBookFile 由 Wails WebView 绑定调用；子路由优先走 /api/wails/delete-book-file。
func (a *App) DeleteBookFile(bookID string) (bool, error) {
	return routers.DeleteBookFileForWails(bookID)
}
