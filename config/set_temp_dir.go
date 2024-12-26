package config

import (
	"os"
	"path"

	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/locale"
	"github.com/yumenaka/comigo/util/logger"
)

// SetTempDir 设置临时文件夹，退出时会被清理
func SetTempDir() {
	//手动设置的临时文件夹
	if Cfg.CachePath != "" && util.IsExist(Cfg.CachePath) && util.ChickIsDir(Cfg.CachePath) {
		Cfg.CachePath = path.Join(Cfg.CachePath)
	} else {
		Cfg.CachePath = path.Join(os.TempDir(), "comigo_cache") //使用系统文件夹
	}
	err := os.MkdirAll(Cfg.CachePath, os.ModePerm)
	if err != nil {
		logger.Infof("%s", locale.GetString("temp_folder_error"))
	} else {
		logger.Infof("%s", locale.GetString("temp_folder_path")+Cfg.CachePath)
	}
}
