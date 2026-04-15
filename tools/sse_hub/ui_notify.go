package sse_hub

import "encoding/json"

// 前端 ui_suggest_reload 事件使用的 reason，与 assets/locale 中 i18n key 后缀一致。
const (
	UISuggestReasonLibraryRescan     = "library_rescan_done"
	UISuggestReasonAutoLibraryRescan = "auto_library_rescan_done"
	UISuggestReasonSingleStoreRescan = "single_store_rescan_done"
	UISuggestReasonDebugToggle       = "debug_toggle"
	UISuggestReasonPluginsChanged    = "plugins_changed"
	UISuggestReasonServerConfig      = "server_config_changed"
	UISuggestReasonLoginSettings     = "login_settings_changed"
)

type uiSuggestPayload struct {
	Reason string `json:"reason"`
}

// BroadcastUISuggestReload 通知浏览器：建议用户确认后刷新整页（SSE event: ui_suggest_reload）。
func BroadcastUISuggestReload(reason string) {
	if MessageHub == nil || reason == "" {
		return
	}
	b, err := json.Marshal(uiSuggestPayload{Reason: reason})
	if err != nil {
		return
	}
	MessageHub.Broadcast(Event{Name: "ui_suggest_reload", Data: string(b)})
}
