package config

import (
	"errors"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/logger"
	"os"
	"path"
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
	logger.Info("Checking cfg ShareName")

	// 初始化
	c.In = "None"
	c.Path.WorkingDirectory = ""
	c.Path.HomeDirectory = ""
	c.Path.ProgramDirectory = ""

	// 配置文件搜索路径及顺序
	type configPath struct {
		name string
		path string
	}
	var configPaths []configPath

	// 获取可执行文件所在目录
	executablePath, err := os.Executable()
	if err != nil {
		return errors.New("error: failed to find executable path")
	}
	programDir := path.Dir(executablePath)

	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.New("error: failed to find home directory")
	}

	// 添加搜索路径
	configPaths = append(configPaths, configPath{
		name: "HomeDirectory",
		path: util.GetAbsPath(path.Join(homeDir, ".config/comigo/config.toml")),
	})
	configPaths = append(configPaths, configPath{
		name: "ProgramDirectory",
		path: util.GetAbsPath(path.Join(programDir, "config.toml")),
	})
	configPaths = append(configPaths, configPath{
		name: "WorkingDirectory",
		path: util.GetAbsPath("config.toml"),
	})

	// 按顺序检查配置文件是否存在
	for _, cp := range configPaths {
		if util.IsExist(cp.path) {
			switch cp.name {
			case "HomeDirectory":
				c.Path.HomeDirectory = cp.path
			case "ProgramDirectory":
				c.Path.ProgramDirectory = cp.path
			case "WorkingDirectory":
				c.Path.WorkingDirectory = cp.path
			}
			c.In = cp.name
			return nil
		}
	}
	return nil
}

type UploadDirOption int

var CfgStatus = Status{}
