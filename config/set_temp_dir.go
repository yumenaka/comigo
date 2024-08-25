package config

import (
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/locale"
	"github.com/yumenaka/comigo/util/logger"
	"os"
	"path"
)

// SetTempDir 设置临时文件夹，退出时会被清理
func SetTempDir() {
	//手动设置的临时文件夹
	if Config.CachePath != "" && util.IsExist(Config.CachePath) && util.ChickIsDir(Config.CachePath) {
		Config.CachePath = path.Join(Config.CachePath)
	} else {
		Config.CachePath = path.Join(os.TempDir(), "comigo_cache") //使用系统文件夹
	}
	err := os.MkdirAll(Config.CachePath, os.ModePerm)
	if err != nil {
		logger.Infof("%s", locale.GetString("temp_folder_error"))
	} else {
		logger.Infof("%s", locale.GetString("temp_folder_path")+Config.CachePath)
	}
}
