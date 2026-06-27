//go:build wails && !js && !bindings && linux

package wails_systray

import _ "embed"

//go:embed icon.png
var trayIcon []byte
