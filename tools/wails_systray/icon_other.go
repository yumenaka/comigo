//go:build wails && !js && !bindings && !linux && !darwin

package wails_systray

import _ "embed"

//go:embed icon.ico
var trayIcon []byte
