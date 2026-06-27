//go:build wails && !js && !bindings && darwin

package wails_systray

/*
#cgo darwin CFLAGS: -x objective-c -fobjc-arc
#cgo darwin LDFLAGS: -framework Cocoa
#include <stdlib.h>

void comigoWailsTrayStart(void *iconBytes, int iconLen, char *tooltip, char *showTitle, char *showTip, char *quitTitle, char *quitTip);
void comigoWailsTrayStop(void);
void comigoWailsTraySetWindowVisible(int visible);
*/
import "C"

import (
	"sync"
	"unsafe"

	"github.com/yumenaka/comigo/assets/locale"
)

var darwinTray struct {
	sync.Mutex
	tray *Tray
}

// startPlatform 使用原生 NSStatusItem，避开 systray/Wails delegate 冲突。
func startPlatform(t *Tray) func() {
	darwinTray.Lock()
	darwinTray.tray = t
	darwinTray.Unlock()

	tooltip := C.CString(locale.GetString("systray_tooltip"))
	showTitle := C.CString(locale.GetString("wails_systray_show"))
	showTip := C.CString(locale.GetString("wails_systray_show_tooltip"))
	quitTitle := C.CString(locale.GetString("systray_quit"))
	quitTip := C.CString(locale.GetString("systray_quit_tooltip"))
	defer C.free(unsafe.Pointer(tooltip))
	defer C.free(unsafe.Pointer(showTitle))
	defer C.free(unsafe.Pointer(showTip))
	defer C.free(unsafe.Pointer(quitTitle))
	defer C.free(unsafe.Pointer(quitTip))

	C.comigoWailsTrayStart(unsafe.Pointer(&trayIcon[0]), C.int(len(trayIcon)), tooltip, showTitle, showTip, quitTitle, quitTip)
	return func() {
		C.comigoWailsTrayStop()
		darwinTray.Lock()
		darwinTray.tray = nil
		darwinTray.Unlock()
	}
}

func setPlatformWindowVisible(visible bool) {
	if visible {
		C.comigoWailsTraySetWindowVisible(1)
		return
	}
	C.comigoWailsTraySetWindowVisible(0)
}

func quitPlatformFallback() {
	C.comigoWailsTrayStop()
}

//export comigoWailsTrayShow
func comigoWailsTrayShow() {
	darwinTray.Lock()
	t := darwinTray.tray
	darwinTray.Unlock()
	if t != nil {
		go t.showWindow()
	}
}

//export comigoWailsTrayQuit
func comigoWailsTrayQuit() {
	darwinTray.Lock()
	t := darwinTray.tray
	darwinTray.Unlock()
	if t != nil {
		go t.quit()
	}
}
