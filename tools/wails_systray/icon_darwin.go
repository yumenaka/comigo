//go:build wails && !js && !bindings && darwin

package wails_systray

import _ "embed"

//go:embed icon.png
var trayIcon []byte
