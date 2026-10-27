//go:build wails && !android && !js && !bindings

package main

// ensureWailsAndroidImportStore 只在 Android Wails 中需要默认导入书库。
func ensureWailsAndroidImportStore() error {
	return nil
}
