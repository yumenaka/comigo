//go:build wails && android && !js && !bindings

package main

import (
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools/logger"
)

// ensureWailsAndroidImportStore 给 Android 上传/导入页提供一个真实可写的内部书库。
func ensureWailsAndroidImportStore() error {
	importDir := filepath.Join(application.Android.StoragePath(), "library", "imports")
	if err := os.MkdirAll(importDir, 0o755); err != nil {
		return err
	}
	for _, storeURL := range config.GetCfg().StoreUrls {
		if storeURL == importDir {
			return nil
		}
	}
	if err := config.GetCfg().AddStoreUrl(importDir); err != nil {
		logger.Infof(locale.GetString("log_failed_to_add_store_url"), err)
	}
	return nil
}
