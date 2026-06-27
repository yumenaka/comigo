//go:build wails && !js && !bindings && linux && !android

package wails_systray

import dbus "github.com/godbus/dbus/v5"

// traySupported 检测 Linux 桌面是否有 StatusNotifier/AppIndicator 宿主。
func traySupported() bool {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return false
	}
	defer conn.Close()

	var hasWatcher bool
	err = conn.BusObject().Call("org.freedesktop.DBus.NameHasOwner", 0, "org.kde.StatusNotifierWatcher").Store(&hasWatcher)
	return err == nil && hasWatcher
}
