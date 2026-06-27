//go:build wails && !js && !bindings && !linux && !darwin && !android

package wails_systray

import _ "embed"

//go:embed icon.ico
var trayIcon []byte
