package config

import (
	"runtime"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/logger"
)

type Status struct {
	// 当前生效的配置文件路径 None、HomeDirectory、WorkingDirectory、ProgramDirectory
	// 设置读取顺序：None（默认值） -> HomeDirectory -> ProgramDirectory -> WorkingDirectory
	In   string
	Path struct {
		// 对应配置文件的绝对路径
		WorkingDirectory string
		HomeDirectory    string
		ProgramDirectory string
	}
}

func (c *Status) SetConfigStatus() error {
	// 在js环境下
	if runtime.GOOS == "js" {
		return nil
	}
	// 初始化
	c.In = "None"
	c.Path.WorkingDirectory = ""
	c.Path.HomeDirectory = ""
	c.Path.ProgramDirectory = ""
	if cfg.TemporaryReaderMode {
		return nil
	}
	logger.Info(locale.GetString("log_checking_cfg_sharename"))

	for _, cp := range []struct {
		name string
		path string
	}{
		{name: HomeDirectory, path: GetHomeDirectoryConfig()},
		{name: ProgramDirectory, path: GetProgramDirectoryConfig()},
		{name: WorkingDirectory, path: GetWorkingDirectoryConfig()},
	} {
		if cp.path == "" {
			continue
		}
		switch cp.name {
		case HomeDirectory:
			c.Path.HomeDirectory = cp.path
		case ProgramDirectory:
			c.Path.ProgramDirectory = cp.path
		case WorkingDirectory:
			c.Path.WorkingDirectory = cp.path
		}
		c.In = cp.name
		return nil
	}
	return nil
}

type UploadDirOption int

var CfgStatus = Status{}
