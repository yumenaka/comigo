//go:build wails && !js && !bindings && darwin

package wails_systray

/*
#cgo darwin CFLAGS: -x objective-c -fobjc-arc
#cgo darwin LDFLAGS: -framework Cocoa
#include <stdlib.h>

void comigoWailsTrayStart(void *iconBytes, int iconLen, char *tooltip,
	char *showTitle, char *showTip, char *copyTitle, char *copyTip,
	char *openDirTitle, char *openDirTip, char **storeURLs, int storeCount,
	char *extraTitle, char *extraTip, char *quitTitle, char *quitTip, char *versionTitle);
void comigoWailsTrayStop(void);
void comigoWailsTraySetWindowVisible(int visible);
*/
import "C"

import (
	"sync"
	"unsafe"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
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
	copyTitle := C.CString(locale.GetString("systray_copy_url"))
	copyTip := C.CString(locale.GetString("systray_copy_url_tooltip"))
	openDirTitle := C.CString(locale.GetString("systray_open_directory"))
	openDirTip := C.CString(locale.GetString("systray_open_directory_tooltip"))
	extraTitle := C.CString(locale.GetString("systray_extra"))
	extraTip := C.CString(locale.GetString("systray_extra_tooltip"))
	quitTitle := C.CString(locale.GetString("systray_quit"))
	quitTip := C.CString(locale.GetString("systray_quit_tooltip"))
	versionTitle := C.CString("Comigo " + config.GetVersion())
	defer C.free(unsafe.Pointer(tooltip))
	defer C.free(unsafe.Pointer(showTitle))
	defer C.free(unsafe.Pointer(showTip))
	defer C.free(unsafe.Pointer(copyTitle))
	defer C.free(unsafe.Pointer(copyTip))
	defer C.free(unsafe.Pointer(openDirTitle))
	defer C.free(unsafe.Pointer(openDirTip))
	defer C.free(unsafe.Pointer(extraTitle))
	defer C.free(unsafe.Pointer(extraTip))
	defer C.free(unsafe.Pointer(quitTitle))
	defer C.free(unsafe.Pointer(quitTip))
	defer C.free(unsafe.Pointer(versionTitle))

	storeURLs := config.GetCfg().StoreUrls
	storeCStrings := make([]*C.char, len(storeURLs))
	for i, storeURL := range storeURLs {
		storeCStrings[i] = C.CString(storeURL)
		defer C.free(unsafe.Pointer(storeCStrings[i]))
	}
	var storeURLsPtr **C.char
	if len(storeCStrings) > 0 {
		storeURLsPtr = (**C.char)(unsafe.Pointer(&storeCStrings[0]))
	}

	C.comigoWailsTrayStart(
		unsafe.Pointer(&trayIcon[0]), C.int(len(trayIcon)), tooltip,
		showTitle, showTip, copyTitle, copyTip,
		openDirTitle, openDirTip, storeURLsPtr, C.int(len(storeCStrings)),
		extraTitle, extraTip, quitTitle, quitTip, versionTitle,
	)
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

//export comigoWailsTrayCopyURL
func comigoWailsTrayCopyURL() {
	go copyReaderURL()
}

//export comigoWailsTrayOpenStore
func comigoWailsTrayOpenStore(storeURL *C.char) {
	path := C.GoString(storeURL)
	go openStoreDirectory(path)
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
