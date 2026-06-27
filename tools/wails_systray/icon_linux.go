//go:build wails && !js && !bindings && linux && !android

package wails_systray

import _ "embed"

//go:embed icon.png
var trayIcon []byte
